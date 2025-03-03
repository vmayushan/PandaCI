package grpcOrchestrator

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/pandaci-com/pandaci/app/git"
	grpcMiddleware "github.com/pandaci-com/pandaci/app/grpc/middleware"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/pkg/jwt"
	"github.com/pandaci-com/pandaci/pkg/retryClient"
	"github.com/pandaci-com/pandaci/pkg/utils/env"
	"github.com/pandaci-com/pandaci/platform/storage"
	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
	"github.com/rs/zerolog/log"

	pbConnect "github.com/pandaci-com/pandaci/proto/go/v1/v1connect"
)

type (
	Handler struct {
		queries          *queries.Queries
		jwt              jwt.JWTHandler
		logStorageClient *storage.BucketClient
		git              *git.Handler
		paddleClient     *paddle.SDK
		runnerClient     pbConnect.RunnerServiceClient // TODO - this will need to change since it's not very exensible,
		// Instead, we need a function which can take a runner name. eg. ubuntu and
		// return a RunnerServiceClient we can use. This way we could connect to self hosted runners
		// For now we're just focusing on cloud so we don't need to worry
	}
)

func (h *Handler) addWorkflowRunAlert(ctx context.Context, workflowRun *typesDB.WorkflowRun, alert types.WorkflowRunAlert) {
	if err := workflowRun.AppendAlert(alert); err != nil {
		log.Error().Err(err).Msg("setting alerts")
	}

	if err := h.queries.UpdateWorkflowRun(ctx, workflowRun); err != nil {
		log.Error().Err(err).Msg("updating workflow run")
	}
}

func NewHandler(queries *queries.Queries, jwt jwt.JWTHandler, logStorageClient *storage.BucketClient, git *git.Handler) (*Handler, error) {

	runnerAddress, err := env.GetRunnerAddress()
	if err != nil {
		return nil, err
	}

	runnerClient := pbConnect.NewRunnerServiceClient(&http.Client{
		Transport: &retryClient.RetryRoundTripper{
			Base: http.DefaultTransport,
		},
	}, *runnerAddress, connect.WithInterceptors(grpcMiddleware.NewWorkflowJWTInerceptor(jwt, nil)))

	apiKey, err := env.GetPaddleAPIKey()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get paddle api key")
		return nil, err
	}

	baseURL := paddle.SandboxBaseURL
	if env.GetStage() == "prod" {
		baseURL = paddle.ProductionBaseURL
	}

	client, err := paddle.New(
		*apiKey,
		paddle.WithBaseURL(baseURL),
	)

	return &Handler{
		queries:          queries,
		jwt:              jwt,
		runnerClient:     runnerClient,
		logStorageClient: logStorageClient,
		git:              git,
		paddleClient:     client,
	}, nil
}
