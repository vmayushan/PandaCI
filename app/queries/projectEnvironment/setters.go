package queriesProjectEnvironment

import (
	"context"

	"github.com/alfiejones/panda-ci/pkg/utils"
	typesDB "github.com/alfiejones/panda-ci/types/database"
	nanoid "github.com/matoous/go-nanoid/v2"
)

func (q *ProjectEnvironmentQueries) UpdateProjectEnvironment(ctx context.Context, environment *typesDB.ProjectEnvironment) error {
	environment.UpdatedAt = utils.CurrentTime()

	projectEnvironmentQuery := `UPDATE project_environment
		SET
			name = :name,
			branch_pattern = :branch_pattern,
			updated_at = :updated_at
		WHERE id = :id AND project_id = :project_id`

	_, err := q.NamedExecContext(ctx, projectEnvironmentQuery, environment)

	return err
}

func (q *ProjectEnvironmentQueries) CreateProjectEnvironment(ctx context.Context, environment *typesDB.ProjectEnvironment) error {
	var err error
	if environment.ID, err = nanoid.New(); err != nil {
		return err
	}

	environment.CreatedAt = utils.CurrentTime()
	environment.UpdatedAt = environment.CreatedAt

	projectEnvironmentQuery := `INSERT INTO project_environment
		(id, project_id, name, branch_pattern, updated_at, created_at)
	VALUES
		(:id, :project_id, :name, :branch_pattern, :updated_at, :created_at)`

	_, err = q.NamedExecContext(ctx, projectEnvironmentQuery, environment)

	return err
}

func (q ProjectEnvironmentQueries) DeleteProjectEnvironmentByID(ctx context.Context, project typesDB.Project, id string) error {
	projectEnvironmentQuery := `DELETE FROM project_environment WHERE id = $1 AND project_id = $2`

	_, err := q.ExecContext(ctx, projectEnvironmentQuery, id, project.ID)

	return err
}
