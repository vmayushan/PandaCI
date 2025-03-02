package gitIntegrationQueries

import "github.com/jmoiron/sqlx"

type GitIntegrationQueries struct {
	*sqlx.DB
}
