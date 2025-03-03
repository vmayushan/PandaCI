package grpcRunner

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	grpcMiddleware "github.com/pandaci-com/pandaci/app/grpc/middleware"
	"github.com/pandaci-com/pandaci/pkg/utils"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
	"github.com/pandaci-com/pandaci/types"
	"github.com/rs/zerolog/log"
)

var jobStartedSources = make(map[string]chan string)

func SubscribeToJobStarted(ctx context.Context, workflowID string, jobID string) <-chan string {
	key := fmt.Sprintf("%s-%s", workflowID, jobID)

	ch := make(chan string)
	jobStartedSources[key] = ch

	go func() {
		<-ctx.Done()
		UnsubscribeFromJobStarted(ctx, workflowID, jobID)
	}()

	return ch
}

func EmitJobStarted(ctx context.Context, workflowID string, jobID string, res string) {
	key := fmt.Sprintf("%s-%s", workflowID, jobID)
	if ch, ok := jobStartedSources[key]; ok {
		ch <- res
	} else {
		log.Warn().Msgf("no subscribers for job started: %s", jobID)
	}
}

func UnsubscribeFromJobStarted(ctx context.Context, workflowID string, jobID string) {
	key := fmt.Sprintf("%s-%s", workflowID, jobID)
	if ch, ok := jobStartedSources[key]; ok {
		close(ch)
		delete(jobStartedSources, key)
	}
}

func (h *Handler) StartWorkflow(ctx context.Context, req *connect.Request[pb.RunnerServiceStartWorkflowRequest]) (*connect.Response[pb.RunnerServiceStartWorkflowResponse], error) {
	defer utils.MeasureTime(time.Now(), "runner service start workflow")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	res, err := h.runner.StartWorkflow(ctx, claims.WorkflowID, req.Msg)
	if err != nil {
		log.Error().Err(err).Msgf("starting workflow")
		return nil, err
	}

	return connect.NewResponse(res), nil
}

func (h *Handler) CleanUpJobs(ctx context.Context, req *connect.Request[pb.RunnerServiceCleanUpJobsRequest]) (*connect.Response[pb.RunnerServiceCleanUpJobsResponse], error) {
	defer utils.MeasureTime(time.Now(), "runner service stop workflow")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	res, err := h.runner.CleanUpJobs(ctx, claims.WorkflowID, req.Msg)
	if err != nil {
		log.Error().Err(err).Msg("runner failed to stop workflow")
		return nil, err
	}

	return connect.NewResponse(res), nil
}

func (h *Handler) StartJob(ctx context.Context, req *connect.Request[pb.RunnerServiceStartJobRequest]) (*connect.Response[pb.RunnerServiceStartJobResponse], error) {
	defer utils.MeasureTime(time.Now(), "runner service start job")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	sub := SubscribeToJobStarted(ctx, claims.WorkflowID, req.Msg.Id)

	res, err := h.runner.StartJob(ctx, claims.WorkflowID, req.Msg)
	if err != nil {
		log.Error().Err(err).Msg("runner failed to start job")
		h.addWorkflowRunAlert(ctx, req.Msg.WorkflowMeta, pb.WorkflowAlert_TYPE_ERROR, "Failed to start job", "something went wrong")
		return nil, err
	}

	select {
	case <-ctx.Done():
		{
			h.addWorkflowRunAlert(context.Background(), req.Msg.WorkflowMeta, pb.WorkflowAlert_TYPE_ERROR, "Failed to start job", "context cancelled")
			if _, err := h.runner.StopJob(ctx, claims.WorkflowID, &pb.RunnerServiceStopJobRequest{
				WorkflowMeta: req.Msg.WorkflowMeta,
				JobMeta:      res.JobMeta,
			}); err != nil {
				log.Error().Err(err).Msg("runner failed to stop job")
			}

			return nil, fmt.Errorf("context cancelled")
		}
	case address := <-sub:
		{
			// Our fly runner address is already correct
			// TODO - we should try to remove this and always have the corect address sent
			if h.runner.GetRunnerType() == types.RunnerTypeLocal {
				res.JobMeta.Address = address
			}
		}
	}

	log.Info().Msgf("Job started: %s", res.JobMeta.Address)

	if _, err := h.orchestratorClient.JobStarted(ctx, connect.NewRequest(&pb.OrchestratorServiceJobStartedRequest{
		WorkflowMeta: req.Msg.WorkflowMeta,
		JobMeta:      res.JobMeta,
	})); err != nil {
		log.Error().Err(err).Msg("error sending job started request to orchestrator")
		h.addWorkflowRunAlert(ctx, req.Msg.WorkflowMeta, pb.WorkflowAlert_TYPE_ERROR, "Failed to start job", "something went wrong")

		if _, err := h.runner.StopJob(ctx, claims.WorkflowID, &pb.RunnerServiceStopJobRequest{
			WorkflowMeta: req.Msg.WorkflowMeta,
			JobMeta:      res.JobMeta,
		}); err != nil {
			log.Error().Err(err).Msg("runner failed to stop job")
		}

		return nil, err
	}

	return connect.NewResponse(res), nil
}

func (h *Handler) StopJob(ctx context.Context, req *connect.Request[pb.RunnerServiceStopJobRequest]) (*connect.Response[pb.RunnerServiceStopJobResponse], error) {
	defer utils.MeasureTime(time.Now(), "runner service stop job")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	res, err := h.runner.StopJob(ctx, claims.WorkflowID, req.Msg)
	if err != nil {
		log.Error().Err(err).Msg("runner failed to stop job")
		h.addWorkflowRunAlert(ctx, req.Msg.WorkflowMeta, pb.WorkflowAlert_TYPE_ERROR, "Failed to stop job", "something went wrong")
		return nil, err
	}

	return connect.NewResponse(res), nil
}

func (h *Handler) JobStarted(ctx context.Context, req *connect.Request[pb.RunnerServiceJobStartedRequest]) (*connect.Response[pb.RunnerServiceJobStartedResponse], error) {
	defer utils.MeasureTime(time.Now(), "runner service job started")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	EmitJobStarted(ctx, claims.WorkflowID, req.Msg.Id, req.Msg.Address)

	return connect.NewResponse(&pb.RunnerServiceJobStartedResponse{}), nil
}

func (h *Handler) GetLogStream(ctx context.Context, req *connect.Request[pb.RunnerServiceGetLogStreamRequest]) (*connect.Response[pb.RunnerServiceGetLogStreamResponse], error) {
	defer utils.MeasureTime(time.Now(), "runner service get log stream")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	res, err := h.runner.GetLogStream(ctx, claims.WorkflowID, req.Msg)
	if err != nil {
		log.Error().Err(err).Msg("runner failed to get log stream")
		return nil, err
	}

	return connect.NewResponse(res), nil
}
