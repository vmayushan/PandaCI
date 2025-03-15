package handlersOry

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pandaci-com/pandaci/pkg/gravatar"
	"github.com/pandaci-com/pandaci/pkg/utils/env"
	"github.com/pandaci-com/pandaci/platform/identity"
	typesDB "github.com/pandaci-com/pandaci/types/database"
	"github.com/rs/zerolog/log"
)

func (h *Handler) HandleAfterRegistration(c echo.Context) error {
	ctx := context.Background()

	type Payload struct {
		ID     string `json:"id"`
		Traits struct {
			Email  string `json:"email"`
			Name   string `json:"name"`
			Avatar string `json:"avatar"`
		} `json:"traits"`
	}

	apiKey := strings.Trim(c.Request().Header.Get("Authorization"), "Basic ")

	oryKey, err := env.GetOryWebhookAPIKey()
	if err != nil {
		return err
	}

	if apiKey != *oryKey {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var payload Payload
	if err := c.Bind(&payload); err != nil {
		return err
	}

	// if payload.Traits.Avatar == "" {
	gravatar := gravatar.NewGravatarFromEmail(payload.Traits.Email)
	payload.Traits.Avatar = gravatar.GetURL()

	if err := identity.UpdateUserTraits(ctx, payload.ID, payload.Traits); err != nil {
		log.Error().Err(err).Msg("failed to update user traits")
	}
	// }

	invites, err := h.queries.GetOrgInvitesByEmail(ctx, payload.Traits.Email)
	if err != nil {
		return err
	}

	for _, invite := range invites {
		if err := h.queries.AddUserToOrg(ctx, &typesDB.OrgUsersDB{
			OrgID:  invite.OrgID,
			Role:   invite.Role,
			UserID: payload.ID,
		}); err != nil {
			log.Error().Err(err).Msg("failed to create org user")
		}
	}

	// Cleanup old invite
	// TODO: This should be done in a cron job
	if err := h.queries.CleanOldOrgInvites(ctx); err != nil {
		log.Error().Err(err).Msg("failed to clean old org invites")
	}

	return c.NoContent(http.StatusOK)
}
