package queriesWorkflow

import (
	"context"
	"math"

	queries_utils "github.com/pandaci-com/pandaci/app/queries/utils"
	"github.com/pandaci-com/pandaci/pkg/utils"
	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
	"github.com/rs/zerolog/log"

	nanoid "github.com/matoous/go-nanoid/v2"
)

func (q *WorkflowQueries) CreateFailedWorkflowRun(ctx context.Context, workflowRun *typesDB.WorkflowRun) error {
	workflowRun.Status = types.RunStatusCompleted
	workflowRun.Conclusion = types.Pointer(types.RunConclusionFailure)
	workflowRun.FinishedAt = types.Pointer(utils.CurrentTime())

	return q.CreateWorkflowRun(ctx, workflowRun)
}

func (q *WorkflowQueries) CreateWorkflowRun(ctx context.Context, workflowRun *typesDB.WorkflowRun) error {

	for i := 0; i < 5; i++ {
		if err := q.createWorkflowRun(ctx, workflowRun); err != nil {
			// We can't lock the rows to count if we don't have any rows yet
			if utils.CheckConstraintError(err, "workflow_run_project_id_number_key") {
				log.Error().Str("project_id", workflowRun.ProjectID).Msg("Failed to create workflow run due to duplicate number")
				continue
			}
			return err
		}

		return nil
	}

	return fmt.Errorf("failed to create workflow run after 5 attempts due to duplicate project ID constraints")
}

func (q *WorkflowQueries) createWorkflowRun(ctx context.Context, workflowRun *typesDB.WorkflowRun) error {
	var err error
	if workflowRun.ID, err = nanoid.New(); err != nil {
		return err
	}

	// Ensures that alerts are formatted correctly
	workflowRun.GetAlerts()

	if workflowRun.GitTitle != nil && len(*workflowRun.GitTitle) > 255 {
		workflowRun.GitTitle = types.Pointer((*workflowRun.GitTitle)[:252] + "...")
	}

	workflowRun.CreatedAt = utils.CurrentTime()

	workflowRunQuery := `WITH locked_rows AS (
			SELECT number
			FROM workflow_run
			WHERE project_id = :project_id
			FOR UPDATE
		),
		next_number AS (
			SELECT COALESCE(MAX(number), 0) + 1 AS number FROM locked_rows
		)
		INSERT INTO workflow_run
			(id, project_id, name, status, conclusion, number, created_at, finished_at, git_title, git_sha, git_branch, pr_number, trigger, committer_email, user_id, alerts)
		VALUES
			(:id, :project_id, :name, :status, :conclusion, (SELECT number FROM next_number), :created_at, :finished_at, :git_title, :git_sha, :git_branch, :pr_number, :trigger, :committer_email, :user_id, :alerts)
		RETURNING number;`

	tx, err := q.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	completed := false
	defer queries_utils.Rollback(&completed, tx)

	rows, err := tx.NamedQuery(workflowRunQuery, workflowRun)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	completed = true

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&workflowRun.Number); err != nil {
			return err
		}
	}
	return nil
}

func (q *WorkflowQueries) UpdateWorkflowRunStatus(ctx context.Context, workflowID string, status types.RunStatus) error {
	workflowRunQuery := `UPDATE workflow_run SET
		status = $1
	WHERE id = $2`

	_, err := q.ExecContext(ctx, workflowRunQuery, status, workflowID)

	return err
}

func (q *WorkflowQueries) FailStuckRuns(ctx context.Context, workflowID string) error {
	// TODO - add billing to this

	jobQuery := `UPDATE job_run SET
		status = 'completed',
		conclusion = 'failure',
		finished_at = NOW()
	WHERE workflow_run_id = $1 AND NOT status = 'completed'`

	if _, err := q.ExecContext(ctx, jobQuery, workflowID); err != nil {
		return err
	}

	tasksQuery := `UPDATE task_run SET
		status = 'completed',
		conclusion = 'failure',
		finished_at = NOW()
	WHERE workflow_run_id = $1 AND NOT status = 'completed'`

	if _, err := q.ExecContext(ctx, tasksQuery, workflowID); err != nil {
		return err
	}

	stepsQuery := `UPDATE step_run SET
		status = 'completed',
		conclusion = 'failure',
		finished_at = NOW()
	WHERE workflow_run_id = $1 AND NOT status = 'completed'`

	if _, err := q.ExecContext(ctx, stepsQuery, workflowID); err != nil {
		return err
	}

	return nil
}

func (q *WorkflowQueries) UpdateWorkflowRun(ctx context.Context, workflowRun *typesDB.WorkflowRun) error {

	workflowRunQuery := `UPDATE workflow_run SET
		status = :status,
		conclusion = :conclusion,
		finished_at = :finished_at,
		build_minutes = :build_minutes
	WHERE id = :id`

	_, err := q.NamedExecContext(ctx, workflowRunQuery, workflowRun)

	return err
}

func (q *WorkflowQueries) AppendAlertToWorkflowRun(ctx context.Context, workflowRunID string, alert types.WorkflowRunAlert) error {
	tx, err := q.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	completed := false
	defer queries_utils.Rollback(&completed, tx)

	getWorkflowQuery := `SELECT
      id,
      project_id,
      name,
      status,
      conclusion,
      created_at,
      finished_at,
      number,
      git_sha,
      git_branch,
      runner,
   	  committer_email,
	  user_id,
      alerts,
      pr_number,
      trigger,
      git_title
	FROM workflow_run
	WHERE id = $1
	FOR UPDATE`

	var workflowRun typesDB.WorkflowRun
	if err := tx.GetContext(ctx, &workflowRun, getWorkflowQuery, workflowRunID); err != nil {
		return err
	}

	if err := typesDB.AppendAlert(&workflowRun, alert); err != nil {
		return err
	}

	updateQuery := `UPDATE workflow_run SET
		alerts = $1
	WHERE id = $2`

	if _, err := tx.ExecContext(ctx, updateQuery, workflowRun.Alerts, workflowRunID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	completed = true

	return nil
}

func (q *WorkflowQueries) FailWorkflowRun(ctx context.Context, workflowRun *typesDB.WorkflowRun, alert types.WorkflowRunAlert) error {

	// TOOD - we should do this in one transaction
	if err := q.AppendAlertToWorkflowRun(ctx, workflowRun.ID, alert); err != nil {
		return err
	}

	workflowRun.Status = types.RunStatusCompleted
	workflowRun.Conclusion = types.Pointer(types.RunConclusionFailure)
	workflowRun.FinishedAt = types.Pointer(utils.CurrentTime())
	workflowRun.BuildMinutes = int(math.Round(workflowRun.FinishedAt.Sub(workflowRun.CreatedAt).Minutes()+0.5)) * types.GetBuildMinutesScale(types.CloudRunner(workflowRun.Runner))

	return q.UpdateWorkflowRun(ctx, workflowRun)
}
