package git

import (
	"fmt"

	gitGithub "github.com/pandaci-com/pandaci/app/git/github"
	gitShared "github.com/pandaci-com/pandaci/app/git/shared"
	"github.com/pandaci-com/pandaci/types"
)

func (h *Handler) GetClient(clientType types.GitProviderType) (gitShared.Client, error) {
	switch clientType {
	case types.GitProviderTypeGithub:
		return gitGithub.NewGithubClient(h.queries)
	default:
		return nil, fmt.Errorf("invalid git integration type: %s", clientType)
	}
}
