package grpcMiddleware

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/log"
)

// If an error is not a connect.Error, it is replaced with a connect.Error with code CodeInternal
// This is to prevent internal errors from being leaked to the client
func NewSanitiseErrorsInerceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				panic("should not be called on client")
			}

			res, err := next(ctx, req)
			if err != nil {
				var _err *connect.Error
				if !errors.As(err, &_err) {
					log.Error().Err(err).Msg("error occured during rpc call, sanitising error")
					return res, connect.NewError(
						connect.CodeInternal,
						errors.New("internal server error"),
					)
				} else {
					log.Error().Err(err).Msg("error occured during rpc call")
				}
			}

			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
