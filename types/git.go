package types

import "fmt"

type GitProviderType string

type GitRepoData struct {
	GitURL        string
	HTMLURL       string
	DefaultBranch string
	Type          GitProviderType
}

func (g *GitRepoData) GetPRURL(prNumber *int32) *string {
	if prNumber == nil {
		return nil
	}

	if g.Type == GitProviderTypeGithub {
		return Pointer(fmt.Sprintf("%s/pull/%d", g.HTMLURL, *prNumber))
	}

	return nil
}

func (g *GitRepoData) GetCommitURL(sha string, prNumber *int32) *string {

	if g.Type == GitProviderTypeGithub {
		if prNumber != nil {
			// return Pointer(g.HTMLURL + "/pull/" + string(*prNumber) + "/commits/" + sha)
			return Pointer(fmt.Sprintf("%s/pull/%d/commits/%s", g.HTMLURL, *prNumber, sha))
		} else {
			return Pointer(fmt.Sprintf("%s/commit/%s", g.HTMLURL, sha))
		}
	}
	return nil
}

const (
	GitProviderTypeGithub GitProviderType = "github"
)

type GitCloneOptions struct {
	URL        string
	Sha        string
	FetchDepth int
}
