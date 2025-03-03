package handlersProject

import (
	"fmt"
	"net/http"

	"github.com/pandaci-com/pandaci/app/api/middleware"
	middlewareLoaders "github.com/pandaci-com/pandaci/app/api/middleware/loaders"
	"github.com/pandaci-com/pandaci/pkg/utils"
	"github.com/pandaci-com/pandaci/platform/analytics"
	typesDB "github.com/pandaci-com/pandaci/types/database"
	typesHTTP "github.com/pandaci-com/pandaci/types/http"
	"github.com/labstack/echo/v4"
	"github.com/posthog/posthog-go"
	"github.com/rs/zerolog/log"
)

func (h *Handler) CreateProject(c echo.Context) error {
	user := middleware.GetUser(c)

	org, err := middlewareLoaders.GetOrg(c)
	if err != nil {
		return err
	}

	license, err := org.GetLicense()
	if err != nil {
		return err
	}

	projectCount, err := h.queries.CountOrgProjects(c.Request().Context(), org)
	if err != nil {
		return err
	}

	if projectCount >= license.Features.MaxProjects {
		return echo.NewHTTPError(http.StatusBadRequest, "You have reached the maximum number of projects allowed for your plan")
	}

	projectReq := typesHTTP.CreateProjectHTTP{}

	if err := c.Bind(&projectReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	log.Debug().Any("project", projectReq).Msg("Creating new project")

	if !utils.IsURLNameValid(projectReq.Slug) {
		return echo.NewHTTPError(http.StatusBadRequest, "urlName is reserved, too short or contains invalid characters, use a different urlName with at least 3 characters")
	}

	gitClient, err := h.gitHandler.GetClient(projectReq.GitProviderType)
	if err != nil {
		return err
	}

	gitAppClient, err := gitClient.NewAppClient(c.Request().Context())
	if err != nil {
		return err
	}

	installation, err := gitAppClient.GetInstallation(c.Request().Context(), projectReq.GitProviderIntegrationID)
	if err != nil {
		return err
	}

	if installation == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "installation doesn't exist")
	}

	integration := &typesDB.GitIntegration{
		ProviderID:        projectReq.GitProviderIntegrationID,
		Type:              projectReq.GitProviderType,
		ProviderAccountID: installation.AccountID,
	}

	if err := h.queries.GetInsertOrUpdateGitIntegration(c.Request().Context(), integration); err != nil {
		return err
	}

	gitInstallation, err := gitClient.NewInstallationClient(c.Request().Context(), installation.ID)
	if err != nil {
		return err
	}

	if exists, err := gitInstallation.RepoExists(c.Request().Context(), projectReq.GitProviderRepoID); err != nil {
		return err
	} else if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "repo doesn't exist")
	}

	projectDB := typesDB.Project{
		OrgID:             org.ID,
		Slug:              projectReq.Slug,
		Name:              projectReq.Name,
		GitIntegrationID:  integration.ID,
		GitProviderRepoID: projectReq.GitProviderRepoID,
	}

	if err := h.queries.CreateProject(c.Request().Context(), org, &projectDB); err != nil {

		if utils.CheckConstraintError(err, "project_url_name") {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("a project with the urlName: %s already exists", projectReq.Slug))
		}

		return err
	}

	analytics.TrackUserProjectEvent(user, projectDB, posthog.Capture{
		Event: "project_created",
	})

	return c.JSON(http.StatusCreated, typesHTTP.ProjectHTTP{
		ID:        projectDB.ID,
		Name:      projectDB.Name,
		Slug:      projectDB.Slug,
		OrgID:     projectDB.OrgID,
		AvatarURL: projectDB.AvatarURL,
	})
}

func (h *Handler) GetOrgProjects(c echo.Context) error {
	org, err := middlewareLoaders.GetOrg(c)
	if err != nil {
		return err
	}

	projects, err := h.queries.GetOrgProjects(c.Request().Context(), org)
	if err != nil {
		return err
	}

	httpProjects := make([]typesHTTP.ProjectHTTP, len(*projects))
	for i, project := range *projects {
		httpProjects[i] = typesHTTP.ProjectHTTP{
			ID:        project.ID,
			Slug:      project.Slug,
			Name:      project.Name,
			OrgID:     project.OrgID,
			AvatarURL: project.AvatarURL,
			LastBuild: project.LastBuild,
		}
	}

	return c.JSON(http.StatusOK, httpProjects)
}

func (h *Handler) GetProject(c echo.Context) error {
	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, typesHTTP.ProjectHTTP{
		ID:        project.ID,
		Name:      project.Name,
		Slug:      project.Slug,
		OrgID:     project.OrgID,
		AvatarURL: project.AvatarURL,
	})
}

func (h *Handler) UpdateProject(c echo.Context) error {
	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	projectReq := typesHTTP.UpdateProjectHTTP{}

	if err := c.Bind(&projectReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !utils.IsURLNameValid(projectReq.Slug) {
		return echo.NewHTTPError(http.StatusBadRequest, "urlName is reserved, too short or contains invalid characters, use a different urlName with at least 3 characters")
	}

	if projectReq.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "displayName is required")
	}

	project.Slug = projectReq.Slug
	project.Name = projectReq.Name
	project.AvatarURL = projectReq.AvatarURL

	if err := h.queries.UpdateProject(c.Request().Context(), project); err != nil {
		if utils.CheckConstraintError(err, "project.org_id, project.url_name") {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("a project with the urlName: %s already exists", projectReq.Slug))
		}
		return err
	}

	return c.JSON(http.StatusOK, typesHTTP.ProjectHTTP{
		ID:        project.ID,
		Name:      project.Name,
		Slug:      project.Slug,
		OrgID:     project.OrgID,
		AvatarURL: project.AvatarURL,
	})
}

func (h *Handler) DeleteProject(c echo.Context) error {
	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	if err := h.queries.DeleteProject(c.Request().Context(), project); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
