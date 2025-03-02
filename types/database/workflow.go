package typesDB

import (
	"encoding/json"
	"time"

	utilsValidator "github.com/alfiejones/panda-ci/pkg/utils/validator"
	"github.com/alfiejones/panda-ci/types"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	"github.com/rs/zerolog/log"
)

type WorkflowRun struct {
	ID        string `db:"id"`
	Number    int    `db:"number"`
	ProjectID string `db:"project_id"`

	Runner string `db:"runner"`

	BuildMinutes int `db:"build_minutes"`

	Name       string               `db:"name"`
	Status     types.RunStatus      `db:"status"`
	Conclusion *types.RunConclusion `db:"conclusion"`

	CreatedAt  time.Time  `db:"created_at"`
	FinishedAt *time.Time `db:"finished_at"`

	GitSha    string `db:"git_sha"`
	GitBranch string `db:"git_branch"`

	GitTitle *string `db:"git_title"`

	CommitterEmail *string `db:"committer_email"` // The email of the user who triggered the run (if known)
	UserID         *string `db:"user_id"`         // The user who triggered the run (if known)

	PrNumber *int32 `db:"pr_number"`

	Trigger types.RunTrigger `db:"trigger"`

	Alerts sqlxTypes.JSONText `db:"alerts"`
}

func (wr WorkflowRun) GetAlerts() []types.WorkflowRunAlert {
	var alerts []types.WorkflowRunAlert

	if err := wr.Alerts.Unmarshal(&alerts); err != nil {
		log.Debug().Str("err", err.Error()).Msg("Failed to unmarshal alerts")
		return []types.WorkflowRunAlert{}
	}

	return alerts
}

func (wr *WorkflowRun) AppendAlert(alert types.WorkflowRunAlert) error {
	alerts := wr.GetAlerts()

	alerts = append(alerts, alert)

	return wr.SetAlerts(alerts)
}

func (wr *WorkflowRun) SetAlerts(alerts []types.WorkflowRunAlert) error {

	validator := utilsValidator.NewValidator()

	validAlerts := []types.WorkflowRunAlert{}

	for _, alert := range alerts {

		if len(alert.Message) > 512 {
			alert.Message = alert.Message[:509] + "..."
		}

		if len(alert.Title) > 256 {
			alert.Title = alert.Title[:253] + "..."
		}

		if err := validator.Struct(alert); err != nil {
			log.Err(err).Msg("Failed to validate alert")
			validAlerts = append(alerts, types.WorkflowRunAlert{
				Type:    types.WorkflowRunAlertTypeError,
				Title:   "Failed to validate alert",
				Message: err.Error(),
			})
		} else {
			validAlerts = append(validAlerts, alert)
		}
	}

	rawBytes, err := json.Marshal(validAlerts)
	if err != nil {
		log.Err(err).Msg("Failed to marshal alerts")
		return err
	}

	return wr.Alerts.Scan(rawBytes)
}

type JobRun struct {
	ID            string `db:"id"`
	Number        int    `db:"number"`
	WorkflowRunID string `db:"workflow_run_id"`

	Runner string `db:"runner"`

	BuildMinutes int `db:"build_minutes"`

	CreatedAt  time.Time  `db:"created_at"`
	FinishedAt *time.Time `db:"finished_at"`

	Name string `db:"name"`

	Status     types.RunStatus      `db:"status"`
	Conclusion *types.RunConclusion `db:"conclusion"`
}

type TaskRun struct {
	ID            string `db:"id"`
	JobRunID      string `db:"job_run_id"`
	WorkflowRunID string `db:"workflow_run_id"`

	Name string `db:"name"`

	CreatedAt  time.Time  `db:"created_at"`
	FinishedAt *time.Time `db:"finished_at"`

	Status     types.RunStatus      `db:"status"`
	Conclusion *types.RunConclusion `db:"conclusion"`

	DockerImage *string `db:"docker_image"`
}

type StepRun struct {
	ID            string            `db:"id"`
	Type          types.StepRunType `db:"type"`
	WorkflowRunID string            `db:"workflow_run_id"`

	CreatedAt  time.Time  `db:"created_at"`
	FinishedAt *time.Time `db:"finished_at"`

	TaskRunID string `db:"task_run_id"`
	JobRunID  string `db:"job_run_id"`

	Status     types.RunStatus      `db:"status"`
	Conclusion *types.RunConclusion `db:"conclusion"`

	Name string `db:"name"`
}
