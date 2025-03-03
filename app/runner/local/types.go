package runnerLocal

import (
	"github.com/pandaci-com/pandaci/app/runner"
	"github.com/pandaci-com/pandaci/pkg/jwt"
	"github.com/pandaci-com/pandaci/types"
	"github.com/docker/docker/client"
)

type (
	Handler struct {
		docker     *client.Client
		JWTHandler *jwt.JWTHandler
	}
)

func (h *Handler) GetRunnerType() types.RunnerType {
	return types.RunnerTypeLocal
}

func NewRunner(jwtHandler *jwt.JWTHandler) (runner.Handler, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Handler{
		docker:     docker,
		JWTHandler: jwtHandler,
	}, nil
}
