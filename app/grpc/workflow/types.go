package grpcWorkflow

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	grpcMiddleware "github.com/alfiejones/panda-ci/app/grpc/middleware"
	"github.com/alfiejones/panda-ci/pkg/flyClient"
	"github.com/alfiejones/panda-ci/pkg/jwt"
	"github.com/alfiejones/panda-ci/pkg/stream"
	"github.com/rs/zerolog/log"

	pb "github.com/alfiejones/panda-ci/proto/go/v1"
	pbConnect "github.com/alfiejones/panda-ci/proto/go/v1/v1connect"
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

type (
	Handler struct {
		orchestratorClient   pbConnect.OrchestratorServiceClient
		getJobClient         func(jobMeta *pb.JobMeta) pbConnect.JobServiceClient
		workflowMeta         *pb.WorkflowMeta
		worstJobConclusion   map[string]pb.Conclusion
		stepLogs             map[string]*stream.LogStream
		jobMonitorContextMap map[string]*context.CancelFunc
		eventStream          *stream.EventStream
	}
)

func getGetJobClient(jwt jwt.JWTHandler) func(jobMeta *pb.JobMeta) pbConnect.JobServiceClient {
	return func(jobMeta *pb.JobMeta) pbConnect.JobServiceClient {

		client := &http.Client{
			Transport: &flyClient.FlyRoundTripper{
				Base:    http.DefaultTransport,
				AppName: &jobMeta.GetFlyMeta().AppName,
			},
		}

		return pbConnect.NewJobServiceClient(client, jobMeta.Address+"/grpc", connect.WithInterceptors(grpcMiddleware.NewWorkflowJWTInerceptor(jwt, nil)))
	}
}

func NewHandler(jwt jwt.JWTHandler, orchestratorClient pbConnect.OrchestratorServiceClient, workflowMeta *pb.WorkflowMeta, worstJobConclusion map[string]pb.Conclusion, stepLogs map[string]*stream.LogStream, eventStream *stream.EventStream) *Handler {
	return &Handler{
		orchestratorClient:   orchestratorClient,
		getJobClient:         getGetJobClient(jwt),
		workflowMeta:         workflowMeta,
		worstJobConclusion:   worstJobConclusion,
		stepLogs:             stepLogs,
		eventStream:          eventStream,
		jobMonitorContextMap: make(map[string]*context.CancelFunc),
	}
}
