package grpcOrchestrator

import (
	"context"
	"fmt"
	"math"
	"os"
	"slices"
	"time"

	"connectrpc.com/connect"
	grpcMiddleware "github.com/pandaci-com/pandaci/app/grpc/middleware"
	"github.com/pandaci-com/pandaci/pkg/utils"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
	"github.com/rs/zerolog/log"
)

func (h *Handler) CreateJob(ctx context.Context, req *connect.Request[pb.OrchestratorServiceCreateJobRequest]) (*connect.Response[pb.OrchestratorServiceCreateJobResponse], error) {
	log.Info().Msg("orchestrator create job")
	defer utils.MeasureTime(time.Now(), "orchestrator service create job")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	workflowRun, err := h.queries.Unsafe_GetWorkflowRunByID(ctx, claims.WorkflowID)
	if err != nil {
		log.Error().Err(err).Msg("getting workflow run")
		return nil, err
	}

	runner := "ubuntu-4x"
	if req.Msg.Runner != nil {
		runner = req.Msg.GetRunner()
	}

	if !slices.Contains([]string{"ubuntu-1x", "ubuntu-2x", "ubuntu-4x", "ubuntu-8x", "ubuntu-16x"}, runner) {
		log.Error().Msg("invalid runner")
		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Invalid runner",
			Message: fmt.Sprintf("Invalid runner: %s", runner),
		})
		return nil, fmt.Errorf("invalid runner")
	}

	org, err := h.queries.Unsafe_GetOrgByID(ctx, claims.OrgID)
	if err != nil {
		log.Error().Err(err).Msg("getting org")
		return nil, err
	}

	license, err := org.GetLicense()
	if err != nil {
		log.Error().Err(err).Msg("getting license")
		return nil, err
	}

	if license.Features.MaxCloudRunnerScale < types.GetBuildMinutesScale(types.CloudRunner(runner)) {
		log.Error().Msg("runner scale is too high for license")
		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Please upgrade your license",
			Message: fmt.Sprintf("Your plan limits you to ubuntu-%dx but you are trying to use %s", license.Features.MaxCloudRunnerScale, runner),
		})
		return nil, fmt.Errorf("runner scale is too high for license")
	}

	job := typesDB.JobRun{
		Name:          req.Msg.Name,
		WorkflowRunID: claims.WorkflowID,
		Status:        types.RunStatusPending,
		Runner:        runner,
	}

	if req.Msg.Skipped {
		job.Conclusion = types.Pointer(types.RunConclusionSkipped)
		job.Status = types.RunStatusCompleted
	}

	if err := h.queries.CreateJobRun(ctx, &job); err != nil {
		log.Error().Err(err).Interface("job", job).Msg("creating job in db")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to create job",
			Message: "Failed to create job in db",
		})

		return nil, err
	}

	log.Info().Msgf("created job with id: %s", job.ID)

	if req.Msg.Skipped {
		return connect.NewResponse(&pb.OrchestratorServiceCreateJobResponse{
			JobMeta: &pb.JobMeta{
				Id:     job.ID,
				Name:   job.Name,
				Runner: runner,
			},
		}), nil
	}

	startJobRes, err := h.runnerClient.StartJob(ctx, connect.NewRequest(&pb.RunnerServiceStartJobRequest{
		WorkflowMeta: req.Msg.WorkflowMeta,
		Id:           job.ID,
		Name:         job.Name,
		Runner:       runner,
	}))
	if err != nil {
		log.Error().Err(err).Msg("calling runner client to start job")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to create job",
			Message: "Failed to start job in runner",
		})

		if err := h.queries.FinishJobRunWithError(ctx, claims.WorkflowID, job.ID, "failed to start job"); err != nil {
			log.Error().Err(err).Msg("finishing job in db")
		}

		return nil, err
	}

	return connect.NewResponse(&pb.OrchestratorServiceCreateJobResponse{
		JobMeta: startJobRes.Msg.JobMeta,
	}), nil
}

