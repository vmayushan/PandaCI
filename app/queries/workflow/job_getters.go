package queriesWorkflow

import (
	"context"

	typesDB "github.com/pandaci-com/pandaci/types/database"
)

type GetJobRunsByWorkflowIDOptions struct {
	OrderByID bool
}

func (q *WorkflowQueries) GetJobRunsByWorkflowID(ctx context.Context, workflowID string, opts *GetJobRunsByWorkflowIDOptions) (*[]typesDB.JobRun, error) {
	var jobRuns []typesDB.JobRun

	query := `SELECT
			id,
			number,
			workflow_run_id,
			name,
			status,
			conclusion,
			created_at,
			finished_at,
			runner
		FROM job_run
		WHERE workflow_run_id = $1`

	if opts != nil {
		if opts.OrderByID {
			query += " ORDER BY id"
		}
	}

	return &jobRuns, q.SelectContext(ctx, &jobRuns, query, workflowID)
}

func (q *WorkflowQueries) GetJobRunByID(ctx context.Context, jobRunID string) (*typesDB.JobRun, error) {
	var jobRun typesDB.JobRun

	query := `SELECT
			id,
			number,
			workflow_run_id,
			name,
			status,
			conclusion,
			created_at,
			finished_at,
			runner
		FROM job_run
		WHERE id = $1`

	return &jobRun, q.GetContext(ctx, &jobRun, query, jobRunID)
}
