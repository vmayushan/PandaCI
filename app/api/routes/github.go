package routes

import (
	"github.com/labstack/echo/v4"
	handlersGithub "github.com/pandaci-com/pandaci/app/api/handlers/github"
	"github.com/pandaci-com/pandaci/app/api/middleware"
	"github.com/pandaci-com/pandaci/app/git"
	"github.com/pandaci-com/pandaci/app/orchestrator"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/app/scanner"
	"github.com/pandaci-com/pandaci/types"
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
