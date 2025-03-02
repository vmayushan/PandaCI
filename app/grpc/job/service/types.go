package service

import "github.com/docker/docker/client"

type Handler struct {
	docker *client.Client
}

func NewHandler() (*Handler, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Handler{
		docker: docker,
	}, nil
}
