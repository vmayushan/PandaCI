package gitGithub

import (
	gitShared "github.com/pandaci-com/pandaci/app/git/shared"
	"github.com/pandaci-com/pandaci/app/queries"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v68/github"
)

type GithubClient struct {
	queries *queries.Queries
}

type GithubUserClient struct {
	queries      *queries.Queries
	githubClient *github.Client
	gitClient    *GithubClient
}

type GithubAppClient struct {
	githubClient *github.Client
	gitClient    *GithubClient
}

type GithubInstallationClient struct {
	githubClient   *github.Client
	gitClient      *GithubClient
	installationID string
	graphqlClient  *graphqlClient
	itr            *ghinstallation.Transport
}

type graphqlClient struct {
	itr *ghinstallation.Transport
}

func NewGithubClient(queries *queries.Queries) (gitShared.Client, error) {
	return &GithubClient{
		queries: queries,
	}, nil
}
