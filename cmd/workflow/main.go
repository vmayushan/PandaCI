package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"

	"connectrpc.com/connect"
	grpcMiddleware "github.com/pandaci-com/pandaci/app/grpc/middleware"
	"github.com/pandaci-com/pandaci/pkg/jwt"
	"github.com/pandaci-com/pandaci/pkg/retryClient"
	pb "github.com/pandaci-com/pandaci/proto/go/v1"
	"github.com/pandaci-com/pandaci/types"
	"github.com/phayes/freeport"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	pbConnect "github.com/pandaci-com/pandaci/proto/go/v1/v1connect"
)

func getConfig() (*pb.WorkflowRunnerInitConfig, error) {
	// We pass a proto object as an argument
	// this gives us the info we require about this workflow run
	rawConfig := []byte(os.Args[1])
	configBytes, err := base64.StdEncoding.DecodeString(string(rawConfig))
	if err != nil {
		log.Error().Err(err).Msg("decoding workflow config")
		return nil, err
	}

	config := &pb.WorkflowRunnerInitConfig{}
	if err := proto.Unmarshal(configBytes, config); err != nil {
		log.Error().Err(err).Msg("unmarshalling workflow config")
		return nil, err
	}

	return config, nil
}

func main() {
	ctx := context.Background()
	logWriter := initLogs()

	config, err := getConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("getting config")
	}

	jwtHandler := jwt.NewJWTHelper(jwt.JWTOpts{
		PublicKey: &config.PublicJwtKey,
	})

	orchestratorClient := pbConnect.NewOrchestratorServiceClient(&http.Client{
		Transport: &retryClient.RetryRoundTripper{
			Base: http.DefaultTransport,
		},
	}, config.OrchestratorAddress, connect.WithInterceptors(grpcMiddleware.NewWorkflowJWTInerceptor(jwtHandler, nil)))

	errorExit := func(err error, msg string) {
		log.Error().Err(err).Msg(msg)

		if _, err := orchestratorClient.CreateWorkflowAlert(ctx, connect.NewRequest(&pb.OrchestratorServiceCreateWorkflowAlertRequest{
			WorkflowMeta: config.WorkflowMeta,
			Alert: &pb.WorkflowAlert{
				Type:    pb.WorkflowAlert_TYPE_ERROR,
				Title:   "Workflow Error - " + msg,
				Message: err.Error(),
			},
		})); err != nil {
			log.Error().Err(err).Msg("creating workflow alert")
		}

		if err := logWriter.UploadLogs(ctx, config); err != nil {
			log.Error().Err(err).Msg("uploading logs")
		}

		if _, err := orchestratorClient.FinishWorkflow(ctx, connect.NewRequest(&pb.OrchestratorServiceFinishWorkflowRequest{
			WorkflowMeta: config.WorkflowMeta,
			Conclusion:   pb.Conclusion_CONCLUSION_FAILURE,
		})); err != nil {
			log.Error().Err(err).Msg("finishing workflow")
		}

		log.Fatal().Msg("workflow failed")
	}

	workflowClaims, err := jwtHandler.ValidateWorkflowToken(config.WorkflowMeta.WorkflowJwt)
	if err != nil {
		errorExit(err, "validating workflow jwt")
	}

	port := int(config.GetPort())
	if port == 0 {
		log.Info().Msg("no port provided, randomly selecting one")
		port, err = freeport.GetFreePort()
		if err != nil {
			errorExit(err, "getting free port")
		}
	}

	if err := setEnv(workflowClaims.WorkflowID, config, port); err != nil {
		errorExit(err, "setting env")
	}

	worstJobConclusion := make(map[string]pb.Conclusion)

	if err := StartHandler(logWriter, workflowClaims.WorkflowID, orchestratorClient, jwtHandler, config, port, worstJobConclusion); err != nil {
		errorExit(err, "starting handler")
	}

	if err := cloneRepo(ctx, config.WorkflowMeta.Repo); err != nil {
		errorExit(err, "cloning repo")
	}

	if _, err := orchestratorClient.WorkflowStarted(ctx, connect.NewRequest(&pb.OrchestratorServiceWorkflowStartedRequest{
		WorkflowMeta: config.WorkflowMeta,
	})); err != nil {
		errorExit(err, "workflow started")
	}

	if err := runWorkflow(ctx, config); err != nil {
		errorExit(err, "running workflow")
	}

	if err := logWriter.UploadLogs(ctx, config); err != nil {
		errorExit(err, "uploading logs")
	}

	workflowConclusion := pb.Conclusion_CONCLUSION_SUCCESS
	for _, conclusion := range worstJobConclusion {
		log.Info().Msgf("job conclusion: %s", conclusion.String())
		if types.CompareProtoConclusionRank(conclusion, workflowConclusion) {
			workflowConclusion = conclusion
		}
	}

	// This query can trigger the docker/vm/whatever to be destroyed
	/// Make sure we don't do anything after this
	if _, err := orchestratorClient.FinishWorkflow(ctx, connect.NewRequest(&pb.OrchestratorServiceFinishWorkflowRequest{
		WorkflowMeta: config.WorkflowMeta,
		Conclusion:   workflowConclusion,
	})); err != nil {
		errorExit(err, "finishing workflow")
	}

	log.Info().Msg("workflow finished")
}
