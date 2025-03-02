package gitIntegrationQueries

import (
	"context"

	typesDB "github.com/alfiejones/panda-ci/types/database"
	nanoid "github.com/matoous/go-nanoid/v2"
)

func (q *GitIntegrationQueries) GetInsertOrUpdateGitIntegration(ctx context.Context, integration *typesDB.GitIntegration) error {
	// If a user uninstalls our app and then reinstalls, the provider_id can change.
	// To help mitigate this, everytime we fetch the installtion, we update it if it changes

	var err error
	if integration.ID, err = nanoid.New(); err != nil {
		return err
	}

	query := `INSERT INTO git_integration
    (
      id,
      type,
      provider_account_id,
      provider_id
    )
    VALUES (
        :id,
        :type,
        :provider_account_id,
        :provider_id
    )
    ON CONFLICT (provider_id, type)
    DO UPDATE
    SET
        provider_id = EXCLUDED.provider_id
    RETURNING *`

	rows, err := q.NamedQueryContext(ctx, query, integration)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(integration); err != nil {
			return err
		}
	}

	return nil
}
