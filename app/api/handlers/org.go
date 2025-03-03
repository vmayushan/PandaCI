package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"slices"
	"strings"

	middleware_loaders "github.com/pandaci-com/pandaci/app/api/middleware/loaders"
	"github.com/pandaci-com/pandaci/app/email"
	"github.com/pandaci-com/pandaci/pkg/utils"
	"github.com/pandaci-com/pandaci/platform/identity"
	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
	typesHTTP "github.com/pandaci-com/pandaci/types/http"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *Handler) CreateOrg(c echo.Context) error {
	org := typesHTTP.OrgHTTP{}

	if err := c.Bind(&org); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := h.getCurrentUser(c)

	if org.OwnerID != "" && org.OwnerID != user.ID {
		return echo.NewHTTPError(http.StatusUnauthorized, "Cannot create organization for another user")
	} else {
		org.OwnerID = user.ID
	}

	log.Info().Any("org", org).Msg("Org data")

	if !utils.IsURLNameValid(org.Slug) {
		return echo.NewHTTPError(http.StatusBadRequest, "urlName is reserved, too short or contains invlaid characters, use a different urlName with at least 3 characters")
	}

	orgDB := typesDB.OrgDB{
		Slug:      org.Slug,
		Name:      org.Name,
		OwnerID:   org.OwnerID,
		AvatarURL: org.AvatarURL,
	}

	if err := h.queries.CreateOrg(c.Request().Context(), &orgDB); err != nil {
		if utils.CheckConstraintError(err, "org_slug_key") {
			return echo.NewHTTPError(http.StatusBadRequest, "Organization with this urlName already exists")
		}
		return err
	}

	license, err := orgDB.GetLicense()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, typesHTTP.OrgHTTP{
		ID:        orgDB.ID,
		Name:      orgDB.Name,
		Slug:      orgDB.Slug,
		License:   license,
		OwnerID:   orgDB.OwnerID,
		AvatarURL: orgDB.AvatarURL,
	})
}

func (h *Handler) GetCurrentUserOrgs(c echo.Context) error {
	user := h.getCurrentUser(c)

	orgs, err := h.queries.GetUsersOrgs(c.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	httpOrgs := make([]typesHTTP.OrgHTTP, len(*orgs))
	for i, org := range *orgs {
		license, err := org.GetLicense()
		if err != nil {
			return err
		}
		httpOrgs[i] = typesHTTP.OrgHTTP{
			ID:               org.ID,
			Slug:             org.Slug,
			OwnerID:          org.OwnerID,
			License:          license,
			Name:             org.Name,
			AvatarURL:        org.AvatarURL,
			CurrentUsersRole: org.CurrentUsersRole,
		}
	}

	return c.JSON(http.StatusOK, httpOrgs)
}

func (h *Handler) GetOrg(c echo.Context) error {
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

	return c.JSON(http.StatusOK, typesHTTP.OrgHTTP{
		ID:               org.ID,
		Slug:             org.Slug,
		Name:             org.Name,
		OwnerID:          org.OwnerID,
		License:          license,
		AvatarURL:        org.AvatarURL,
		CurrentUsersRole: org.CurrentUsersRole,
	})
}

func (h *Handler) UpdateOrg(c echo.Context) error {
	org, err := middleware_loaders.GetOrg(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get org")
		return err
	}

	orgReq := typesHTTP.UpdateOrgHTTP{}
	if err := c.Bind(&orgReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !utils.IsURLNameValid(orgReq.Slug) {
		return echo.NewHTTPError(http.StatusBadRequest, "urlName is reserved, too short or contains invlaid characters, use a different urlName with at least 3 characters")
	}

	if orgReq.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}

	org.Name = orgReq.Name
	org.Slug = orgReq.Slug
	org.AvatarURL = orgReq.AvatarURL

	if err := h.queries.UpdateOrg(c.Request().Context(), org); err != nil {
		if utils.CheckConstraintError(err, "org_slug_key") {
			return echo.NewHTTPError(http.StatusBadRequest, "Organization with this urlName already exists")
		}

		log.Error().Err(err).Msg("Failed to update org")
		return err
	}

	license, err := org.GetLicense()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get license")
		return err
	}

	return c.JSON(http.StatusOK, typesHTTP.OrgHTTP{
		ID:        org.ID,
		Slug:      org.Slug,
		Name:      org.Name,
		OwnerID:   org.OwnerID,
		License:   license,
		AvatarURL: org.AvatarURL,
	})
}

func (h *Handler) GetOrgUsers(c echo.Context) error {
	org, err := middleware_loaders.GetOrg(c)
	if err != nil {
		return err
	}

	orgUsers, err := h.queries.GetOrgUsers(c.Request().Context(), org.ID)
	if err != nil {
		return err
	}

	userIDs := make([]string, len(*orgUsers))
	for i, orgUser := range *orgUsers {
		userIDs[i] = orgUser.UserID
	}

	oryUsers, err := identity.ListUsersByIDs(c.Request().Context(), userIDs)
	if err != nil {
		return err
	}

	httpOrgUsers := []typesHTTP.OrgUserHTTP{}
	for _, orgUser := range *orgUsers {
		oryUserIndex := slices.IndexFunc(oryUsers, func(oryUser *types.User) bool {
			return oryUser.ID == orgUser.UserID
		})
		if oryUserIndex == -1 {
			log.Info().Str("userID", orgUser.UserID).Msg("User not found in ory")
			go func() {
				if err := h.queries.RemoveUserFromOrg(context.Background(), org.ID, orgUser.UserID); err != nil {
					log.Error().Err(err).Msg("Failed to remove user from org")
				}
			}()
			continue
		}
		httpOrgUsers = append(httpOrgUsers, typesHTTP.OrgUserHTTP{
			User: types.User{
				ID:    orgUser.UserID,
				Email: oryUsers[oryUserIndex].Email,
				Name:  oryUsers[oryUserIndex].Name,
			},
			Role: orgUser.Role,
		})
	}

	return c.JSON(http.StatusOK, httpOrgUsers)
}

