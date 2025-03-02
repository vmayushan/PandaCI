package orchestrator

import (
	"github.com/alfiejones/panda-ci/app/git"
	"github.com/alfiejones/panda-ci/app/queries"
	"github.com/alfiejones/panda-ci/pkg/jwt"
	pbConnect "github.com/alfiejones/panda-ci/proto/go/v1/v1connect"
)

type Handler struct {
	queries      *queries.Queries
	jwt          jwt.JWTHandler
	runnerClient pbConnect.RunnerServiceClient
	git          *git.Handler
}

func NewOrchestratorHandler(queries *queries.Queries, jwt jwt.JWTHandler, runnerClient pbConnect.RunnerServiceClient, git *git.Handler) *Handler {

	return &Handler{
		queries:      queries,
		jwt:          jwt,
		runnerClient: runnerClient,
		git:          git,
	}
}
