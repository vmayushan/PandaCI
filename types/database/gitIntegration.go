package typesDB

import "github.com/pandaci-com/pandaci/types"

type GitIntegration struct {
	ID                string                `db:"id"`
	ProviderID        string                `db:"provider_id"`
	ProviderAccountID string                `db:"provider_account_id"`
	Type              types.GitProviderType `db:"type"`
}
