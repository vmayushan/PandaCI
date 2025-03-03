package main

import (
	"github.com/pandaci-com/pandaci/app/api/middleware"
	"github.com/pandaci-com/pandaci/app/api/routes"
	"github.com/pandaci-com/pandaci/app/git"
	"github.com/pandaci-com/pandaci/app/grpc"
	"github.com/pandaci-com/pandaci/pkg/utils/env"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func main() {
	if env.GetStage() == env.StageDev {
		if err := godotenv.Load(".env"); err != nil {
			log.Info().Msg("No .env file found")
		}
	}

	log.Info().Msgf("Stage: %s", env.GetStage())

	services, err := getServices()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to get services")
	}

	e := services.echo

	e.Use(middleware.Logger())
	e.Use(echoMiddleware.Recover())

	allowedOrigins, err := env.GetAllowedOrigins()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to get allowed origins")
	}

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(echoMiddleware.Secure())
	e.Use(middleware.TranslateErrors)

	routes.RegisterOrgRoutes(e, services.queries, services.gitHandler, services.scannerHandler, services.orchestratorHandler, services.bucketClient, services.jwt)

	if err := routes.RegisterPaddleRoutes(e, services.queries); err != nil {
		log.Fatal().Err(err).Msg("failed to register paddle routes")
	}

	if err := routes.RegisterOryRoutes(e, services.queries); err != nil {
		log.Fatal().Err(err).Msg("failed to register ory routes")
	}

	gitClient := git.NewGitHandler(services.queries)

	if err := routes.RegisterGithubRoutes(e, services.queries, gitClient, services.scannerHandler, services.orchestratorHandler); err != nil {
		log.Fatal().Err(err).Msg("failed to register github routes")
	}

	if err := grpc.RegisterOrchestratorGRPC(e, services.queries, *services.jwt, services.bucketClient, gitClient); err != nil {
		log.Fatal().Err(err).Msg("Unable to register orchestrator grpc")
	}

	apiHost, err := env.GetAPIHost()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to get api host")
	}

	if env.GetStage() == env.StageDev || env.GetStage() == env.StageLocal {
		e.Debug = true
	}

	if err := e.Start(*apiHost); err != nil {
		log.Fatal().Err(err).Msg("Unable to start server")
	}
}
