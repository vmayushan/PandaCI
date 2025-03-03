package userAccountQueries

import (
	"context"

	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
)

func (q *UserAccountQueries) GetUserAccountByType(ctx context.Context, user types.User, accountType typesDB.UserAccountType) (*typesDB.UserAccount, error) {
	var account typesDB.UserAccount

	query := `SELECT
      user_id,
      provider_account_id,
      type,
      refresh_token,
      refresh_token_expires_at,
      access_token,
      access_token_expires_at
    FROM user_account
    WHERE user_id = $1 AND type = $2`

	return &account, q.GetContext(ctx, &account, query, user.ID, accountType)
}

func (q *UserAccountQueries) GetUserAccountByProviderAccountID(ctx context.Context, providerAccountID string, accountType typesDB.UserAccountType) (*typesDB.UserAccount, error) {
	var account typesDB.UserAccount

	query := `SELECT
	  user_id,
	  provider_account_id,
	  type,
	  refresh_token,
	  refresh_token_expires_at,
	  access_token,
	  access_token_expires_at
	FROM user_account
	WHERE provider_account_id = $1 AND type = $2`

	return &account, q.GetContext(ctx, &account, query, providerAccountID, accountType)
}

func (q *UserAccountQueries) GetAccountRefreshStateByID(ctx context.Context, id string) (typesDB.OAuthUserAccountRefresh, error) {
	var refreshState typesDB.OAuthUserAccountRefresh

	query := `SELECT
      id,
      user_id,
      type,
      created_at
    FROM user_account_oauth_state
    WHERE id = $1`

	err := q.GetContext(ctx, &refreshState, query, id)

	return refreshState, err
}
