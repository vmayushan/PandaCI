package git

import (
	"github.com/alfiejones/panda-ci/app/queries"
)

type Handler struct {
	queries *queries.Queries
}

func NewGitHandler(queries *queries.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}
