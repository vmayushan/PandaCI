package queriesWorkflow

import (
	"context"

	"github.com/pandaci-com/pandaci/pkg/utils"
	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
	nanoid "github.com/matoous/go-nanoid/v2"
)

func (q *WorkflowQueries) CreateStepRun(ctx context.Context, stepRun *typesDB.StepRun) error {

	var err error
	stepRun.ID, err = nanoid.New()
	if err != nil {
		return err
	}

	if len(stepRun.Name) > 255 {
		stepRun.Name = stepRun.Name[:252] + "..."
	}

	stepRun.CreatedAt = utils.CurrentTime()

	stepRunQuery := `INSERT INTO step_run
		(id, workflow_run_id, type, created_at, task_run_id, job_run_id, status, name)
	VALUES
		(:id, :workflow_run_id, :type, :created_at, :task_run_id, :job_run_id, :status, :name)`

	_, err = q.NamedExecContext(ctx, stepRunQuery, *stepRun)

	return err
}

func (q *WorkflowQueries) UpdateStepRunStatus(ctx context.Context, workflowID string, stepRunID string, status types.RunStatus) error {

	stepRunQuery := `UPDATE step_run SET
		status = $1
	WHERE id = $2 AND workflow_run_id = $3`

	_, err := q.ExecContext(ctx, stepRunQuery, status, stepRunID, workflowID)

	return err
}

func (q *WorkflowQueries) FinishStepRun(ctx context.Context, workflowID string, stepRunID string, conclusion types.RunConclusion) error {

	stepRunQuery := `UPDATE step_run SET
		status = 'completed',
		conclusion = $1,
		finished_at = NOW()
	WHERE id = $2 AND workflow_run_id = $3`

	_, err := q.ExecContext(ctx, stepRunQuery, conclusion, stepRunID, workflowID)

	return err
}
