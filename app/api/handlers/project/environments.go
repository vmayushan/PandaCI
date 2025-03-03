package handlersProject

import (
	"net/http"

	"github.com/pandaci-com/pandaci/app/api/middleware"
	middlewareLoaders "github.com/pandaci-com/pandaci/app/api/middleware/loaders"
	utilsValidator "github.com/pandaci-com/pandaci/pkg/utils/validator"
	"github.com/pandaci-com/pandaci/platform/analytics"
	typesDB "github.com/pandaci-com/pandaci/types/database"
	typesHTTP "github.com/pandaci-com/pandaci/types/http"
	"github.com/labstack/echo/v4"
	"github.com/posthog/posthog-go"
)

func (h *Handler) CreateProjectEnvironment(c echo.Context) error {

	user := middleware.GetUser(c)

	environmentReq := typesHTTP.CreateProjectEnvironmentBody{}

	if err := c.Bind(&environmentReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	validate := utilsValidator.NewValidator()
	if err := validate.Struct(environmentReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utilsValidator.ValidatorErrors(err))
	}

	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	environment := typesDB.ProjectEnvironment{
		Name:          environmentReq.Name,
		ProjectID:     project.ID,
		BranchPattern: environmentReq.BranchPattern,
	}

	if err := h.queries.CreateProjectEnvironment(c.Request().Context(), &environment); err != nil {
		return err
	}

	analytics.TrackUserProjectEvent(user, *project, posthog.Capture{
		Event: "project_environment_created",
	})

	return c.JSON(http.StatusCreated, typesHTTP.ProjectEnvironment{
		ID:            environment.ID,
		ProjectID:     environment.ProjectID,
		Name:          environment.Name,
		BranchPattern: environment.BranchPattern,
		UpdatedAt:     environment.UpdatedAt,
		CreatedAt:     environment.CreatedAt,
	})
}

func (h *Handler) UpdateProjectEnvironment(c echo.Context) error {
	id := c.Param("environment_id")

	environmentReq := typesHTTP.UpdateProjectEnvironmentBody{}

	if err := c.Bind(&environmentReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	validate := utilsValidator.NewValidator()
	if err := validate.Struct(environmentReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utilsValidator.ValidatorErrors(err))
	}

	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	environment, err := h.queries.GetProjectEnvironmentByID(c.Request().Context(), project, id)
	if err != nil {
		return err
	}

	environment.Name = environmentReq.Name
	environment.BranchPattern = environmentReq.BranchPattern

	if err := h.queries.UpdateProjectEnvironment(c.Request().Context(), environment); err != nil {
		return err
	}

	analytics.TrackUserProjectEvent(middleware.GetUser(c), *project, posthog.Capture{
		Event: "project_environment_updated",
	})

	return c.JSON(http.StatusOK, typesHTTP.ProjectEnvironment{
		ID:            environment.ID,
		ProjectID:     environment.ProjectID,
		Name:          environment.Name,
		BranchPattern: environment.BranchPattern,
		UpdatedAt:     environment.UpdatedAt,
		CreatedAt:     environment.CreatedAt,
	})
}

func (h *Handler) DeleteProjectEnvironment(c echo.Context) error {
	id := c.Param("environment_id")

	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	if err := h.queries.DeleteProjectEnvironmentByID(c.Request().Context(), *project, id); err != nil {
		return err
	}

	analytics.TrackUserProjectEvent(middleware.GetUser(c), *project, posthog.Capture{
		Event: "project_environment_deleted",
	})

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetProjectEnvironments(c echo.Context) error {

	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	environmentsDB, err := h.queries.GetProjectEnvironments(c.Request().Context(), project)
	if err != nil {
		return err
	}

	environments := make([]typesHTTP.ProjectEnvironment, len(environmentsDB))
	for i, environment := range environmentsDB {
		environments[i] = typesHTTP.ProjectEnvironment{
			ID:            environment.ID,
			ProjectID:     environment.ProjectID,
			Name:          environment.Name,
			BranchPattern: environment.BranchPattern,
			UpdatedAt:     environment.UpdatedAt,
			CreatedAt:     environment.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, environments)
}
