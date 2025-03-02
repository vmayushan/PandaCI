package projectVariableQueries

import (
	"context"
	"database/sql"
	"slices"
	"strings"
	"time"

	typesDB "github.com/alfiejones/panda-ci/types/database"
)

func (q *ProjectVariableQueries) GetProjectVariableByID(ctx context.Context, project *typesDB.Project, id string) (*typesDB.ProjectVariable, error) {
	var projectVariable typesDB.ProjectVariable

	query := `SELECT
			id,
			project_id,
			key,
			value,
			updated_at,
			created_at,
			encryption_key_id,
			initialisation_vector,
			sensitive
		FROM project_variable
		WHERE id = $1 AND project_id = $2`

	err := q.GetContext(ctx, &projectVariable, query, id, project.ID)
	if err != nil {
		return nil, err
	}

	return &projectVariable, nil
}

func (q *ProjectVariableQueries) GetProjectVariables(ctx context.Context, project *typesDB.Project) ([]typesDB.ProjectVariable, error) {
	var projectVariables []typesDB.ProjectVariable

	query := `SELECT
			id,
			project_id,
			key,
			value,
			updated_at,
			created_at,
			encryption_key_id,
			initialisation_vector,
			sensitive
		FROM project_variable
		WHERE project_id = $1`

	return projectVariables, q.SelectContext(ctx, &projectVariables, query, project.ID)
}

func (q *ProjectVariableQueries) GetProjectVariablesWithEnvironments(ctx context.Context, project *typesDB.Project) ([]typesDB.ProjectVariableWithEnvironments, error) {
	var projectVariables []typesDB.ProjectVariableWithEnvironments

	type RawValues struct {
		ID                       string         `db:"id"`
		ProjectID                string         `db:"project_id"`
		Key                      string         `db:"key"`
		Value                    string         `db:"value"`
		UpdatedAt                time.Time      `db:"updated_at"`
		CreatedAt                time.Time      `db:"created_at"`
		EncryptionKeyID          string         `db:"encryption_key_id"`
		InitialisationVector     string         `db:"initialisation_vector"`
		Sensitive                bool           `db:"sensitive"`
		EnvironmentID            sql.NullString `db:"environment_id"`
		EnvironmentName          sql.NullString `db:"environment_name"`
		EnvironmentBranchPattern sql.NullString `db:"environment_branch_pattern"`
		EnvironmentCreatedAt     sql.NullTime   `db:"environment_created_at"`
		EnvironmentUpdatedAt     sql.NullTime   `db:"environment_updated_at"`
	}

	query := `SELECT
			pv.id,
			pv.project_id,
			pv.key,
			pv.value,
			pv.updated_at,
			pv.created_at,
			pv.encryption_key_id,
			pv.initialisation_vector,
			pv.sensitive,
			pe.id AS environment_id,
			pe.name AS environment_name,
			pe.branch_pattern AS environment_branch_pattern,
			pe.created_at AS environment_created_at,
			pe.updated_at AS environment_updated_at
		FROM project_variable pv
		LEFT JOIN project_variable_on_project_environment pvope ON pv.id = pvope.project_variable_id
		LEFT JOIN project_environment pe ON pvope.project_environment_id = pe.id
		WHERE pv.project_id = $1
		ORDER BY pv.id, pe.name`

	var rawValues []RawValues

	err := q.SelectContext(ctx, &rawValues, query, project.ID)
	if err != nil {
		return nil, err
	}

	for i, rawValue := range rawValues {

		if i == 0 || projectVariables[len(projectVariables)-1].ID != rawValue.ID {
			projectVariables = append(projectVariables, typesDB.ProjectVariableWithEnvironments{
				ID:                   rawValue.ID,
				ProjectID:            rawValue.ProjectID,
				Key:                  rawValue.Key,
				Value:                rawValue.Value,
				UpdatedAt:            rawValue.UpdatedAt,
				CreatedAt:            rawValue.CreatedAt,
				EncryptionKeyID:      rawValue.EncryptionKeyID,
				InitialisationVector: rawValue.InitialisationVector,
				Sensitive:            rawValue.Sensitive,
				Environments:         []typesDB.ProjectEnvironment{},
			})
		}

		if rawValue.EnvironmentID.Valid {
			projectVariables[len(projectVariables)-1].Environments = append(projectVariables[len(projectVariables)-1].Environments, typesDB.ProjectEnvironment{
				ID:            rawValue.EnvironmentID.String,
				Name:          rawValue.EnvironmentName.String,
				BranchPattern: rawValue.EnvironmentBranchPattern.String,
				ProjectID:     project.ID,
				UpdatedAt:     rawValue.EnvironmentCreatedAt.Time,
				CreatedAt:     rawValue.EnvironmentCreatedAt.Time,
			})
		}
	}

	// Sort by environment name, allowing conflicting enviroments to be resolved consistently
	// in the future we'll probably want to allow
	slices.SortFunc(projectVariables, func(a, b typesDB.ProjectVariableWithEnvironments) int {
		if len(a.Environments) == 0 && len(b.Environments) == 0 {
			return 0
		} else if len(a.Environments) == 0 {
			return 1
		} else if len(b.Environments) == 0 {
			return -1
		}
		return strings.Compare(a.Environments[0].Name, b.Environments[0].Name)
	})

	return projectVariables, nil
}
