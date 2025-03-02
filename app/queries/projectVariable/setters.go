package projectVariableQueries

import (
	"context"

	"github.com/alfiejones/panda-ci/pkg/utils"
	typesDB "github.com/alfiejones/panda-ci/types/database"
	nanoid "github.com/matoous/go-nanoid/v2"
)

// Does not enforce one variable name per environment
func (q *ProjectVariableQueries) CreateProjectVariable(ctx context.Context, variable *typesDB.ProjectVariable, environmentIDs []string) error {
	var err error
	if variable.ID, err = nanoid.New(); err != nil {
		return err
	}


	variable.CreatedAt = utils.CurrentTime()
	variable.UpdatedAt = variable.CreatedAt

	projectVariableQuery := `INSERT INTO project_variable
		(id, project_id, key, value, updated_at, created_at, encryption_key_id, initialisation_vector, sensitive)
	VALUES
		(:id, :project_id, :key, :value, :updated_at, :created_at, :encryption_key_id, :initialisation_vector, :sensitive)`

	projectVariableOnEnvironmentQuery := `INSERT INTO project_variable_on_project_environment
		(project_variable_id, project_environment_id)
	VALUES
		(:project_variable_id, :project_environment_id)`

	tx, err := q.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.NamedExecContext(ctx, projectVariableQuery, variable)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, environmentID := range environmentIDs {
		_, err = tx.NamedExecContext(ctx, projectVariableOnEnvironmentQuery, typesDB.ProjectVariableOnProjectEnvironment{
			ProjectVariableID:    variable.ID,
			ProjectEnvironmentID: environmentID,
		})
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (q ProjectVariableQueries) DeleteProjectVariableByID(ctx context.Context, project typesDB.Project, id string) error {
	projectVariableQuery := `DELETE FROM project_variable WHERE id = $1 AND project_id = $2`

	_, err := q.ExecContext(ctx, projectVariableQuery, id, project.ID)

	return err
}
