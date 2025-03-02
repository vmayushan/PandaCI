package userAccountQueries

import "github.com/jmoiron/sqlx"

type UserAccountQueries struct {
	*sqlx.DB
}
