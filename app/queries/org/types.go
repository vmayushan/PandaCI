package orgQueries

import "github.com/jmoiron/sqlx"

type OrgQueries struct {
	*sqlx.DB
}
