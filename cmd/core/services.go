package main

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/alfiejones/panda-ci/app/git"
	grpcMiddleware "github.com/alfiejones/panda-ci/app/grpc/middleware"
	"github.com/alfiejones/panda-ci/app/orchestrator"
	"github.com/alfiejones/panda-ci/app/queries"
	"github.com/alfiejones/panda-ci/app/scanner"
	scannerGit "github.com/alfiejones/panda-ci/app/scanner/git"
	"github.com/alfiejones/panda-ci/pkg/jwt"
	"github.com/alfiejones/panda-ci/pkg/retryClient"
	"github.com/alfiejones/panda-ci/pkg/utils/env"
	"github.com/alfiejones/panda-ci/platform/database"
	"github.com/alfiejones/panda-ci/platform/storage"
	pbConnect "github.com/alfiejones/panda-ci/proto/go/v1/v1connect"
	"github.com/labstack/echo/v4"
)

type services struct {
	queries             *queries.Queries
	jwt                 *jwt.JWTHandler
	gitHandler          *git.Handler
	scannerHandler      scanner.Handler
	orchestratorHandler orchestrator.Handler
	bucketClient        *storage.BucketClient
	echo                *echo.Echo
}

func getServices() (services, error) {

	e := echo.New()
	e.HideBanner = true

	db, err := database.NewPostgresDatabase()
	if err != nil {
		return services{}, err
	}

	queriesQueries := queries.NewQueries(db)

	jwtHandler := jwt.NewJWTHelper(jwt.JWTOpts{
		ExpiresIn: jwt.DefaultExpireTime,
	})

	gitHandler := git.NewGitHandler(queriesQueries)

	scannerHandler := scannerGit.NewHandler(queriesQueries, gitHandler, &jwtHandler)

	bucketClient, err := storage.GetClient(context.Background())
	if err != nil {
		return services{}, err
	}

	runnerAddress, err := env.GetRunnerAddress()
	if err != nil {
		return services{}, err
	}

	runnerClient := pbConnect.NewRunnerServiceClient(&http.Client{
		Transport: &retryClient.RetryRoundTripper{
			Base: http.DefaultTransport,
		},
	}, *runnerAddress, connect.WithInterceptors(grpcMiddleware.NewWorkflowJWTInerceptor(jwtHandler, nil)))

	orchestratorHandler := orchestrator.NewOrchestratorHandler(queriesQueries, jwtHandler, runnerClient, gitHandler)

	return services{
		queries:             queriesQueries,
		jwt:                 &jwtHandler,
		gitHandler:          gitHandler,
		scannerHandler:      scannerHandler,
		bucketClient:        bucketClient,
		orchestratorHandler: *orchestratorHandler,
		echo:                e,
	}, nil
}
