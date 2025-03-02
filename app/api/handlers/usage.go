package handlers

import (
	"net/http"

	middlewareLoaders "github.com/alfiejones/panda-ci/app/api/middleware/loaders"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetOrgUsage(c echo.Context) error {
	org, err := middlewareLoaders.GetOrg(c)
	if err != nil {
		return err
	}

	type Usage struct {
		ProjectCount     int `json:"projectCount"`
		UsedBuildMinutes int `json:"usedBuildMinutes"`
		UsedCommitters   int `json:"usedCommitters"`
	}

	projectCount, err := h.queries.CountOrgProjects(c.Request().Context(), org)
	if err != nil {
		return err
	}

	license, err := org.GetLicense()
	if err != nil {
		return err
	}

	billingPeriod := license.GetBillingPeriod()

	committersCount, err := h.queries.CountCommitters(c.Request().Context(), org.ID, billingPeriod.StartsAt, billingPeriod.EndsAt)
	if err != nil {
		return err
	}

	buildMinutes, err := h.queries.GetBuildMinutes(c.Request().Context(), org.ID, billingPeriod.StartsAt, billingPeriod.EndsAt)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Usage{
		ProjectCount:     projectCount,
		UsedBuildMinutes: buildMinutes,
		UsedCommitters:   committersCount,
	})
}
