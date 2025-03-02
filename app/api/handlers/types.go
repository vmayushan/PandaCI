package handlers

import (
	"github.com/alfiejones/panda-ci/app/api/middleware"
	"github.com/alfiejones/panda-ci/app/email"
	"github.com/alfiejones/panda-ci/app/queries"
	"github.com/alfiejones/panda-ci/types"
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
