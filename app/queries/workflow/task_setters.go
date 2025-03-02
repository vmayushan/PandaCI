package queriesWorkflow

import (
	"context"

	"github.com/alfiejones/panda-ci/pkg/utils"
	"github.com/alfiejones/panda-ci/types"
	typesDB "github.com/alfiejones/panda-ci/types/database"

	nanoid "github.com/matoous/go-nanoid/v2"
)

func (q *WorkflowQueries) CreateTaskRun(ctx context.Context, task *typesDB.TaskRun) error {

	id, err := nanoid.New()
	if err != nil {
		return err
	}

	task.ID = id

	task.CreatedAt = utils.CurrentTime()

	if task.Conclusion != nil && *task.Conclusion == types.RunConclusionSkipped {
		task.FinishedAt = &task.CreatedAt
	}

	taskRunQuery := `
		INSERT INTO task_run
			(id, workflow_run_id, job_run_id, name, created_at, finished_at, status, conclusion, docker_image)
		VALUES
			(:id, :workflow_run_id, :job_run_id, :name, :created_at, :finished_at, :status, :conclusion, :docker_image)
	`

	_, err = q.NamedExecContext(ctx, taskRunQuery, task)

	return err

}

func (q *WorkflowQueries) FinishTaskRun(ctx context.Context, workflowID string, taskRunID string, conclusion types.RunConclusion) error {

	taskRunQuery := `UPDATE task_run SET
		status = 'completed',
		conclusion = $1,
		finished_at = NOW()
	WHERE id = $2 AND workflow_run_id = $3`

	_, err := q.ExecContext(ctx, taskRunQuery, conclusion, taskRunID, workflowID)

	return err
}
