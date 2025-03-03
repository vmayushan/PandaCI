package grpcWorkflow

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/log"

	logStream "github.com/pandaci-com/pandaci/pkg/stream"
	"github.com/pandaci-com/pandaci/pkg/uploads"
	"github.com/pandaci-com/pandaci/pkg/utils"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
	"github.com/pandaci-com/pandaci/types"
)

func (h *Handler) StartJob(ctx context.Context, req *connect.Request[pb.WorkflowServiceStartJobRequest]) (*connect.Response[pb.WorkflowServiceStartJobResponse], error) {
	log.Info().Msg("Starting job")
	defer utils.MeasureTime(time.Now(), "workflow service start job")

	createdJob, err := h.orchestratorClient.CreateJob(ctx, connect.NewRequest(&pb.OrchestratorServiceCreateJobRequest{
		WorkflowMeta: h.workflowMeta,
		Name:         req.Msg.Name,
		Runner:       req.Msg.Runner,
		Skipped:      req.Msg.Skipped,
	}))
	if err != nil {
		log.Error().Err(err).Msg("creating job")
		return nil, err
	}

	if !req.Msg.Skipped {
		go h.monitorJobStatus(createdJob.Msg.JobMeta, h.workflowMeta)
	}

	if err := h.eventStream.UnknownEvent(); err != nil {
		log.Error().Err(err).Msg("writing unknown event")
	}

	res := &pb.WorkflowServiceStartJobResponse{
		JobMeta: createdJob.Msg.JobMeta,
	}

	h.worstJobConclusion[createdJob.Msg.JobMeta.Id] = pb.Conclusion_CONCLUSION_SUCCESS

	return connect.NewResponse(res), nil
}

func (h *Handler) StopJob(ctx context.Context, req *connect.Request[pb.WorkflowServiceStopJobRequest]) (*connect.Response[pb.WorkflowServiceStopJobResponse], error) {
	log.Info().Msg("Stopping job")
	defer utils.MeasureTime(time.Now(), "workflow service stop job")

	if req.Msg.Conclusion == pb.Conclusion_CONCLUSION_UNSPECIFIED {
		req.Msg.Conclusion = h.worstJobConclusion[req.Msg.JobMeta.Id]
	}

	h.worstJobConclusion[req.Msg.JobMeta.Id] = req.Msg.Conclusion

	if cancel, ok := h.jobMonitorContextMap[req.Msg.JobMeta.Id]; ok && cancel != nil {
		// Stop pinging the job
		(*cancel)()
	}

	if _, err := h.orchestratorClient.FinishJob(ctx, connect.NewRequest(&pb.OrchestratorServiceFinishJobRequest{
		WorkflowMeta: h.workflowMeta,
		JobMeta:      req.Msg.JobMeta,
		Conclusion:   req.Msg.Conclusion,
	})); err != nil {
		log.Error().Err(err).Msg("Error sending finish request to orchestrator")
		return nil, err
	}

	if err := h.eventStream.UnknownEvent(); err != nil {
		log.Error().Err(err).Msg("writing unknown event")
	}

	return connect.NewResponse(&pb.WorkflowServiceStopJobResponse{
		Conclusion: req.Msg.Conclusion,
	}), nil
}

func (h *Handler) StartTask(ctx context.Context, req *connect.Request[pb.WorkflowServiceStartTaskRequest]) (*connect.Response[pb.WorkflowServiceStartTaskResponse], error) {
	log.Info().Msg("Starting task")
	defer utils.MeasureTime(time.Now(), "workflow service start task")

	createdTask, err := h.orchestratorClient.CreateTask(ctx, connect.NewRequest(&pb.OrchestratorServiceCreateTaskRequest{
		WorkflowMeta: h.workflowMeta,
		JobMeta:      req.Msg.JobMeta,
		Data:         req.Msg.Data,
		Skipped:      req.Msg.Skipped,
	}))
	if err != nil {
		log.Error().Err(err).Msg("creating task")
		return nil, err
	}

	if err := h.eventStream.UnknownEvent(); err != nil {
		log.Error().Err(err).Msg("writing unknown event")
	}

	if req.Msg.Skipped {
		return connect.NewResponse(&pb.WorkflowServiceStartTaskResponse{
			TaskMeta: &pb.TaskMeta{Id: createdTask.Msg.Id,
				Name: req.Msg.Data.Name,
			},
		}), nil
	}

	runningTask, err := h.getJobClient(req.Msg.JobMeta).StartTask(ctx, connect.NewRequest(&pb.JobServiceStartTaskRequest{
		WorkflowMeta: h.workflowMeta,
		JobMeta:      req.Msg.JobMeta,
		Data:         req.Msg.Data,
		Id:           createdTask.Msg.Id,
	}))
	if err != nil {
		log.Error().Err(err).Msg("starting task")
		h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Failed to start task", err.Error())
		return nil, err
	}

	res := &pb.WorkflowServiceStartTaskResponse{
		TaskMeta: runningTask.Msg.TaskMeta,
	}

	return connect.NewResponse(res), nil
}

