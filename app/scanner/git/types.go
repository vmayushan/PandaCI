package scannerGit

import (
	"github.com/pandaci-com/pandaci/app/git"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/app/scanner"
	"github.com/pandaci-com/pandaci/pkg/jwt"
)

type Handler struct {
	queries    *queries.Queries
	gitHandler *git.Handler
	jwtHandler *jwt.JWTHandler
}

func NewHandler(queries *queries.Queries, gitHandler *git.Handler, jwtHandler *jwt.JWTHandler) scanner.Handler {
	return &Handler{
		queries:    queries,
		gitHandler: gitHandler,
		jwtHandler: jwtHandler,
	}
}
