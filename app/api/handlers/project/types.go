package handlersProject

import (
	"github.com/pandaci-com/pandaci/app/api/middleware"
	"github.com/pandaci-com/pandaci/app/git"
	"github.com/pandaci-com/pandaci/app/orchestrator"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/app/scanner"
	"github.com/pandaci-com/pandaci/pkg/jwt"
	"github.com/pandaci-com/pandaci/platform/storage"
	"github.com/pandaci-com/pandaci/types"
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
