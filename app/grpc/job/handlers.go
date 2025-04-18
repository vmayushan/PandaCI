package grpcJob

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	grpcMiddleware "github.com/pandaci-com/pandaci/app/grpc/middleware"
	"github.com/pandaci-com/pandaci/pkg/utils"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
	"github.com/rs/zerolog/log"
)

func (h *Handler) StartTask(ctx context.Context, req *connect.Request[pb.JobServiceStartTaskRequest]) (*connect.Response[pb.JobServiceStartTaskResponse], error) {
	defer utils.MeasureTime(time.Now(), "job service start task")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	if req.Msg.Data.GetDockerData() == nil {
		return nil, fmt.Errorf("docker data is required")
	}

	res, err := h.service.StarDockerTask(ctx, claims.WorkflowID, req.Msg)
	if err != nil {
		log.Error().Err(err).Msg("starting docker task")
		h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Failed to start task", err.Error())
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(res), nil
}

func (h *Handler) StopTask(ctx context.Context, req *connect.Request[pb.JobServiceStopTaskRequest]) (*connect.Response[pb.JobServiceStopTaskResponse], error) {
	defer utils.MeasureTime(time.Now(), "job service stop task")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	res, err := h.service.StopDockerTask(ctx, claims.WorkflowID, req.Msg)
	if err != nil {
		log.Error().Err(err).Msg("stopping docker task")
		h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Failed to stop task", err.Error())
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(res), nil
}

func (h *Handler) CreateJobVolume(ctx context.Context, req *connect.Request[pb.JobServiceCreateJobVolumeRequest]) (*connect.Response[pb.JobServiceCreateJobVolumeResponse], error) {
	defer utils.MeasureTime(time.Now(), "job service create job volume")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	res, err := h.service.CreateJobVolume(ctx, claims.WorkflowID, req.Msg)
	if err != nil {
		log.Error().Err(err).Msg("creating job volume")
		h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Failed to create job volume", err.Error())
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(res), nil
}

func (h *Handler) StartStep(ctx context.Context, req *connect.Request[pb.JobServiceStartStepRequest], stream *connect.ServerStream[pb.JobServiceStartStepResponse]) error {
	defer utils.MeasureTime(time.Now(), "job service start step")

	claims := grpcMiddleware.GetWorkflowClaims(ctx)

	if err := h.service.StartDockerStep(ctx, claims.WorkflowID, req.Msg, func(msg *pb.JobServiceStartStepResponse) error {
		if err := stream.Send(msg); err != nil {
			log.Error().Err(err).Msg("sending step message")
			h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Failed to send step message", err.Error())
			return err
		}
		return nil
	}); err != nil {
		log.Error().Err(err).Msg("starting step")
		h.addWorkflowRunAlert(ctx, pb.WorkflowAlert_TYPE_ERROR, "Failed to start step", err.Error())
		return err
	}

	return nil
}

func (h *Handler) StopStep(ctx context.Context, req *connect.Request[pb.JobServiceStopStepRequest]) (*connect.Response[pb.JobServiceStopStepResponse], error) {
	defer utils.MeasureTime(time.Now(), "job service stop step")

	return nil, nil
}

func (h *Handler) Ping(ctx context.Context, req *connect.Request[pb.JobServicePingRequest]) (*connect.Response[pb.JobServicePingResponse], error) {
	defer utils.MeasureTime(time.Now(), "job service ping")

	return connect.NewResponse(&pb.JobServicePingResponse{}), nil
}
