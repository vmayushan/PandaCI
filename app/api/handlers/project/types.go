package handlersProject

import (
	"github.com/alfiejones/panda-ci/app/api/middleware"
	"github.com/alfiejones/panda-ci/app/git"
	"github.com/alfiejones/panda-ci/app/orchestrator"
	"github.com/alfiejones/panda-ci/app/queries"
	"github.com/alfiejones/panda-ci/app/scanner"
	"github.com/alfiejones/panda-ci/pkg/jwt"
	"github.com/alfiejones/panda-ci/platform/storage"
	"github.com/alfiejones/panda-ci/types"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	queries        *queries.Queries
	getCurrentUser func(echo.Context) types.User
	gitHandler     *git.Handler
	scanner        scanner.Handler
	orchestrator   orchestrator.Handler
	bucketClient   *storage.BucketClient
	jwtHandler     *jwt.JWTHandler
}

func NewHandler(queries *queries.Queries, gitHandler *git.Handler, scanner scanner.Handler, orchestrator orchestrator.Handler, bucketClient *storage.BucketClient, jwtHandler *jwt.JWTHandler) *Handler {
	return &Handler{
		queries:        queries,
		getCurrentUser: middleware.GetUser,
		gitHandler:     gitHandler,
		scanner:        scanner,
		bucketClient:   bucketClient,
		orchestrator:   orchestrator,
		jwtHandler:     jwtHandler,
	}
}
