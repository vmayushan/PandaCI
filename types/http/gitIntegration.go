package typesHTTP

import "github.com/pandaci-com/pandaci/types"

type GitIntegration struct {
	ID                string                `json:"id"`
	ProviderID        string                `json:"provider_id"`
	ProviderAccountID string                `json:"provider_account_id"`
	Type              types.GitProviderType `json:"type"`
}
