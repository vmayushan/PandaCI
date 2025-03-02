package handlersGithub

import (
	"github.com/labstack/echo/v4"
	"github.com/alfiejones/panda-ci/app/api/middleware"
	gitShared "github.com/alfiejones/panda-ci/app/git/shared"
	"github.com/alfiejones/panda-ci/app/orchestrator"
	"github.com/alfiejones/panda-ci/app/queries"
	"github.com/alfiejones/panda-ci/app/scanner"
	"github.com/alfiejones/panda-ci/types"
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
