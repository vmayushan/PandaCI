package handlersGithub

import (
	"github.com/labstack/echo/v4"
	"github.com/pandaci-com/pandaci/app/api/middleware"
	gitShared "github.com/pandaci-com/pandaci/app/git/shared"
	"github.com/pandaci-com/pandaci/app/orchestrator"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/app/scanner"
	"github.com/pandaci-com/pandaci/types"
)

type Handler struct {
	queries        *queries.Queries
	getCurrentUser func(echo.Context) types.User
	githubClient   gitShared.Client
	scanner        scanner.Handler
	orchestrator   orchestrator.Handler
}

func NewHandler(queries *queries.Queries, githubClient gitShared.Client, scanner scanner.Handler, orchestrator orchestrator.Handler) *Handler {
	return &Handler{
		queries:        queries,
		getCurrentUser: middleware.GetUser,
		githubClient:   githubClient,
		scanner:        scanner,
		orchestrator:   orchestrator,
	}
}