func (h *Handler) InviteUserToOrg(c echo.Context) error {

	currentUser := h.getCurrentUser(c)

	type inviteUserToOrgRequest struct {
		Email string `json:"email"`
	}

	org, err := middleware_loaders.GetOrg(c)
	if err != nil {
		return err
	}

	req := inviteUserToOrgRequest{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if user exists
	user, err := identity.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by email")
		return err
	}

	invitedByName := strings.Split(currentUser.Email, "@")[0]
	if currentUser.Name != nil {
		invitedByName = *currentUser.Name
	}

	if user == nil {
		// no user so we create an invite which is then used in the registration webhook
		if err := h.queries.CreateOrgInvite(c.Request().Context(), typesDB.PendingOrgInviteDB{
			OrgID: org.ID,
			Email: req.Email,
			Role:  types.ORG_USERS_ROLE_MEMBER,
		}); err != nil {
			log.Error().Err(err).Msg("Failed to create org invite")
			return err
		}

		// using background ctx since we've already added the invite to the DB
		if err := h.email.SendOrgInviteEmailToNewUser(context.Background(), email.OrgInviteEmailNewUserData{
			To:             req.Email,
			InvitedByEmail: currentUser.Email,
			InvitedByName:  invitedByName,
			TeamName:       org.Name,
			TeamSlug:       org.Slug,
			TeamAvatarURL:  org.AvatarURL,
		}); err != nil {
			log.Error().Err(err).Msg("Failed to send invite email")
			return err
		}

		return c.NoContent(http.StatusOK)
	}

	orgUser := typesDB.OrgUsersDB{
		OrgID:  org.ID,
		UserID: user.ID,
		Role:   types.ORG_USERS_ROLE_MEMBER,
	}

	if err := h.queries.AddUserToOrg(c.Request().Context(), &orgUser); err != nil {
		if utils.CheckConstraintError(err, "org_users_pkey") {
			return echo.NewHTTPError(http.StatusBadRequest, "User is already a member of this organization")
		}
		return err
	}

	inviteeName := strings.Split(user.Email, "@")[0]
	if user.Name != nil {
		inviteeName = *user.Name
	}

	// using background ctx since we've already added the user to the org
	if err := h.email.SendOrgInviteEmailToExistingUser(context.Background(), email.OrgInviteEmailExistingUserData{
		To:             req.Email,
		InvitedByEmail: currentUser.Email,
		InvitedByName:  invitedByName,
		TeamName:       org.Name,
		TeamSlug:       org.Slug,
		TeamAvatarURL:  org.AvatarURL,
		InviteeName:    inviteeName,
	}); err != nil {
		log.Error().Err(err).Msg("Failed to send invite email")
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) RemoveUserFromOrg(c echo.Context) error {
	org, err := middleware_loaders.GetOrg(c)
	if err != nil {
		return err
	}

	userID := c.Param("user_id")

	if org.OwnerID == userID {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot remove the owner of the organization")
	}

	if err := h.queries.RemoveUserFromOrg(c.Request().Context(), org.ID, userID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) DeleteOrg(c echo.Context) error {
	user := h.getCurrentUser(c)

	org, err := middleware_loaders.GetOrg(c)
	if err != nil {
		return err
	}

	if err := h.queries.DeleteOrg(c.Request().Context(), org.ID); err != nil {
		log.Error().Err(err).Msg("Failed to delete org")
		return err
	}

	if _, err := h.queries.GetUsersFreeOrg(context.Background(), user.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error().Err(err).Msg("Failed to get free org")
		return err
	} else if err != nil {
		// User does not have a free org so we pick a paused one (if any) and make it free
		orgs, err := h.queries.GetUsersOrgs(context.Background(), user.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get user orgs")
			return err
		}

		for _, org := range *orgs {
			if org.OwnerID != user.ID {
				continue
			}
			license, err := org.GetLicense()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get license")
				return err
			}

			if license.Plan == types.CloudSubscriptionPlanPaused {
				log.Debug().Str("orgID", org.ID).Msg("Setting org to free")
				license.Plan = types.CloudSubscriptionPlanFree
				license.Features = types.Features{
					BuildMinutes:        6000,
					MaxBuildMinutes:     6000,
					Committers:          5,
					MaxCommitters:       5,
					MaxCloudRunnerScale: 4,
					MaxProjects:         10,
				}
				if err := org.SetLicense(license); err != nil {
					log.Error().Err(err).Msg("Failed to set license")
					return err
				}
				if err := h.queries.UpdateOrgLicense(context.Background(), &org); err != nil {
					log.Error().Err(err).Msg("Failed to update org")
					return err
				}
				break
			}
		}
	}

	return c.NoContent(http.StatusOK)
}
