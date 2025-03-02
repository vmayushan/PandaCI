package queriesProjectEnvironment

import (
	"context"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/alfiejones/panda-ci/pkg/utils"
	typesDB "github.com/alfiejones/panda-ci/types/database"
)

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
