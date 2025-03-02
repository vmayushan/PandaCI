package queriesProjectEnvironment

import "github.com/jmoiron/sqlx"

type ProjectEnvironmentQueries struct {
	*sqlx.DB
}
