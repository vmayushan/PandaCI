package grpcMiddleware

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/alfiejones/panda-ci/pkg/jwt"
	pb "github.com/alfiejones/panda-ci/proto/go/v1"
	"github.com/rs/zerolog/log"
)

func getWorkflowJWT(msg any) (string, error) {
	if m, ok := msg.(interface{ GetWorkflowMeta() *pb.WorkflowMeta }); ok {
		return m.GetWorkflowMeta().GetWorkflowJwt(), nil
	} else if m, ok := msg.(interface{ GetWorkflowJwt() string }); ok {
		return m.GetWorkflowJwt(), nil
	}
	return "", errors.New("no workflow meta found")
}

func NewWorkflowJWTInerceptor(jwt jwt.JWTHandler, expectedJWT *string) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {

			if req.Spec().IsClient {
				workflowJWT, err := getWorkflowJWT(req.Any())
				if err != nil {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						err,
					)
				}

				req.Header().Set("Authorization", "Bearer "+workflowJWT)
			} else {

				workflowJWT := strings.TrimPrefix(req.Header().Get("Authorization"), "Bearer ")

				if workflowJWT == "" {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("no token provided"),
					)
				}

				if expectedJWT != nil && workflowJWT != *expectedJWT {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("invalid token"),
					)
				}

				workflowClaims, err := jwt.ValidateWorkflowToken(workflowJWT)
				if err != nil {
					log.Error().Err(err).Msg("Failed to validate workflow token")
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("invalid token"),
					)
				}

				ctx = context.WithValue(ctx, "workflowClaims", workflowClaims)
			}

			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func GetWorkflowClaims(ctx context.Context) jwt.WorkflowClaims {
	claims, ok := ctx.Value("workflowClaims").(jwt.WorkflowClaims)
	if !ok {
		return jwt.WorkflowClaims{}
	}
	return claims
}
