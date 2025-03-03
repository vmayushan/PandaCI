package middlewareLoaders

import (
	"github.com/labstack/echo/v4"
	"github.com/pandaci-com/pandaci/app/api/middleware"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/types"
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
