package typesHTTP

import (
	"time"

	"github.com/alfiejones/panda-ci/types"
)

type GitRepository struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Description string    `json:"description"`
	Public      bool      `json:"public"`
	URL         string    `json:"url"`
}

type GitRepositories struct {
	Repos         []GitRepository `json:"repos"`
	LimitExceeded bool            `json:"limitExceeded"`
	Limit         int             `json:"limit"`
}

type GitInstallation struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	AvatarURL        string                `json:"avatarURL"`
	IsUser           bool                  `json:"isUser"`
	RepositoryScopes string                `json:"repositoryScopes"`
	AccountID        string                `json:"accountID"`
	Type             types.GitProviderType `json:"type"`
}

type GitInstallations struct {
	Installations []GitInstallation `json:"installations"`
	IsLastPage    bool              `json:"isLastPage"`
}
