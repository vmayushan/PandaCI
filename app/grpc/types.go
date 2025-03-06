package grpc

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/pandaci-com/pandaci/app/git"
	grpcJob "github.com/pandaci-com/pandaci/app/grpc/job"
	grpcMiddleware "github.com/pandaci-com/pandaci/app/grpc/middleware"
	grpcOrchestrator "github.com/pandaci-com/pandaci/app/grpc/orchestrator"
	grpcRunner "github.com/pandaci-com/pandaci/app/grpc/runner"
	grpcWorkflow "github.com/pandaci-com/pandaci/app/grpc/workflow"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/app/runner"
	"github.com/pandaci-com/pandaci/pkg/jwt"
	"github.com/pandaci-com/pandaci/pkg/stream"
	"github.com/pandaci-com/pandaci/platform/storage"
	pbConnect "github.com/pandaci-com/pandaci/proto/go/v1/v1connect"

	pb "github.com/pandaci-com/pandaci/proto/go/v1"
)

type GRPCServer struct {
	server *http.ServeMux
}

func (s *GRPCServer) Start(address string) error {
	return http.ListenAndServe(address, h2c.NewHandler(s.server, &http2.Server{}))
}

func RegisterOrchestratorGRPC(e *echo.Echo, queries *queries.Queries, jwt jwt.JWTHandler, bucketClient *storage.BucketClient, git *git.Handler) error {

	orchestratorService, err := grpcOrchestrator.NewHandler(queries, jwt, bucketClient, git)
	if err != nil {
		return err
	}

	grpcMiddlewareHandler := connect.WithInterceptors(grpcMiddleware.NewWorkflowJWTInerceptor(jwt, nil), grpcMiddleware.NewSanitiseErrorsInerceptor())

	orchestratorPath, orchestratorHandler := pbConnect.NewOrchestratorServiceHandler(orchestratorService, grpcMiddlewareHandler)

	grpcGroup := e.Group("/grpc" + orchestratorPath)
	grpcGroup.Any("*", echo.WrapHandler(http.StripPrefix("/grpc", orchestratorHandler)))

	return nil
}

func RegisterWorkflowRPC(e *echo.Echo, jwt jwt.JWTHandler, orchestratorClient pbConnect.OrchestratorServiceClient, workflowMeta *pb.WorkflowMeta, worstJobConclusion map[string]pb.Conclusion, stepLogs map[string]*stream.LogStream, eventStream *stream.EventStream) error {

	grpcMiddlewareHandler := connect.WithInterceptors(grpcMiddleware.NewWorkflowJWTInerceptor(jwt, &workflowMeta.WorkflowJwt), grpcMiddleware.NewSanitiseErrorsInerceptor())

	workflowService := grpcWorkflow.NewHandler(jwt, orchestratorClient, workflowMeta, worstJobConclusion, stepLogs, eventStream)
	workflowPath, workflowHandler := pbConnect.NewWorkflowServiceHandler(workflowService, grpcMiddlewareHandler)

	grpcGroup := e.Group("/grpc" + workflowPath)
	grpcGroup.Any("*", echo.WrapHandler(http.StripPrefix("/grpc", workflowHandler)))

	return nil
}

func NewRunnerServer(jwt jwt.JWTHandler, orchestratorClient pbConnect.OrchestratorServiceClient, runner runner.Handler) (*GRPCServer, error) {
	grpc := http.NewServeMux()

	grpcMiddlewareHandler := connect.WithInterceptors(grpcMiddleware.NewWorkflowJWTInerceptor(jwt, nil), grpcMiddleware.NewSanitiseErrorsInerceptor())

	runnerService, err := grpcRunner.NewHandler(jwt, orchestratorClient, runner)
	if err != nil {
		return nil, err
	}

	runnerPath, runnerHandler := pbConnect.NewRunnerServiceHandler(runnerService, grpcMiddlewareHandler)

	grpc.Handle(runnerPath, runnerHandler)

	mux := http.NewServeMux()
	mux.Handle("/grpc/", http.StripPrefix("/grpc", grpc))

	return &GRPCServer{server: mux}, nil
}

func NewJobServer(orchestratorClient pbConnect.OrchestratorServiceClient, jwt jwt.JWTHandler, address string, workflowMeta *pb.WorkflowMeta) (*GRPCServer, error) {
	grpc := http.NewServeMux()

	jobService, err := grpcJob.NewHandler(orchestratorClient, workflowMeta, jwt, address)
	if err != nil {
		return nil, err
	}

	grpcMiddlewareHandler := connect.WithInterceptors(grpcMiddleware.NewWorkflowJWTInerceptor(jwt, &workflowMeta.WorkflowJwt), grpcMiddleware.NewSanitiseErrorsInerceptor())

	jobPath, jobHandler := pbConnect.NewJobServiceHandler(jobService, grpcMiddlewareHandler)

	grpc.Handle(jobPath, jobHandler)

	mux := http.NewServeMux()
	mux.Handle("/grpc/", http.StripPrefix("/grpc", grpc))

	return &GRPCServer{server: mux}, nil
}
