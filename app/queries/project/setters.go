package projectQueries

import (
	"context"
	"fmt"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/alfiejones/panda-ci/pkg/utils"
	typesDB "github.com/alfiejones/panda-ci/types/database"
)

func (q *ProjectQueries) CreateProject(ctx context.Context, org *typesDB.OrgDB, project *typesDB.Project) error {
	if !utils.IsURLNameValid(project.Slug) {
		return fmt.Errorf("Invalid project url")
	}

	var err error
	if project.ID, err = nanoid.New(); err != nil {
		return err
	}

	projectQuery := `INSERT INTO project 
      (id, org_id, slug, name, git_integration_id, git_provider_repo_id, avatar_url)
    VALUES
  (:id, :org_id, :slug, :name, :git_integration_id, :git_provider_repo_id, :avatar_url)`

	_, err = q.NamedExecContext(ctx, projectQuery, project)

	return err
}

func (q *ProjectQueries) UpdateProject(ctx context.Context, project *typesDB.Project) error {
	if !utils.IsURLNameValid(project.Slug) {
		return fmt.Errorf("Invalid project url")
	}

	projectQuery := `UPDATE project
	SET name = :name, slug = :slug, avatar_url = :avatar_url
	WHERE id = :id`

	_, err := q.NamedExecContext(ctx, projectQuery, project)

	return err
}

func (q *ProjectQueries) DeleteProject(ctx context.Context, project *typesDB.Project) error {
	projectQuery := `DELETE FROM project WHERE id = :id`

	_, err := q.NamedExecContext(ctx, projectQuery, project)

	return err
}
