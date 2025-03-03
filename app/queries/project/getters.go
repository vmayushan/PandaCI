package projectQueries

import (
	"context"

	typesDB "github.com/pandaci-com/pandaci/types/database"
)

func (q *ProjectQueries) CountOrgProjects(ctx context.Context, org *typesDB.OrgDB) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM project WHERE org_id = $1`

	return count, q.GetContext(ctx, &count, query, org.ID)
}

func (q *ProjectQueries) GetOrgProjects(ctx context.Context, org *typesDB.OrgDB) (*[]typesDB.ProjectWithLastBuild, error) {
	var projects []typesDB.ProjectWithLastBuild

	query := `SELECT
      id,
      org_id,
      slug,
      name,
      git_integration_id,
      git_provider_repo_id,
      avatar_url,
      (SELECT MAX(created_at) FROM workflow_run WHERE workflow_run.project_id = project.id) AS last_build
    FROM project
    WHERE org_id = $1
    ORDER BY last_build DESC NULLS LAST`

	return &projects, q.SelectContext(ctx, &projects, query, org.ID)
}

func (q *ProjectQueries) Unsafe_GetProjectByID(ctx context.Context, projectID string) (*typesDB.Project, error) {
	var project typesDB.Project

	query := `SELECT
	  project.id,
	  project.org_id,
	  project.slug,
	  project.name,
	  project.git_integration_id,
	  project.git_provider_repo_id,
	  project.avatar_url
	FROM project
	WHERE project.id = $1`

	return &project, q.GetContext(ctx, &project, query, projectID)
}

func (q *ProjectQueries) GetOrgProjectByID(ctx context.Context, org *typesDB.OrgDB, projectID string) (*typesDB.Project, error) {
	var project typesDB.Project

	query := `SELECT
      project.id,
      project.org_id,
      project.slug,
      project.name,
      project.git_integration_id,
      project.git_provider_repo_id,
      project.avatar_url
    FROM project
    WHERE project.org_id = $1 AND project.id = $2`

	return &project, q.GetContext(ctx, &project, query, org.ID, projectID)
}

func (q *ProjectQueries) GetOrgProjectByName(ctx context.Context, org *typesDB.OrgDB, projectURLName string) (*typesDB.Project, error) {
	var project typesDB.Project

	query := `SELECT
      id,
      org_id,
      slug,
      name,
      git_integration_id,
      git_provider_repo_id,
      avatar_url
    FROM project
    WHERE org_id = $1 AND slug = $2`

	return &project, q.GetContext(ctx, &project, query, org.ID, projectURLName)
}

func (q *ProjectQueries) GetProjectsByGitIntegrationID(ctx context.Context, gitIntegrationID string, repoID string) (*[]typesDB.Project, error) {
	var projects []typesDB.Project

	query := `SELECT
	  id,
	  org_id,
	  slug,
	  name,
	  git_integration_id,
	  git_provider_repo_id,
	  avatar_url
	FROM project
	WHERE git_integration_id = $1 AND git_provider_repo_id = $2`

	return &projects, q.SelectContext(ctx, &projects, query, gitIntegrationID, repoID)
}
