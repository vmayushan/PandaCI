package queriesWorkflow

import (
	"context"

	"github.com/jmoiron/sqlx"
	typesDB "github.com/alfiejones/panda-ci/types/database"
)

type GetStepRunsByJobRunIDsOptions struct {
	OrderByTaskRunIDAndCreatedAt bool
}

func (q *WorkflowQueries) GetStepRunsByStepRunIDs(ctx context.Context, jobRunIDs []string, opts *GetStepRunsByJobRunIDsOptions) ([]typesDB.StepRun, error) {
	var stepRuns []typesDB.StepRun

	if len(jobRunIDs) == 0 {
		return stepRuns, nil
	}

	stepQuery := `SELECT
			id,
			type,
			job_run_id,
			workflow_run_id,
			name,
			created_at,
			finished_at	,
			task_run_id,
			status,
			conclusion
		FROM step_run
		WHERE job_run_id IN (?)`

	if opts != nil {
		if opts.OrderByTaskRunIDAndCreatedAt {
			stepQuery += " ORDER BY job_run_id ASC, task_run_id ASC, created_at ASC"
		}
	}

	query, args, err := sqlx.In(stepQuery, jobRunIDs)
	if err != nil {
		return nil, err
	}

	query = q.DB.Rebind(query)

	if err := q.SelectContext(ctx, &stepRuns, query, args...); err != nil {
		return nil, err
	}

	return stepRuns, nil
}
