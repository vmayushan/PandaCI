package middlewareLoaders

import (
	"github.com/labstack/echo/v4"
	"github.com/alfiejones/panda-ci/app/api/middleware"
	"github.com/alfiejones/panda-ci/app/queries"
	"github.com/alfiejones/panda-ci/types"
)

type middlewareHandler struct {
	queries        *queries.Queries
	getCurrentUser func(echo.Context) types.User
}

func New(queries *queries.Queries) *middlewareHandler {
	return &middlewareHandler{
		queries:        queries,
		getCurrentUser: middleware.GetUser,
	}
}
