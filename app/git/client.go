package git

import (
	"fmt"

	gitGithub "github.com/alfiejones/panda-ci/app/git/github"
	gitShared "github.com/alfiejones/panda-ci/app/git/shared"
	"github.com/alfiejones/panda-ci/types"
)

func (h *Handler) GetClient(clientType types.GitProviderType) (gitShared.Client, error) {
	switch clientType {
	case types.GitProviderTypeGithub:
		return gitGithub.NewGithubClient(h.queries)
	default:
		return nil, fmt.Errorf("invalid git integration type: %s", clientType)
	}
}
