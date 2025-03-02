package routes

import (
	"github.com/labstack/echo/v4"
	handlersGithub "github.com/alfiejones/panda-ci/app/api/handlers/github"
	"github.com/alfiejones/panda-ci/app/api/middleware"
	"github.com/alfiejones/panda-ci/app/git"
	"github.com/alfiejones/panda-ci/app/orchestrator"
	"github.com/alfiejones/panda-ci/app/queries"
	"github.com/alfiejones/panda-ci/app/scanner"
	"github.com/alfiejones/panda-ci/types"
)

func RegisterGithubRoutes(e *echo.Echo, queries *queries.Queries, gitHandler *git.Handler, scanner scanner.Handler, orchestrator orchestrator.Handler) error {
	githubClient, err := gitHandler.GetClient(types.GitProviderTypeGithub)
	if err != nil {
		return err
	}

	handler := handlersGithub.NewHandler(queries, githubClient, scanner, orchestrator)

	v1 := e.Group("/v1/git/github")
	v1.Use(middleware.NewOryMiddleware().Session)

	v1Unauthenticated := e.Group("/v1/git/github")

	v1.GET("/installations", handler.GetUserInstallations)
	v1.GET("/callback", handler.GithubAccountCallback)

	githubWebhookFn, err := handler.GetGithubWebhook()
	if err != nil {
		return err
	}
	v1Unauthenticated.POST("/webhook", githubWebhookFn)

	githubInstallation := v1.Group("/installations/:installation_id")

	githubInstallation.GET("/repos", handler.GetInstallationRepos)

	return nil
}
