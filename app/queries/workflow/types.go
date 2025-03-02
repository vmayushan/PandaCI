package queriesWorkflow

import "github.com/jmoiron/sqlx"

type WorkflowQueries struct {
	*sqlx.DB
}
