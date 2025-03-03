package handlersPaddle

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	middleware_loaders "github.com/pandaci-com/pandaci/app/api/middleware/loaders"
)

func (h *Handler) HandlePortalRequest(c echo.Context) error {
	org, err := middleware_loaders.GetOrg(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get org")
		return err
	}

	license, err := org.GetLicense()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get license")
		return err
	}

	if license.PaddleData == nil {
		return echo.NewHTTPError(404, "No paddle data found")
	}

	res, err := h.paddleClient.CustomerPortalSessionsClient.CreateCustomerPortalSession(c.Request().Context(), &paddle.CreateCustomerPortalSessionRequest{
		CustomerID:      license.PaddleData.CustomerID,
		SubscriptionIDs: []string{license.PaddleData.SubscriptionID},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create customer portal session")
		return err
	}

	type Response struct {
		GeneralURL string `json:"generalURL"`
	}

	return c.JSON(200, Response{
		GeneralURL: res.URLs.General.Overview,
	})
}