func (h *Handler) StopTask(ctx context.Context, req *connect.Request[pb.WorkflowServiceStopTaskRequest]) (*connect.Response[pb.WorkflowServiceStopTaskResponse], error) {
	log.Info().Msg("Stopping task")
	defer utils.MeasureTime(time.Now(), "workflow service stop task")

	if types.CompareProtoConclusionRank(req.Msg.Conclusion, h.worstJobConclusion[req.Msg.JobMeta.Id]) {
		h.worstJobConclusion[req.Msg.JobMeta.Id] = req.Msg.Conclusion
	}

	// TODO - we can probably do these in parallel
	if _, err := h.orchestratorClient.FinishTask(ctx, connect.NewRequest(&pb.OrchestratorServiceFinishTaskRequest{
		WorkflowMeta: h.workflowMeta,
		JobMeta:      req.Msg.JobMeta,
		TaskMeta:     req.Msg.TaskMeta,
		Conclusion:   req.Msg.Conclusion,
	})); err != nil {
		log.Error().Err(err).Msg("Error sending finish request to orchestrator")
		return nil, err
	}

	if err := h.eventStream.UnknownEvent(); err != nil {
		log.Error().Err(err).Msg("writing unknown event")
	}

	if _, err := h.getJobClient(req.Msg.JobMeta).StopTask(ctx, connect.NewRequest(&pb.JobServiceStopTaskRequest{
		WorkflowMeta: h.workflowMeta,
		JobMeta:      req.Msg.JobMeta,
		TaskMeta:     req.Msg.TaskMeta,
	})); err != nil {
		log.Error().Err(err).Msg("Error stopping task")

		h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Failed to stop task", err.Error())

		return nil, err
	}

	return connect.NewResponse(&pb.WorkflowServiceStopTaskResponse{}), nil
}

func (h *Handler) StartStep(ctx context.Context, req *connect.Request[pb.WorkflowServiceStartStepRequest], stream *connect.ServerStream[pb.WorkflowServiceStartStepResponse]) error {
	log.Info().Msg("Starting step")
	defer utils.MeasureTime(time.Now(), "workflow service start step")

	stepName := fmt.Sprintf("$ %s", req.Msg.GetExecData().GetCmd())
	if req.Msg.GetExecData().Cwd != "" {
		stepName = fmt.Sprintf("%s %s", req.Msg.GetExecData().Cwd, stepName)
	}

	stepRun, err := h.orchestratorClient.CreateStep(ctx, connect.NewRequest(&pb.OrchestratorServiceCreateStepRequest{
		WorkflowMeta: h.workflowMeta,
		JobMeta:      req.Msg.JobMeta,
		TaskMeta:     req.Msg.TaskMeta,
		Name:         stepName,
		StepData: &pb.OrchestratorServiceCreateStepRequest_ExecData{
			ExecData: req.Msg.GetExecData(),
		},
	}))
	if err != nil {
		log.Error().Err(err).Msg("creating step")
		return err
	}

	if err := h.eventStream.UnknownEvent(); err != nil {
		log.Error().Err(err).Msg("writing unknown event")
	}

	h.stepLogs[stepRun.Msg.Id] = logStream.NewLogStream([]string{"timestamp,type,data\n"})

	defer func() {
		defer delete(h.stepLogs, stepRun.Msg.Id)

		// upload logs to the object store
		var buffer bytes.Buffer

		for _, entry := range h.stepLogs[stepRun.Msg.Id].Entries() {
			if _, err := buffer.WriteString(entry); err != nil {
				log.Error().Err(err).Msg("writing log")
				h.addWorkflowRunAlert(context.Background(), pb.WorkflowAlert_TYPE_ERROR, "Failed to write log", err.Error())
				break
			}
		}

		if err := uploads.UploadFile(context.Background(), stepRun.Msg.PresignedOutputUrl, &buffer, "text/csv"); err != nil {
			log.Error().Err(err).Msg("uploading logs")
			h.addWorkflowRunAlert(context.Background(), pb.WorkflowAlert_TYPE_ERROR, "Failed to upload logs", err.Error())
		}
	}()

	stepStream, err := h.getJobClient(req.Msg.JobMeta).StartStep(ctx, connect.NewRequest(&pb.JobServiceStartStepRequest{
		WorkflowJwt:        req.Msg.WorkflowJwt,
		JobMeta:            req.Msg.JobMeta,
		TaskMeta:           req.Msg.TaskMeta,
		Id:                 stepRun.Msg.Id,
		PresignedOutputUrl: stepRun.Msg.PresignedOutputUrl,
		Data: &pb.JobServiceStartStepRequest_ExecData{
			ExecData: req.Msg.GetExecData(),
		},
	}))
	if err != nil {
		log.Error().Err(err).Msg("starting step")
		h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Failed to start step", err.Error())
		return err
	}

	if err := h.eventStream.UnknownEvent(); err != nil {
		log.Error().Err(err).Msg("writing unknown event")
	}

	conclusion := pb.Conclusion_CONCLUSION_SUCCESS

	defer stepStream.Close()

	for stepStream.Receive() {

		execData := stepStream.Msg().GetExec()

		if execData.GetExitCode() != 0 {
			conclusion = pb.Conclusion_CONCLUSION_FAILURE
		}

		// Write the logs to the log stream
		// This will then send the logs to any clients that are listening
		// e.g. the UI

		var logType logStream.LogType
		var data string

		if _, ok := execData.LogData.(*pb.StepExecPayload_Stdout); ok {
			logType = "stdout"
			data = string(execData.GetStdout())
		} else if _, ok := execData.LogData.(*pb.StepExecPayload_Stderr); ok {
			logType = "stderr"
			data = string(execData.GetStderr())
		} else {
			logType = "exit"
			data = strconv.Itoa(int(execData.GetExitCode()))
		}

		if err := h.stepLogs[stepRun.Msg.Id].WriteLog(logStream.Log{
			Timestamp: execData.Timestamp.AsTime(),
			Type:      logType,
			Data:      data,
		}); err != nil {
			log.Error().Err(err).Msg("writing log")
			h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Failed to write log", err.Error())
			return err
		}

		if err := stream.Send(&pb.WorkflowServiceStartStepResponse{
			Payload: &pb.WorkflowServiceStartStepResponse_Exec{
				Exec: execData,
			},
		}); err != nil {
			log.Error().Err(err).Msg("sending step data")
			h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Failed to send step data", err.Error())
			return err
		}
	}

	if err := stepStream.Err(); err != nil {
		log.Error().Err(err).Msg("stream error")
		h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Stream error", err.Error())
	}

	if _, err := h.orchestratorClient.FinishStep(ctx, connect.NewRequest(&pb.OrchestratorServiceFinishStepRequest{
		WorkflowMeta: h.workflowMeta,
		JobMeta:      req.Msg.JobMeta,
		TaskMeta:     req.Msg.TaskMeta,
		Id:           stepRun.Msg.Id,
		Conclusion:   conclusion,
	})); err != nil {
		log.Error().Err(err).Msg("finishing step")
		return err
	}

	return nil
}

