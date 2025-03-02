package typesDB

import "time"

type UserAccountType string

type UserAccount struct {
	ProviderAccountID string          `db:"provider_account_id"`
	Type              UserAccountType `db:"type"`
	UserID            string          `db:"user_id"`

	RefreshToken          string    `db:"refresh_token"`
	AccessToken           string    `db:"access_token"`
	AccessTokenExpiresAt  time.Time `db:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `db:"refresh_token_expires_at"`
}

const (
	UserAccountTypeGithub = "github"
)

type OAuthUserAccountRefresh struct {
	ID string `db:"id"`

	Type   UserAccountType `db:"type"`
	UserID string          `db:"user_id"`

	CreatedAt time.Time `db:"created_at"`
}
