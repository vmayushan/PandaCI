package identity

import (
	"context"
	"fmt"

	"github.com/alfiejones/panda-ci/pkg/utils/env"
	"github.com/alfiejones/panda-ci/types"
	typesDB "github.com/alfiejones/panda-ci/types/database"
	client "github.com/ory/client-go"
	"github.com/rs/zerolog/log"
)

var _oryClient *client.APIClient

func getOryAdminClient() (*client.APIClient, error) {
	url, err := env.GetOryAdminURL()
	if err != nil {
		return nil, err
	}

	if _oryClient != nil {
		return _oryClient, nil
	}

	cfg := client.NewConfiguration()
	cfg.Servers = client.ServerConfigurations{
		{URL: *url},
	}

	_oryClient = client.NewAPIClient(cfg)
	return _oryClient, nil
}

func getConfigs(ctx context.Context, identityID string) ([]interface{}, error) {
	authed := context.WithValue(ctx, client.ContextAccessToken, env.GetOryAdminToken())

	ory, err := getOryAdminClient()
	if err != nil {
		return []any{}, err
	}

	identity, _, err := ory.IdentityAPI.
		GetIdentity(authed, identityID).
		IncludeCredential([]string{"oidc"}).Execute()
	if err != nil {
		return []any{}, err
	}

	tokens := (*identity.Credentials)["oidc"]

	configsRaw, ok := tokens.GetConfig()["providers"]
	if !ok {
		return []any{}, nil
	}

	configs, ok := configsRaw.([]interface{})
	if !ok {
		return configs, fmt.Errorf("Failed to cast ocid providers to map")
	}

	return configs, nil
}

type Provider struct {
	Type                typesDB.UserAccountType
	InitialRefreshToken string
	ProviderAccountID   string
}

func GetGithubProvider(ctx context.Context, identityID string) (*Provider, error) {
	configs, err := getConfigs(ctx, identityID)
	if err != nil {
		return nil, err
	}

	for _, c := range configs {
		config, ok := c.(map[string]interface{})
		if !ok {
			log.Error().Msg("Failed to cast ocid provider to map")
			continue
		}

		switch config["provider"] {
		case "github":
			{
				return &Provider{
					Type:                typesDB.UserAccountTypeGithub,
					InitialRefreshToken: config["initial_refresh_token"].(string),
					ProviderAccountID:   config["subject"].(string),
				}, nil
			}
		}
	}

	return nil, nil
}

func GetUserByID(ctx context.Context, id string) (*types.User, error) {

	ory, err := getOryAdminClient()
	if err != nil {
		return nil, err
	}

	authed := context.WithValue(ctx, client.ContextAccessToken, env.GetOryAdminToken())

	identity, _, err := ory.IdentityAPI.GetIdentity(authed, id).Execute()
	if err != nil {
		return nil, err
	}

	traits := identity.GetTraits().(map[string]interface{})

	user := &types.User{
		ID: identity.GetId(),
	}

	if name, ok := traits["name"].(string); ok {
		user.Name = &name
	}

	if email, ok := traits["email"].(string); ok {
		user.Email = email
	}

	if avatar, ok := traits["avatar"].(string); ok {
		user.Avatar = &avatar
	}

	return user, nil
}

func ListUsersByIDs(ctx context.Context, ids []string) ([]*types.User, error) {
	ory, err := getOryAdminClient()
	if err != nil {
		return nil, err
	}

	authed := context.WithValue(ctx, client.ContextAccessToken, env.GetOryAdminToken())

	identities, _, err := ory.IdentityAPI.ListIdentities(authed).Ids(ids).Execute()
	if err != nil {
		return nil, err
	}

	users := make([]*types.User, len(identities))
	for i, identity := range identities {
		traits := identity.GetTraits().(map[string]interface{})

		user := &types.User{
			ID: identity.GetId(),
		}

		if name, ok := traits["name"].(string); ok {
			user.Name = &name
		}

		if email, ok := traits["email"].(string); ok {
			user.Email = email
		}

		if avatar, ok := traits["avatar"].(string); ok {
			user.Avatar = &avatar
		}

		users[i] = user
	}

	return users, nil
}

func GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	ory, err := getOryAdminClient()
	if err != nil {
		return nil, err
	}

	authed := context.WithValue(ctx, client.ContextAccessToken, env.GetOryAdminToken())

	identities, _, err := ory.IdentityAPI.ListIdentities(authed).CredentialsIdentifier(email).Execute()
	if err != nil {
		return nil, err
	}

	if len(identities) == 0 {
		return nil, nil
	} else if len(identities) > 1 {
		log.Error().Msgf("Multiple identities found for email %s", email)
		return nil, fmt.Errorf("Multiple identities found for email %s", email)
	}

	identity := identities[0]

	traits := identity.GetTraits().(map[string]interface{})

	user := &types.User{
		ID: identity.GetId(),
	}

	if name, ok := traits["name"].(string); ok {
		user.Name = &name
	}

	if email, ok := traits["email"].(string); ok {
		user.Email = email
	}

	return user, nil
}
