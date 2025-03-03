package orchestrator

import (
	"github.com/pandaci-com/pandaci/app/git"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/pkg/jwt"
	pbConnect "github.com/pandaci-com/pandaci/proto/go/v1/v1connect"
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
