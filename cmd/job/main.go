package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"connectrpc.com/connect"
	"github.com/alfiejones/panda-ci/app/grpc"
	grpcMiddleware "github.com/alfiejones/panda-ci/app/grpc/middleware"
	"github.com/alfiejones/panda-ci/pkg/jwt"
	"github.com/alfiejones/panda-ci/pkg/retryClient"
	"github.com/phayes/freeport"
	"github.com/rs/zerolog/log"

	"google.golang.org/protobuf/proto"

	pb "github.com/alfiejones/panda-ci/proto/go/v1"
	pbConnect "github.com/alfiejones/panda-ci/proto/go/v1/v1connect"
)

func startDocker(ctx context.Context) {
	if output, err := exec.CommandContext(ctx, "dockerd", "--data-root", "/home/.runner/docker").CombinedOutput(); err != nil {
		log.Warn().Err(err).Bytes("out", output).Msg("failed to start docker, it might already be running")
	}
}

// TODO - have better error handling. We should ideally fail the build if before we exit

func main() {
	ctx := context.Background()

	if local := os.Getenv("PANDACI_LOCAL"); local != "true" {
		go startDocker(ctx)
	}

	// We pass a proto object as an argument
	// this gives us the info we require about this job run
	rawConfig := []byte(os.Args[1])
	configBytes, err := base64.StdEncoding.DecodeString(string(rawConfig))
	if err != nil {
		log.Fatal().Err(err).Msg("decoding config")
	}

	config := &pb.JobRunnerInitConfig{}

	if err := proto.Unmarshal(configBytes, config); err != nil {
		log.Fatal().Err(err).Msg("parsing job config")
	}

	jwtHandler := jwt.NewJWTHelper(jwt.JWTOpts{
		PublicKey: &config.PublicJwtKey,
	})

	go monitorWorkflowStatus(ctx, config, jwtHandler)

	workflowClaims, err := jwtHandler.ValidateWorkflowToken(config.WorkflowMeta.WorkflowJwt)
	if err != nil {
		log.Fatal().Err(err).Msg("validating workflow jwt")
	}

	log.Info().Msg("cloning repo")

	if err := cloneRepo(ctx, config.WorkflowMeta.Repo); err != nil {
		log.Fatal().Err(err).Msg("cloning repo")
	}

	port := int(config.GetPort())
	if port == 0 {
		log.Info().Msg("no port provided, randomly selecting one")
		port, err = freeport.GetFreePort()
		if err != nil {
			log.Fatal().Err(err).Msg("getting free port")
		}

		config.JobMeta.Address = fmt.Sprintf("http://%s:%d", config.Host, port)
	}

	if err := setEnv(workflowClaims.WorkflowID, config, port); err != nil {
		log.Fatal().Err(err).Msg("setting env")
	}

	log.Info().Msg("starting grpc server")

	grpcServer, err := grpc.NewJobServer(jwtHandler, fmt.Sprintf("%s:%d", config.Host, port), config.WorkflowMeta)
	if err != nil {
		log.Fatal().Err(err).Msg("starting grpc server")
	}

	go func() {
		log.Info().Msg("notifying runner that job has started")

		runnerClient := pbConnect.NewRunnerServiceClient(&http.Client{
			Transport: &retryClient.RetryRoundTripper{
				Base: http.DefaultTransport,
				Headers: map[string]string{
					// We need to target the machine that started this job since that machine is waiting for our response
					"fly-force-instance-id": config.JobMeta.GetFlyMeta().GetParentRunnerMachineId(),
				},
			},
		}, config.RunnerAddress, connect.WithInterceptors(grpcMiddleware.NewWorkflowJWTInerceptor(jwtHandler, nil)))

		// The runner is waiting on us to send them our address
		// it knows the host but not the port, it wants to wait until the server is up
		if _, err := runnerClient.JobStarted(ctx, connect.NewRequest(&pb.RunnerServiceJobStartedRequest{
			WorkflowMeta: config.WorkflowMeta,
			Id:           config.JobMeta.Id,
			Address:      fmt.Sprintf("http://%s:%d", config.Host, port),
		})); err != nil {
			log.Fatal().Err(err).Msg("sending job started")
		}

		log.Info().Msg("runner notified")
	}()

	if err := grpcServer.Start(fmt.Sprintf("%s:%d", config.Host, port)); err != nil {
		log.Fatal().Err(err).Msg("starting grpc server")
	}
}
