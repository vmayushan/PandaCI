package handlersGithub

import (
	"database/sql"
	"net/http"

	gitShared "github.com/alfiejones/panda-ci/app/git/shared"
	"github.com/alfiejones/panda-ci/pkg/utils"
	"github.com/alfiejones/panda-ci/pkg/utils/env"
	typesDB "github.com/alfiejones/panda-ci/types/database"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *Handler) GithubAccountCallback(c echo.Context) error {
	user := h.getCurrentUser(c)

	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Code is required")
	}
	if state == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "State is required")
	}

	if accountState, err := h.queries.GetAccountRefreshStateByID(c.Request().Context(), state); err != nil && err != sql.ErrNoRows {
		return err
	} else if err == sql.ErrNoRows || accountState.UserID != user.ID || accountState.Type != typesDB.UserAccountTypeGithub {
		log.Info().Msg("Bad github oauth state")
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid oauth state")
	}

	account := &typesDB.UserAccount{
		UserID: user.ID,
	}

	if err := h.githubClient.RefreshOAuthTokens(c.Request().Context(), user, account, nil, &code); err != nil {
		return err
	}

	frontendURL, err := env.GetFrontendURL()
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, *frontendURL)
}

func (h *Handler) GetUserInstallations(c echo.Context) error {
	user := h.getCurrentUser(c)

	page, err := utils.GetQueryParamInt(c, "page")
	if err != nil {
		return err
	}

	perPage, err := utils.GetQueryParamInt(c, "perPage")
	if err != nil {
		return err
	}

	userClient, err := h.githubClient.NewUserClient(c.Request().Context(), user)
	if err != nil {
		return err
	}

	installations, err := userClient.GetInstallations(c.Request().Context(), gitShared.GetInstallationsOptions{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {

		return err
	}

	return c.JSON(http.StatusOK, installations)
}

func (h *Handler) GetInstallationRepos(c echo.Context) error {

	installationID := c.Param("installation_id")

	query := c.QueryParam("query")

	owner := c.QueryParam("owner")
	name := c.QueryParam("name")

	if (name != "" && owner == "") || (owner != "" && name == "") {
		return echo.NewHTTPError(http.StatusBadRequest, "either owner or name cannot be set whilst the other is empty")
	}

	user := h.getCurrentUser(c)

	log.Debug().Msg("getting user client")

	userClient, err := h.githubClient.NewUserClient(c.Request().Context(), user)
	if err != nil {
		return err
	}

	log.Debug().Msg("got user client")

	repos, err := userClient.GetRepositories(c.Request().Context(), installationID, gitShared.GetRepositoriesOptions{
		Query: query,
		Owner: owner,
		Name:  name,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, repos)
}
