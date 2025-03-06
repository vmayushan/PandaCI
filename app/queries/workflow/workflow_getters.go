package queriesWorkflow

import (
	"context"
	"database/sql"
	"errors"

	typesDB "github.com/pandaci-com/pandaci/types/database"
	typesHTTP "github.com/pandaci-com/pandaci/types/http"
)

func (q *WorkflowQueries) Unsafe_GetWorkflowRunByID(ctx context.Context, runID string) (*typesDB.WorkflowRun, error) {
	var workflowRun typesDB.WorkflowRun

	query := `SELECT
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
	WHERE id = $1`

	return &workflowRun, q.GetContext(ctx, &workflowRun, query, runID)
}

func (q *WorkflowQueries) GetWorkflowRunByID(ctx context.Context, project *typesDB.Project, runID string) (*typesDB.WorkflowRun, error) {
	var workflowRun typesDB.WorkflowRun

	query := `SELECT
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
    WHERE id = $1 AND project_id = $2`

	return &workflowRun, q.GetContext(ctx, &workflowRun, query, runID, project.ID)
}

func (q *WorkflowQueries) GetWorkflowRunByNumber(ctx context.Context, project *typesDB.Project, number int) (*typesDB.WorkflowRun, error) {
	var workflowRun typesDB.WorkflowRun

	query := `SELECT
      id,
      project_id,
      name,
      status,
      conclusion,
      created_at,
      finished_at,
      number,
      runner,
      git_sha,
      git_branch,
   	  committer_email,
	  user_id,
      alerts,
      pr_number,
      trigger,
      git_title
    FROM workflow_run
    WHERE number = $1 AND project_id = $2`

	return &workflowRun, q.GetContext(ctx, &workflowRun, query, number, project.ID)
}

func (q *WorkflowQueries) GetProjectWorkflows(ctx context.Context, project *typesDB.Project, page int, perPage int) (*[]typesDB.WorkflowRun, bool, error) {
	var workflowRun []typesDB.WorkflowRun

	query := `SELECT
      id,
      project_id,
      name,
      status,
      conclusion,
      created_at,
      finished_at,
      git_sha,
      runner,
      git_branch,
   	  committer_email,
	  user_id,
      number,
      alerts,
      pr_number,
      trigger,
      git_title
    FROM workflow_run
    WHERE project_id = $1
    ORDER BY created_at DESC
    LIMIT $2 OFFSET $3`

	err := q.SelectContext(ctx, &workflowRun, query, project.ID, perPage+1, (page)*perPage)
	if err != nil {
		return nil, false, err
	}

	hasMore := false
	if len(workflowRun) > perPage {
		hasMore = true
		workflowRun = workflowRun[:len(workflowRun)-1]
	}

	return &workflowRun, hasMore, nil
}

func (q *WorkflowQueries) GetWorkflowRunJobsHTTP(ctx context.Context, workflow *typesDB.WorkflowRun) ([]typesHTTP.JobRun, error) {

	jobsDB, err := q.GetJobRunsByWorkflowID(ctx, workflow.ID, &GetJobRunsByWorkflowIDOptions{
		OrderByID: true,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) {
		jobsDB = &[]typesDB.JobRun{}
	}

	jobIDs := make([]string, len(*jobsDB))
	for i, job := range *jobsDB {
		jobIDs[i] = job.ID
	}

	if len(jobIDs) == 0 {
		return nil, nil
	}

	tasksDB, err := q.GetTaskRunsByJobRunIDs(ctx, jobIDs, &GetTaskRunsByJobRunIDsOptions{
		OrderByJobRunIDAndCreatedAt: true,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) {
		tasksDB = []typesDB.TaskRun{}
	}

	stepsDB, err := q.GetStepRunsByStepRunIDs(ctx, jobIDs, &GetStepRunsByJobRunIDsOptions{
		OrderByTaskRunIDAndCreatedAt: true,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if errors.Is(err, sql.ErrNoRows) {
		stepsDB = []typesDB.StepRun{}
	}

	jobs := make([]typesHTTP.JobRun, len(*jobsDB))

	stepPointer := 0
	taskPointer := 0

	for i, job := range *jobsDB {
		jobs[i] = typesHTTP.JobRun{
			ID:            job.ID,
			Number:        job.Number,
			WorkflowRunID: job.WorkflowRunID,
			Name:          job.Name,
			CreatedAt:     job.CreatedAt,
			FinishedAt:    job.FinishedAt,
			Status:        job.Status,
			Conclusion:    job.Conclusion,
			Runner:        job.Runner,
			Tasks:         []typesHTTP.TaskRun{},
		}

		for taskPointer < len(tasksDB) {

			task := tasksDB[taskPointer]

			if task.JobRunID != job.ID {
				break
			}

			taskHTTP := typesHTTP.TaskRun{
				ID:          task.ID,
				CreatedAt:   task.CreatedAt,
				FinishedAt:  task.FinishedAt,
				Name:        task.Name,
				Status:      task.Status,
				Conclusion:  task.Conclusion,
				DockerImage: task.DockerImage,
				Steps:       []typesHTTP.StepRun{},
			}

			for stepPointer < len(stepsDB) {
				step := (stepsDB)[stepPointer]

				if step.TaskRunID != task.ID {
					break
				}

				taskHTTP.Steps = append(taskHTTP.Steps, typesHTTP.StepRun{
					ID:         step.ID,
					Type:       step.Type,
					CreatedAt:  step.CreatedAt,
					FinishedAt: step.FinishedAt,
					Status:     step.Status,
					Conclusion: step.Conclusion,
					Name:       step.Name,
				})
				stepPointer++
			}

			// We want to order the tasks by created_at
			// we can't to this at the db level since we are ordering by id to ensure we have the right order for steps
			insertIndex := len(jobs[i].Tasks)
			for j, t := range jobs[i].Tasks {
				if t.CreatedAt.After(task.CreatedAt) {
					insertIndex = j
					break
				}
			}
			jobs[i].Tasks = append(jobs[i].Tasks[:insertIndex], append([]typesHTTP.TaskRun{taskHTTP}, jobs[i].Tasks[insertIndex:]...)...)

			taskPointer++
		}

	}

	return jobs, nil
}
