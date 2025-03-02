package queriesProjectEnvironment

import (
	"context"

	typesDB "github.com/alfiejones/panda-ci/types/database"
)

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
