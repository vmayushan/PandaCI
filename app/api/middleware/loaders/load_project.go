package middlewareLoaders

import (
	"fmt"
	"net/http"

	typesDB "github.com/pandaci-com/pandaci/types/database"
	"github.com/labstack/echo/v4"
)

func GetProject(c echo.Context) (*typesDB.Project, error) {
	project := c.Get("project")

	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	return project.(*typesDB.Project), nil
}

func (m *middlewareHandler) LoadProject(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		projectURLName := c.Param("project_url_name")

		if projectURLName == "" {
			return echo.NewHTTPError(http.StatusNotFound, "urlName or id required")
		}

		org, err := GetOrg(c)
		if err != nil {
			return err
		}

		project, err := m.queries.GetOrgProjectByName(c.Request().Context(), org, projectURLName)
		if err != nil {
			return err
		}

		c.Set("project", project)

		return next(c)
	}
}
