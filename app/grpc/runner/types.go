package grpcRunner

import (
	"context"

	"connectrpc.com/connect"
	"github.com/alfiejones/panda-ci/app/runner"
	"github.com/alfiejones/panda-ci/pkg/jwt"
	"github.com/rs/zerolog/log"

	pb "github.com/alfiejones/panda-ci/proto/go/v1"
	pbConnect "github.com/alfiejones/panda-ci/proto/go/v1/v1connect"
)

type (
	Handler struct {
		runner             runner.Handler
		jwt                jwt.JWTHandler
		orchestratorClient pbConnect.OrchestratorServiceClient
	}
)

func (h *Handler) addWorkflowRunAlert(ctx context.Context, workflowMeta *pb.WorkflowMeta, alertType pb.WorkflowAlert_Type, title string, message string) {
	if _, err := h.orchestratorClient.CreateWorkflowAlert(ctx, connect.NewRequest(&pb.OrchestratorServiceCreateWorkflowAlertRequest{
		WorkflowMeta: workflowMeta,
		Alert: &pb.WorkflowAlert{
			Type:    alertType,
			Title:   title,
			Message: message,
		},
	})); err != nil {
		log.Error().Err(err).Msg("creating workflow alert")
	}
}

func NewHandler(jwt jwt.JWTHandler, orchestratorClient pbConnect.OrchestratorServiceClient, runner runner.Handler) (*Handler, error) {
	return &Handler{
		jwt:                jwt,
		orchestratorClient: orchestratorClient,
		runner:             runner,
	}, nil
}
