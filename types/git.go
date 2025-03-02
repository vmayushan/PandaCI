package types

type GitProviderType string

type GitRepoData struct {
	URL           string
	DefaultBranch string
}

const (
	GitProviderTypeGithub GitProviderType = "github"
)

type GitCloneOptions struct {
	URL        string
	Sha        string
	FetchDepth int
}
