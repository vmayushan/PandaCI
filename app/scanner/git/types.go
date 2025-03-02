package scannerGit

import (
	"github.com/alfiejones/panda-ci/app/git"
	"github.com/alfiejones/panda-ci/app/queries"
	"github.com/alfiejones/panda-ci/app/scanner"
	"github.com/alfiejones/panda-ci/pkg/jwt"
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
