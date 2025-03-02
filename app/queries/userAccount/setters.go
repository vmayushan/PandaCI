package userAccountQueries

import (
	"context"
	"fmt"
	"time"

	typesDB "github.com/alfiejones/panda-ci/types/database"
	nanoid "github.com/matoous/go-nanoid/v2"
)

func (q *UserAccountQueries) UpsertAccount(ctx context.Context, account *typesDB.UserAccount) error {
	accountQuery := `INSERT INTO user_account
    (
      user_id,
      type,
      provider_account_id,
      access_token,
      access_token_expires_at,
      refresh_token,
      refresh_token_expires_at
    )
    VALUES (
        :user_id,
        :type,
        :provider_account_id,
        :access_token,
        :access_token_expires_at,
        :refresh_token,
        :refresh_token_expires_at
    )
    ON CONFLICT %s
    DO UPDATE
    SET
        provider_account_id = EXCLUDED.provider_account_id,
        access_token = EXCLUDED.access_token,
        access_token_expires_at = EXCLUDED.access_token_expires_at,
        refresh_token = EXCLUDED.refresh_token,
        refresh_token_expires_at = EXCLUDED.refresh_token_expires_at,
        user_id = EXCLUDED.user_id;`

	userConflictQuery := fmt.Sprintf(accountQuery, "(user_id, type)")
	providerConflictQuery := fmt.Sprintf(accountQuery, "(provider_account_id, type)")

	_, err := q.NamedExecContext(ctx, userConflictQuery, account)
	if err != nil {
		_, err = q.NamedExecContext(ctx, providerConflictQuery, account)
	}

	return err
}

func (q *UserAccountQueries) CreateAccountRefreshState(ctx context.Context, account *typesDB.UserAccount) (*typesDB.OAuthUserAccountRefresh, error) {
	id, err := nanoid.New()
	if err != nil {
		return nil, err
	}

	refreshState := &typesDB.OAuthUserAccountRefresh{
		UserID:    account.UserID,
		Type:      account.Type,
		CreatedAt: time.Now(),
		ID:        id,
	}

	refreshState.ID = id

	refreshQuery := `INSERT INTO user_account_oauth_state
      (id, user_id, type, created_at)
    VALUES
      (:id, :user_id, :type, :created_at)`

	_, err = q.NamedExecContext(ctx, refreshQuery, refreshState)

	return refreshState, err
}
