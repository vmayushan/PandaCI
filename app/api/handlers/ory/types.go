package handlersOry

import (
	"github.com/alfiejones/panda-ci/app/queries"
)

type Handler struct {
	queries *queries.Queries
}

func NewHandler(queries *queries.Queries) (*Handler, error) {
	return &Handler{
		queries: queries,
	}, nil
}
