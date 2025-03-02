package queries

import (
	gitIntegrationQueries "github.com/alfiejones/panda-ci/app/queries/gitIntegration"
	orgQueries "github.com/alfiejones/panda-ci/app/queries/org"
	projectQueries "github.com/alfiejones/panda-ci/app/queries/project"
	queriesProjectEnvironment "github.com/alfiejones/panda-ci/app/queries/projectEnvironment"
	projectVariableQueries "github.com/alfiejones/panda-ci/app/queries/projectVariable"
	userAccountQueries "github.com/alfiejones/panda-ci/app/queries/userAccount"
	queriesWorkflow "github.com/alfiejones/panda-ci/app/queries/workflow"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	*orgQueries.OrgQueries
	*projectQueries.ProjectQueries
	*userAccountQueries.UserAccountQueries
	*gitIntegrationQueries.GitIntegrationQueries
	*queriesWorkflow.WorkflowQueries
	*projectVariableQueries.ProjectVariableQueries
	*queriesProjectEnvironment.ProjectEnvironmentQueries
}

func NewQueries(db *sqlx.DB) *Queries {
	return &Queries{
		OrgQueries: &orgQueries.OrgQueries{
			DB: db,
		},
		ProjectQueries: &projectQueries.ProjectQueries{
			DB: db,
		},
		UserAccountQueries: &userAccountQueries.UserAccountQueries{
			DB: db,
		},
		GitIntegrationQueries: &gitIntegrationQueries.GitIntegrationQueries{
			DB: db,
		},
		WorkflowQueries: &queriesWorkflow.WorkflowQueries{
			DB: db,
		},
		ProjectVariableQueries: &projectVariableQueries.ProjectVariableQueries{
			DB: db,
		},
		ProjectEnvironmentQueries: &queriesProjectEnvironment.ProjectEnvironmentQueries{
			DB: db,
		},
	}
}
