package middlewareLoaders

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	typesDB "github.com/pandaci-com/pandaci/types/database"
)

func GetOrg(c echo.Context) (*typesDB.OrgDB, error) {
	org := c.Get("org")

	if org == nil {
		return nil, fmt.Errorf("org not found")
	}

	return org.(*typesDB.OrgDB), nil
}

func (m *middlewareHandler) LoadOrg(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		orgURLName := c.Param("org_url_name")
		orgID := c.Param("org_id")

		if orgURLName == "" && orgID == "" {
			return echo.NewHTTPError(http.StatusNotFound, "urlName or id required")
		}

		user := m.getCurrentUser(c)

		var org *typesDB.OrgDB
		var err error
		if orgID != "" {
			org, err = m.queries.GetOrgByIDAsUser(c.Request().Context(), user, orgID)
		} else {
			org, err = m.queries.GetOrgByURLNameAsUser(c.Request().Context(), user, orgURLName)
		}

		if err != nil {
			return err
		}

		c.Set("org", org)

		return next(c)
	}
}
