package main

import (
	"context"
	"net/http"
	"time"

	"connectrpc.com/connect"
	grpcMiddleware "github.com/alfiejones/panda-ci/app/grpc/middleware"
	"github.com/alfiejones/panda-ci/pkg/flyClient"
	"github.com/alfiejones/panda-ci/pkg/jwt"
	"github.com/alfiejones/panda-ci/pkg/retryClient"
	pb "github.com/alfiejones/panda-ci/proto/go/v1"
	pbConnect "github.com/alfiejones/panda-ci/proto/go/v1/v1connect"
	"github.com/rs/zerolog/log"
)

// Pings the workflow service every 5 seconds to check its still alive
// If dead, we'll call the orchestrator to fail the workflow
// this should cause this job to also be stopped
func monitorWorkflowStatus(ctx context.Context, config *pb.JobRunnerInitConfig, jwtHandler jwt.JWTHandler) {

	log.Debug().Str("workflow_address", config.WorkflowMeta.Address).Msg("monitoring workflow status")

	workflowClient := pbConnect.NewWorkflowServiceClient(&http.Client{
		Transport: &flyClient.FlyRoundTripper{
			Base:    http.DefaultTransport,
			AppName: &config.GetWorkflowMeta().GetFlyMeta().AppName,
			Headers: map[string]string{
				"Authorization": "Bearer " + config.WorkflowMeta.GetWorkflowJwt(),
			},
		},
	}, config.WorkflowMeta.Address+"/grpc")

	ticker := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, err := workflowClient.Ping(ctx, connect.NewRequest(&pb.WorkflowServicePingRequest{}))
			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Error().Err(err).Msg("workflow is dead")
				handleUnresponsiveWorkflow(ctx, config, jwtHandler)
				continue
			}
		}
	}
}

func handleUnresponsiveWorkflow(ctx context.Context, config *pb.JobRunnerInitConfig, jwtHandler jwt.JWTHandler) {
	orchestratorClient := pbConnect.NewOrchestratorServiceClient(&http.Client{
		Transport: &retryClient.RetryRoundTripper{
			Base: http.DefaultTransport,
		},
	}, config.OrchestratorAddress, connect.WithInterceptors(grpcMiddleware.NewWorkflowJWTInerceptor(jwtHandler, nil)))

	if _, err := orchestratorClient.CreateWorkflowAlert(ctx, connect.NewRequest(&pb.OrchestratorServiceCreateWorkflowAlertRequest{
		WorkflowMeta: config.WorkflowMeta,
		Alert: &pb.WorkflowAlert{
			Type:    pb.WorkflowAlert_TYPE_ERROR,
			Title:   "Unresponsive Workflow",
			Message: "The workflow became unresponsive. If this issue persists, please contact support.",
		},
	})); err != nil {
		log.Error().Err(err).Msg("creating workflow alert")
	}

	if _, err := orchestratorClient.FinishWorkflow(ctx, connect.NewRequest(&pb.OrchestratorServiceFinishWorkflowRequest{
		WorkflowMeta: config.WorkflowMeta,
		Conclusion:   pb.Conclusion_CONCLUSION_FAILURE,
	})); err != nil {
		log.Error().Err(err).Msg("finishing workflow")
	}

	log.Fatal().Msg("workflow is dead")
}
