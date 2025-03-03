package handlersOry

import (
	"github.com/pandaci-com/pandaci/app/queries"
)

type Handler struct {
	queries *queries.Queries
}

func NewHandler(queries *queries.Queries) (*Handler, error) {
	return &Handler{
		queries: queries,
	}, nil
}
