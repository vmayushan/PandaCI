package runnerLocal

import (
	"github.com/alfiejones/panda-ci/app/runner"
	"github.com/alfiejones/panda-ci/pkg/jwt"
	"github.com/alfiejones/panda-ci/types"
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