func (h *Handler) FinishJob(ctx context.Context, req *connect.Request[pb.OrchestratorServiceFinishJobRequest]) (*connect.Response[pb.OrchestratorServiceFinishJobResponse], error) {
	defer utils.MeasureTime(time.Now(), "orchestrator service finish job")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	workflowRun, err := h.queries.Unsafe_GetWorkflowRunByID(ctx, claims.WorkflowID)
	if err != nil {
		log.Error().Err(err).Msg("getting workflow run")
		return nil, err
	}

	conclusion, err := types.RunOutputFromProto(req.Msg.Conclusion)
	if err != nil {
		log.Error().Err(err).Msg("getting output from proto")
		if err := h.queries.FinishJobRunWithError(ctx, claims.WorkflowID, req.Msg.JobMeta.Id, "failed to get output from proto"); err != nil {
			log.Error().Err(err).Msg("finishing job in db")
		}

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to finish job",
			Message: "Failed to get output from proto, please try again",
		})

		return nil, err
	}

	if err := h.queries.FinishJobRun(ctx, claims.WorkflowID, req.Msg.GetJobMeta().Id, conclusion); err != nil {
		log.Error().Err(err).Msg("finishing job in db")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to finish job",
			Message: "Failed to finish job in db, please try again",
		})

		return nil, err
	}

	if _, err := h.runnerClient.StopJob(ctx, connect.NewRequest(&pb.RunnerServiceStopJobRequest{
		WorkflowMeta: req.Msg.WorkflowMeta,
		JobMeta:      req.Msg.JobMeta,
	})); err != nil {
		log.Error().Err(err).Msg("calling runner client to stop job")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to finish job",
			Message: "Failed to stop job in runner",
		})

		if err := h.queries.UpdateWorkflowRun(ctx, workflowRun); err != nil {
			log.Error().Err(err).Msg("updating workflow run")
		}

		return nil, err
	}

	return connect.NewResponse(&pb.OrchestratorServiceFinishJobResponse{}), nil
}

func (h *Handler) CreateTask(ctx context.Context, req *connect.Request[pb.OrchestratorServiceCreateTaskRequest]) (*connect.Response[pb.OrchestratorServiceCreateTaskResponse], error) {
	defer utils.MeasureTime(time.Now(), "orchestrator service create task")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	workflowRun, err := h.queries.Unsafe_GetWorkflowRunByID(ctx, claims.WorkflowID)
	if err != nil {
		log.Error().Err(err).Msg("getting workflow run")
		return nil, err
	}

	task := typesDB.TaskRun{
		JobRunID:      req.Msg.JobMeta.Id,
		Name:          req.Msg.Data.Name,
		Status:        types.RunStatusRunning,
		WorkflowRunID: claims.WorkflowID,
		DockerImage:   types.Pointer(req.Msg.Data.GetDockerData().GetImage()),
	}

	if req.Msg.Skipped {
		task.Conclusion = types.Pointer(types.RunConclusionSkipped)
		task.Status = types.RunStatusCompleted
	}

	if err := h.queries.CreateTaskRun(ctx, &task); err != nil {
		log.Error().Err(err).Msg("creating task in db")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to create task",
			Message: "Failed to create task in db",
		})

		return nil, err
	}

	return connect.NewResponse(&pb.OrchestratorServiceCreateTaskResponse{
		Id: task.ID,
	}), nil

}

