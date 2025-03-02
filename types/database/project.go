package typesDB

import "time"

type Project struct {
	ID                string `db:"id"`
	OrgID             string `db:"org_id"`
	Slug              string `db:"slug"`
	Name              string `db:"name"`
	GitIntegrationID  string `db:"git_integration_id"`
	GitProviderRepoID string `db:"git_provider_repo_id"`
	AvatarURL         string `db:"avatar_url"`
}

type ProjectWithLastBuild struct {
	ID                string     `db:"id"`
	OrgID             string     `db:"org_id"`
	Slug              string     `db:"slug"`
	Name              string     `db:"name"`
	GitIntegrationID  string     `db:"git_integration_id"`
	GitProviderRepoID string     `db:"git_provider_repo_id"`
	AvatarURL         string     `db:"avatar_url"`
	LastBuild         *time.Time `db:"last_build"`
}
