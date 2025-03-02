package handlersProject

import (
	"net/http"

	"github.com/alfiejones/panda-ci/app/api/middleware"
	middlewareLoaders "github.com/alfiejones/panda-ci/app/api/middleware/loaders"
	"github.com/alfiejones/panda-ci/platform/analytics"
	"github.com/alfiejones/panda-ci/types"
	typesHTTP "github.com/alfiejones/panda-ci/types/http"
	"github.com/labstack/echo/v4"
	"github.com/posthog/posthog-go"
)

func (h *Handler) TriggerRun(c echo.Context) error {
	user := middleware.GetUser(c)

	triggerRequest := typesHTTP.TriggerRunRequest{}

	if err := c.Bind(&triggerRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	project, err := middlewareLoaders.GetProject(c)
	if err != nil {
		return err
	}

	triggerEvent := types.TriggerEvent{
		Source:   types.TriggerEventSourceManual,
		SHA:      triggerRequest.SHA,
		Branch:   triggerRequest.Branch,
		Trigger:  types.RunTriggerManual,
		GitTitle: "Manual Trigger",
		Committer: types.Committer{
			Email:  &user.Email,
			UserID: &user.ID,
		},
	}

	workflows, err := h.scanner.GetWorkflowDefinitions(c.Request().Context(), *project, triggerEvent)
	if err != nil {
		return err
	}

	runsDB, err := h.orchestrator.StartWorkflows(c.Request().Context(), project, workflows)
	if err != nil {
		return err
	}

	runs := make([]typesHTTP.WorkflowRun, len(runsDB))
	for i, run := range runsDB {
		runs[i] = typesHTTP.WorkflowRun{
			ID:        run.ID,
			ProjectID: run.ProjectID,
			Name:      run.Name,
			CreatedAt: run.CreatedAt,
			Status:    run.Status,
			Number:    run.Number,
			GitSha:    run.GitSha,
			GitBranch: run.GitBranch,
			PrNumber:  run.PrNumber,
			Trigger:   run.Trigger,
			Committer: typesHTTP.Committer{
				Email:  user.Email,
				Name:   *user.Name,
				Avatar: *user.Avatar,
			},
		}
	}

	analytics.TrackUserProjectEvent(user, *project, posthog.Capture{
		Event: "project_run_triggered",
	})

	return c.JSON(http.StatusAccepted, runs)
}