func (h *Handler) FinishTask(ctx context.Context, req *connect.Request[pb.OrchestratorServiceFinishTaskRequest]) (*connect.Response[pb.OrchestratorServiceFinishTaskResponse], error) {
	defer utils.MeasureTime(time.Now(), "orchestrator service finish task")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	workflowRun, err := h.queries.Unsafe_GetWorkflowRunByID(ctx, claims.WorkflowID)
	if err != nil {
		log.Error().Err(err).Msg("getting workflow run")
		return nil, err
	}

	conclusion, err := types.RunOutputFromProto(req.Msg.Conclusion)
	if err != nil {
		log.Error().Err(err).Msg("getting output from proto")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to finish task",
			Message: "Failed to get output from proto, please try again",
		})

		return nil, err
	}

	if err := h.queries.FinishTaskRun(ctx, claims.WorkflowID, req.Msg.TaskMeta.Id, conclusion); err != nil {
		log.Error().Err(err).Msg("updating task status in db")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to finish task",
			Message: "Failed to update task status in db, please try again",
		})

		return nil, err
	}

	return connect.NewResponse(&pb.OrchestratorServiceFinishTaskResponse{}), nil
}

func (h *Handler) CreateStep(ctx context.Context, req *connect.Request[pb.OrchestratorServiceCreateStepRequest]) (*connect.Response[pb.OrchestratorServiceCreateStepResponse], error) {
	defer utils.MeasureTime(time.Now(), "orchestrator service create step")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	workflowRun, err := h.queries.Unsafe_GetWorkflowRunByID(ctx, claims.WorkflowID)
	if err != nil {
		log.Error().Err(err).Msg("getting workflow run")
		return nil, err
	}

	step := typesDB.StepRun{
		Status:        types.RunStatusRunning,
		Type:          types.StepRunTypeExec,
		JobRunID:      req.Msg.JobMeta.Id,
		Name:          req.Msg.Name,
		WorkflowRunID: claims.WorkflowID,
	}

	if req.Msg.GetTaskMeta() != nil {
		step.TaskRunID = req.Msg.GetTaskMeta().Id
	}

	if err := h.queries.CreateStepRun(ctx, &step); err != nil {
		log.Error().Err(err).Msg("creating step in db")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to create step",
			Message: "Failed to create step in db",
		})

		return nil, err
	}

	// TODO - get env from env.go
	presignedWorkflowLogsURL, err := h.logStorageClient.PutObject(ctx, os.Getenv("WORKFLOW_LOGS_BUCKET"), fmt.Sprintf("%s/%s/%s/%s.csv", claims.OrgID, claims.ProjectID, claims.WorkflowID, step.ID), "text/csv", 60*60*24*7)
	if err != nil {
		log.Error().Err(err).Msg("getting step presigned url")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to create step",
			Message: "Failed to create presigned workflow logs url, please try again",
		})

		return nil, err
	}

	return connect.NewResponse(&pb.OrchestratorServiceCreateStepResponse{
		Id:                 step.ID,
		PresignedOutputUrl: presignedWorkflowLogsURL.URL,
	}), nil
}

func (h *Handler) FinishStep(ctx context.Context, req *connect.Request[pb.OrchestratorServiceFinishStepRequest]) (*connect.Response[pb.OrchestratorServiceFinishStepResponse], error) {
	defer utils.MeasureTime(time.Now(), "orchestrator service finish step")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	workflowRun, err := h.queries.Unsafe_GetWorkflowRunByID(ctx, claims.WorkflowID)
	if err != nil {
		log.Error().Err(err).Msg("getting workflow run")
		return nil, err
	}

	conclusion, err := types.RunOutputFromProto(req.Msg.Conclusion)
	if err != nil {
		log.Error().Err(err).Msg("getting output from proto")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to finish step",
			Message: "Failed to get output from proto, please try again",
		})

		return nil, err
	}

	if err := h.queries.FinishStepRun(ctx, claims.WorkflowID, req.Msg.Id, conclusion); err != nil {
		log.Error().Err(err).Msg("updating step status in db")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to finish step",
			Message: "Failed to update step status in db, please try again",
		})

		return nil, err
	}

	return connect.NewResponse(&pb.OrchestratorServiceFinishStepResponse{}), nil
}

