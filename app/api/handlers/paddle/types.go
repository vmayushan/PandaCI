package handlersPaddle

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/pandaci-com/pandaci/pkg/utils/env"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	queries      *queries.Queries
	paddleClient *paddle.SDK
}

func NewHandler(queries *queries.Queries) (*Handler, error) {

	apiKey, err := env.GetPaddleAPIKey()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get paddle api key")
		return nil, err
	}

	baseURL := paddle.SandboxBaseURL
	if env.GetStage() == "prod" {
		baseURL = paddle.ProductionBaseURL
	}

	client, err := paddle.New(
		*apiKey,
		paddle.WithBaseURL(baseURL),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create paddle client")
		return nil, err
	}

	return &Handler{
		queries:      queries,
		paddleClient: client,
	}, nil
}
