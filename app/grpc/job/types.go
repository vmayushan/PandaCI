package grpcJob

import (
	"github.com/alfiejones/panda-ci/app/grpc/job/service"
	"github.com/alfiejones/panda-ci/pkg/jwt"
)

type (
	Handler struct {
		jwt     jwt.JWTHandler
		service *service.Handler
		address string
	}
)

func NewHandler(jwt jwt.JWTHandler, address string) (*Handler, error) {

	service, err := service.NewHandler()
	if err != nil {
		return nil, err
	}

	return &Handler{
		jwt:     jwt,
		service: service,
		address: address,
	}, nil
}