func (h *Handler) FinishWorkflow(ctx context.Context, req *connect.Request[pb.OrchestratorServiceFinishWorkflowRequest]) (*connect.Response[pb.OrchestratorServiceFinishWorkflowResponse], error) {
	defer utils.MeasureTime(time.Now(), "orchestrator service finish workflow")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	workflowRun, err := h.queries.Unsafe_GetWorkflowRunByID(ctx, claims.WorkflowID)
	if err != nil {
		log.Error().Err(err).Msg("getting workflow run")
		return nil, err
	}

	defer func() {
		go func() {
			// This cleans up the runner jobs
			if _, err := h.runnerClient.CleanUpJobs(context.Background(), connect.NewRequest(&pb.RunnerServiceCleanUpJobsRequest{
				WorkflowJwt: req.Msg.WorkflowMeta.WorkflowJwt,
			})); err != nil {
				log.Error().Err(err).Msg("calling runner client to clean up jobs")

				h.addWorkflowRunAlert(context.Background(), workflowRun, types.WorkflowRunAlert{
					Type:    types.WorkflowRunAlertTypeError,
					Title:   "Failed to finish workflow",
					Message: "Failed to clean up jobs in runner",
				})
			}
		}()
	}()

	output, err := types.RunOutputFromProto(req.Msg.Conclusion)
	if err != nil {
		log.Error().Err(err).Msg("output from proto")

		if err := h.queries.FailWorkflowRun(context.Background(), workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to finish workflow",
			Message: "Failed to get output from proto, please try again",
		}); err != nil {
			log.Error().Err(err).Msg("failing workflow run")
		}

		return nil, err
	}

	workflowRun.Status = types.RunStatusCompleted
	workflowRun.Conclusion = &output
	workflowRun.FinishedAt = types.Pointer(utils.CurrentTime())

	workflowRun.BuildMinutes = int(math.Round(workflowRun.FinishedAt.Sub(workflowRun.CreatedAt).Minutes()+0.5)) * types.GetBuildMinutesScale(types.CloudRunner(workflowRun.Runner))

	if err := h.queries.UpdateWorkflowRun(context.Background(), workflowRun); err != nil {
		log.Error().Err(err).Msg("updating workflow status")

		if err := h.queries.FailWorkflowRun(context.Background(), workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to finish workflow",
			Message: "Failed to update workflow status, please try again",
		}); err != nil {
			log.Error().Err(err).Msg("failing workflow run")
		}

		return nil, err
	}

	if err := h.queries.FailStuckRuns(context.Background(), workflowRun.ID); err != nil {
		log.Error().Err(err).Msg("failing stuck runs")

		if err := h.queries.FailWorkflowRun(context.Background(), workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to finish workflow",
			Message: "Whilst finishing the workflow, we failed to fail any stuck runs",
		}); err != nil {
			log.Error().Err(err).Msg("failing workflow run")
		}
	}

	org, err := h.queries.Unsafe_GetOrgByID(context.Background(), claims.OrgID)
	if err != nil {
		log.Error().Err(err).Msg("getting org")
	} else {
		if err := h.chargeForMinutes(context.Background(), org); err != nil {
			log.Error().Err(err).Msg("charging for minutes")
		}
		if err := h.chargeForCommitters(context.Background(), org); err != nil {
			log.Error().Err(err).Msg("charging for committers")
		}
	}

	if err := h.git.UpdateRunStatusInRepo(context.Background(), *workflowRun); err != nil {
		log.Error().Err(err).Msg("updating run status in repo")

		h.addWorkflowRunAlert(context.Background(), workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to finish workflow",
			Message: "Failed to update status in git repo",
		})
	}

	return connect.NewResponse(&pb.OrchestratorServiceFinishWorkflowResponse{}), nil
}

