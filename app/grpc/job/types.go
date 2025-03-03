package grpcJob

import (
	"github.com/pandaci-com/pandaci/app/grpc/job/service"
	"github.com/pandaci-com/pandaci/pkg/jwt"
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