func (h *Handler) StopStep(ctx context.Context, req *connect.Request[pb.WorkflowServiceStopStepRequest]) (*connect.Response[pb.WorkflowServiceStopStepResponse], error) {
	log.Info().Msg("Stopping step")
	defer utils.MeasureTime(time.Now(), "workflow service stop step")
	return nil, nil
}

func (h *Handler) CreateJobVolume(ctx context.Context, req *connect.Request[pb.WorkflowServiceCreateJobVolumeRequest]) (*connect.Response[pb.WorkflowServiceCreateJobVolumeResponse], error) {
	log.Info().Msg("Creating job volume")
	defer utils.MeasureTime(time.Now(), "workflow service create job volume")

	res, err := h.getJobClient(req.Msg.JobMeta).CreateJobVolume(ctx, connect.NewRequest(&pb.JobServiceCreateJobVolumeRequest{
		WorkflowMeta: h.workflowMeta,
		JobMeta:      req.Msg.JobMeta,
		Host:         req.Msg.Host,
	}))
	if err != nil {
		log.Error().Err(err).Msg("creating job volume")
		h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Failed to create job volume", err.Error())
		return nil, err
	}

	return connect.NewResponse(
		&pb.WorkflowServiceCreateJobVolumeResponse{
			Source: res.Msg.Source,
		},
	), nil
}

func (h *Handler) CreateWorkflowAlert(ctx context.Context, req *connect.Request[pb.WorkflowServiceCreateWorkflowAlertRequest]) (*connect.Response[pb.WorkflowServiceCreateWorkflowAlertResponse], error) {
	log.Info().Msg("Creating workflow alert")
	defer utils.MeasureTime(time.Now(), "workflow service create workflow alert")

	if _, err := h.orchestratorClient.CreateWorkflowAlert(ctx, connect.NewRequest(&pb.OrchestratorServiceCreateWorkflowAlertRequest{
		WorkflowMeta: h.workflowMeta,
		Alert:        req.Msg.Alert,
	})); err != nil {
		log.Error().Err(err).Msg("creating workflow alert")
		return nil, err
	}

	return connect.NewResponse(&pb.WorkflowServiceCreateWorkflowAlertResponse{}), nil
}

func (h *Handler) Ping(ctx context.Context, req *connect.Request[pb.WorkflowServicePingRequest]) (*connect.Response[pb.WorkflowServicePingResponse], error) {
	defer utils.MeasureTime(time.Now(), "workflow service ping")

	return connect.NewResponse(&pb.WorkflowServicePingResponse{}), nil
}