func (h *Handler) WorkflowStarted(ctx context.Context, req *connect.Request[pb.OrchestratorServiceWorkflowStartedRequest]) (*connect.Response[pb.OrchestratorServiceWorkflowStartedResponse], error) {
	defer utils.MeasureTime(time.Now(), "orchestrator service workflow started")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	if err := h.queries.UpdateWorkflowRunStatus(ctx, claims.WorkflowID, types.RunStatusRunning); err != nil {
		log.Error().Err(err).Msg("updating workflow status")
		return nil, err
	}

	workflowRun, err := h.queries.Unsafe_GetWorkflowRunByID(ctx, claims.WorkflowID)
	if err != nil {
		log.Error().Err(err).Msg("getting workflow run")
		return nil, err
	}

	if err := h.git.UpdateRunStatusInRepo(ctx, *workflowRun); err != nil {
		log.Error().Err(err).Msg("updating run status in repo")

		if err := h.queries.FailWorkflowRun(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to start workflow",
			Message: "Failed to update run status in repo, please try again",
		}); err != nil {
			log.Error().Err(err).Msg("failing workflow run")
		}

		return nil, err
	}

	return connect.NewResponse(&pb.OrchestratorServiceWorkflowStartedResponse{}), nil
}

func (h *Handler) JobStarted(ctx context.Context, req *connect.Request[pb.OrchestratorServiceJobStartedRequest]) (*connect.Response[pb.OrchestratorServiceJobStartedResponse], error) {
	defer utils.MeasureTime(time.Now(), "orchestrator service job started")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	workflowRun, err := h.queries.Unsafe_GetWorkflowRunByID(ctx, claims.WorkflowID)
	if err != nil {
		log.Error().Err(err).Msg("getting workflow run")
		return nil, err
	}

	if err := h.queries.UpdateJobRunStatus(ctx, claims.WorkflowID, req.Msg.JobMeta.Id, types.RunStatusRunning); err != nil {
		log.Error().Err(err).Msg("updating job status")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to start job",
			Message: "Failed to update job status, please try again",
		})

		return nil, err
	}

	return connect.NewResponse(&pb.OrchestratorServiceJobStartedResponse{}), nil
}

func (h *Handler) CreateWorkflowAlert(ctx context.Context, req *connect.Request[pb.OrchestratorServiceCreateWorkflowAlertRequest]) (*connect.Response[pb.OrchestratorServiceCreateWorkflowAlertResponse], error) {

	defer utils.MeasureTime(time.Now(), "orchestrator service create workflow alert started")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	workflowRun, err := h.queries.Unsafe_GetWorkflowRunByID(ctx, claims.WorkflowID)
	if err != nil {
		log.Error().Err(err).Msg("getting workflow run")
		return nil, err
	}

	alertType := types.WorkflowRunAlertTypeError
	switch req.Msg.Alert.Type {
	case *pb.WorkflowAlert_TYPE_INFO.Enum():
		alertType = types.WorkflowRunAlertTypeInfo
	case pb.WorkflowAlert_TYPE_WARNING:
		alertType = types.WorkflowRunAlertTypeWarning
	}

	if err := workflowRun.AppendAlert(types.WorkflowRunAlert{
		Type:    alertType,
		Message: req.Msg.Alert.Message,
		Title:   req.Msg.Alert.Title,
	}); err != nil {
		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to create alert",
			Message: err.Error(), // This error shouldn't contain any sensitive information
		})

		return nil, err
	}

	if err := h.queries.UpdateWorkflowRun(ctx, workflowRun); err != nil {
		log.Error().Err(err).Msg("updating workflow run")

		h.addWorkflowRunAlert(ctx, workflowRun, types.WorkflowRunAlert{
			Type:    types.WorkflowRunAlertTypeError,
			Title:   "Failed to create alert",
			Message: "Failed to update workflow run",
		})

		return nil, err
	}

	return connect.NewResponse(&pb.OrchestratorServiceCreateWorkflowAlertResponse{}), nil
}
