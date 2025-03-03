package queriesWorkflow

import (
	"context"

	"github.com/jmoiron/sqlx"
	typesDB "github.com/pandaci-com/pandaci/types/database"
)

type GetTaskRunsByJobRunIDsOptions struct {
	OrderByJobRunIDAndCreatedAt bool
	OrderByCreatedAt            bool
}

func (q *WorkflowQueries) GetTaskRunsByJobRunIDs(ctx context.Context, jobRunIDs []string, opts *GetTaskRunsByJobRunIDsOptions) ([]typesDB.TaskRun, error) {
	var taskRuns []typesDB.TaskRun

	if len(jobRunIDs) == 0 {
		return taskRuns, nil
	}

	taskQuery := `SELECT
			id,
			job_run_id,
			workflow_run_id,
			name,
			created_at,
			finished_at,
			status,
			conclusion,
			docker_image
		FROM task_run
		WHERE job_run_id IN (?)`

	if opts != nil {
		if opts.OrderByJobRunIDAndCreatedAt {
			taskQuery += " ORDER BY job_run_id ASC, id ASC"
		}
	}

	query, args, err := sqlx.In(taskQuery, jobRunIDs)
	if err != nil {
		return nil, err
	}

	query = q.DB.Rebind(query)

	if err := q.SelectContext(ctx, &taskRuns, query, args...); err != nil {
		return nil, err
	}

	return taskRuns, nil
}
