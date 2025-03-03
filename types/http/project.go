package typesHTTP

import (
	"time"

	"github.com/pandaci-com/pandaci/types"
)

type CreateProjectHTTP struct {
	OrgID                    string                `json:"orgID"`
	Slug                     string                `json:"slug"`
	Name                     string                `json:"name"`
	AvatarURL                string                `json:"avatarURL"`
	GitProviderIntegrationID string                `json:"gitProviderIntegrationID"`
	GitProviderRepoID        string                `json:"gitProviderRepoID"`
	GitProviderType          types.GitProviderType `json:"gitProviderType"`
}

type ProjectHTTP struct {
	ID        string     `json:"id"`
	OrgID     string     `json:"orgID"`
	Slug      string     `json:"slug"`
	Name      string     `json:"name"`
	AvatarURL string     `json:"avatarURL"`
	LastBuild *time.Time `json:"lastBuild"`
}

type UpdateProjectHTTP struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	AvatarURL string `json:"avatarURL"`
}
