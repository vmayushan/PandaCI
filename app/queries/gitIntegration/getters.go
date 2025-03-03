package gitIntegrationQueries

import (
	"context"

	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
)

func (q *GitIntegrationQueries) GetGitIntegration(ctx context.Context, id string) (*typesDB.GitIntegration, error) {
	query := `SELECT
		id,
		type,
		provider_account_id,
		provider_id
	FROM  git_integration
	WHERE  id = $1`

	var integration typesDB.GitIntegration

	return &integration, q.GetContext(ctx, &integration, query, id)
}

func (q *GitIntegrationQueries) GetGitIntegrationByProviderID(ctx context.Context, providerID string, providerType types.GitProviderType) (*typesDB.GitIntegration, error) {
	query := `SELECT
		id,
		type,
		provider_account_id,
		provider_id
	FROM git_integration
	WHERE provider_id = $1 AND type = $2`

	var integration typesDB.GitIntegration

	return &integration, q.GetContext(ctx, &integration, query, providerID, providerType)
}
