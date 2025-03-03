package queriesProjectEnvironment

import (
	"context"

	typesDB "github.com/alfiejones/panda-ci/types/database"
)

func (q *ProjectEnvironmentQueries) GetProjectEnvironmentByID(ctx context.Context, project *typesDB.Project, id string) (*typesDB.ProjectEnvironment, error) {
	var projectEnvironment typesDB.ProjectEnvironment

	query := `SELECT
			id,
			project_id,
			name,
			branch_pattern,
			updated_at,
			created_at
		FROM project_environment
		WHERE id = $1 AND project_id = $2`

	err := q.GetContext(ctx, &projectEnvironment, query, id, project.ID)
	if err != nil {
		return nil, err
	}

	return &projectEnvironment, nil
}

func (q *ProjectEnvironmentQueries) GetProjectEnvironments(ctx context.Context, project *typesDB.Project) ([]typesDB.ProjectEnvironment, error) {
	var projectEnvironments []typesDB.ProjectEnvironment

	query := `SELECT
			id,
			project_id,
			name,
			branch_pattern,
			updated_at,
			created_at
		FROM project_environment
		WHERE project_id = $1`

	return projectEnvironments, q.SelectContext(ctx, &projectEnvironments, query, project.ID)
}
