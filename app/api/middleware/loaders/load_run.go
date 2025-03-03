package middlewareLoaders

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	typesDB "github.com/pandaci-com/pandaci/types/database"
)

func GetWorkflowRun(c echo.Context) (*typesDB.WorkflowRun, error) {
	workflowRun := c.Get("workflow_run")

	if workflowRun == nil {
		return nil, fmt.Errorf("workflow run not found")
	}

	return workflowRun.(*typesDB.WorkflowRun), nil
}

func (m *middlewareHandler) LoadWorkflowRun(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		workflowRunNumber := c.Param("run_number")
		workflowRunID := c.Param("run_id")

		if workflowRunID == "" && workflowRunNumber == "" {
			return echo.NewHTTPError(http.StatusNotFound, "run number or id required")
		}

		project, err := GetProject(c)
		if err != nil {
			return err
		}

		var workflowRun *typesDB.WorkflowRun
		if workflowRunNumber != "" {
			workflowRunNumberInt, err := strconv.Atoi(workflowRunNumber)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid run number")
			}
			workflowRun, err = m.queries.GetWorkflowRunByNumber(c.Request().Context(), project, workflowRunNumberInt)
		} else {
			workflowRun, err = m.queries.GetWorkflowRunByID(c.Request().Context(), project, workflowRunID)
		}

		if err != nil {
			return err
		}

		c.Set("workflow_run", workflowRun)

		return next(c)
	}
}
