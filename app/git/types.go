package git

import (
	"github.com/pandaci-com/pandaci/app/queries"
)

type Handler struct {
	queries *queries.Queries
}

func NewGitHandler(queries *queries.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}
