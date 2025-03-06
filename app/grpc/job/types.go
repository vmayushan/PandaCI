package grpcJob

import (
	"context"

	"connectrpc.com/connect"
	"github.com/pandaci-com/pandaci/app/grpc/job/service"
	"github.com/pandaci-com/pandaci/pkg/jwt"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
	pbConnect "github.com/pandaci-com/pandaci/proto/go/v1/v1connect"
	"github.com/rs/zerolog/log"
)

type (
	Handler struct {
		jwt                jwt.JWTHandler
		service            *service.Handler
		address            string
		workflowMeta       *pb.WorkflowMeta
		orchestratorClient pbConnect.OrchestratorServiceClient
	}
)

func (h *Handler) addWorkflowRunAlert(ctx context.Context, alertType pb.WorkflowAlert_Type, title string, message string) {
	if _, err := h.orchestratorClient.CreateWorkflowAlert(ctx, connect.NewRequest(&pb.OrchestratorServiceCreateWorkflowAlertRequest{
		WorkflowMeta: h.workflowMeta,
		Alert: &pb.WorkflowAlert{
			Type:    alertType,
			Title:   title,
			Message: message,
		},
	})); err != nil {
		log.Error().Err(err).Msg("creating workflow alert")
	}
}

func NewHandler(orchestratorClient pbConnect.OrchestratorServiceClient, workflowMeta *pb.WorkflowMeta, jwt jwt.JWTHandler, address string) (*Handler, error) {

	service, err := service.NewHandler()
	if err != nil {
		return nil, err
	}

	return &Handler{
		jwt:                jwt,
		service:            service,
		address:            address,
		orchestratorClient: orchestratorClient,
		workflowMeta:       workflowMeta,
	}, nil
}
