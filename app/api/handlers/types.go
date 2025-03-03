package handlers

import (
	"github.com/pandaci-com/pandaci/app/api/middleware"
	"github.com/pandaci-com/pandaci/app/email"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/types"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	queries        *queries.Queries
	getCurrentUser func(echo.Context) types.User
	email          *email.Handler
}

func NewHandler(queries *queries.Queries) *Handler {

	email := email.NewHandler()

	return &Handler{
		queries:        queries,
		getCurrentUser: middleware.GetUser,
		email:          email,
	}
}
