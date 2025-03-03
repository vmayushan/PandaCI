package queriesWorkflow

import (
	"context"
	"math"
	"time"

	"github.com/pandaci-com/pandaci/pkg/utils"
	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
	nanoid "github.com/matoous/go-nanoid/v2"
)

func (q *WorkflowQueries) CreateJobRun(ctx context.Context, jobRun *typesDB.JobRun) error {

	var err error
	if jobRun.ID, err = nanoid.New(); err != nil {
		return err
	}

	jobRun.CreatedAt = utils.CurrentTime()

	if jobRun.Conclusion != nil && *jobRun.Conclusion == types.RunConclusionSkipped {
		jobRun.FinishedAt = &jobRun.CreatedAt
	}

	jobRunQuery := `
		INSERT INTO job_run
			(id, workflow_run_id, name, status, conclusion, created_at, finished_at, runner, number)
		VALUES
			(:id, :workflow_run_id, :name, :status, :conclusion, :created_at, :finished_at, :runner,
			(SELECT COALESCE(MAX(number), 0) + 1 FROM job_run WHERE workflow_run_id = :workflow_run_id))
	`

	_, err = q.NamedExecContext(ctx, jobRunQuery, jobRun)

	// TODO - there's a chance we could create a job run with the same number if two jobs are created at the same time
	// we should retry the insert if we get a unique constraint error
	// This is probably very unlikely to happen, but it's better to be safe than sorry

	if utils.CheckConstraintError(err, "job_run_workflow_run_id_number_key") {
		return q.CreateJobRun(ctx, jobRun)
	}

	return err
}

func (q *WorkflowQueries) UpdateJobRunStatus(ctx context.Context, workflowID string, jobRunID string, status types.RunStatus) error {

	jobRunQuery := `UPDATE job_run SET
		status = $1
	WHERE id = $2 AND workflow_run_id = $3`

	_, err := q.ExecContext(ctx, jobRunQuery, status, jobRunID, workflowID)

	return err
}

func (q *WorkflowQueries) FinishJobRun(ctx context.Context, workflowID string, jobRunID string, conclusion types.RunConclusion) error {

	jobRun, err := q.GetJobRunByID(ctx, jobRunID)
	if err != nil {
		return err
	}

	buildMinutes := int(math.Round(time.Now().Sub(jobRun.CreatedAt).Minutes()+0.5)) * types.GetBuildMinutesScale(types.CloudRunner(jobRun.Runner))

	jobRunQuery := `UPDATE job_run SET
		status = 'completed',
		conclusion = $1,
		finished_at = NOW(),
		build_minutes = $2
	WHERE id = $3 AND workflow_run_id = $4`

	_, err = q.ExecContext(ctx, jobRunQuery, conclusion, buildMinutes, jobRunID, workflowID)

	return err
}

func (q *WorkflowQueries) FinishJobRunWithError(ctx context.Context, workflowID string, jobRunID string, msg string) error {

	jobRun, err := q.GetJobRunByID(ctx, jobRunID)

	if err != nil {
		return err
	}

	buildMinutes := int(math.Round(time.Now().Sub(jobRun.CreatedAt).Minutes()+0.5)) * types.GetBuildMinutesScale(types.CloudRunner(jobRun.Runner))

	jobRunQuery := `UPDATE job_run SET
		status = 'completed',
		conclusion = 'failure',
		finished_at = NOW(),
		build_minutes = $1
	WHERE id = $2 AND workflow_run_id = $3`

	_, err = q.ExecContext(ctx, jobRunQuery, buildMinutes, jobRunID, workflowID)

	return err
}
