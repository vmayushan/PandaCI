package main

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/pandaci-com/pandaci/app/git"
	grpcMiddleware "github.com/pandaci-com/pandaci/app/grpc/middleware"
	"github.com/pandaci-com/pandaci/app/orchestrator"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/app/scanner"
	scannerGit "github.com/pandaci-com/pandaci/app/scanner/git"
	"github.com/pandaci-com/pandaci/pkg/jwt"
	"github.com/pandaci-com/pandaci/pkg/retryClient"
	"github.com/pandaci-com/pandaci/pkg/utils/env"
	"github.com/pandaci-com/pandaci/platform/database"
	"github.com/pandaci-com/pandaci/platform/storage"
	pbConnect "github.com/pandaci-com/pandaci/proto/go/v1/v1connect"
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
