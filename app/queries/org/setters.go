package orgQueries

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	queries_utils "github.com/alfiejones/panda-ci/app/queries/utils"
	"github.com/alfiejones/panda-ci/pkg/utils"
	"github.com/alfiejones/panda-ci/types"
	typesDB "github.com/alfiejones/panda-ci/types/database"
	nanoid "github.com/matoous/go-nanoid/v2"
)

// Insert a given org and admin user objects into the database
// ID is generated in this function overriding anything that was there before
func (q *OrgQueries) CreateOrg(ctx context.Context, org *typesDB.OrgDB) error {
	if !utils.IsURLNameValid(org.Slug) {
		return fmt.Errorf("Invalid team url")
	}

	var err error
	if org.ID, err = nanoid.New(); err != nil {
		return err
	}

	usersFreeOrg, err := q.GetUsersFreeOrg(ctx, org.OwnerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if usersFreeOrg != nil {
		// User already has a free org, so this one requires a paid plan
		org.SetLicense(types.CloudLicense{
			Plan: types.CloudSubscriptionPlanPaused,
		})
	} else {
		org.SetLicense(types.CloudLicense{
			Plan: types.CloudSubscriptionPlanFree,
			Features: types.Features{
				BuildMinutes:        6000,
				MaxBuildMinutes:     6000,
				Committers:          5,
				MaxCommitters:       5,
				MaxCloudRunnerScale: 4,
				MaxProjects:         10,
			},
		})
	}

	orgUser := typesDB.OrgUsersDB{
		OrgID:  org.ID,
		UserID: org.OwnerID,
		Role:   types.ORG_USERS_ROLE_ADMIN,
	}

	tx, err := q.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	completed := false
	defer queries_utils.Rollback(&completed, tx)

	orgQuery := `INSERT INTO org
      (id, slug, name, license, owner_user_id, avatar_url)
    VALUES
      (:id, :slug, :name, :license, :owner_user_id, :avatar_url)`

	userQuery := `INSERT INTO org_users
      (user_id, org_id, role)
    VALUES
      (:user_id, :org_id, :role)`

	if _, err := tx.NamedExecContext(ctx, orgQuery, org); err != nil {
		return err
	}

	if _, err := tx.NamedExecContext(ctx, userQuery, orgUser); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	completed = true

	return nil
}

func (q *OrgQueries) AddUserToOrg(ctx context.Context, orgUser *typesDB.OrgUsersDB) error {

	query := `INSERT INTO org_users
	  (user_id, org_id, role)
	VALUES
	  (:user_id, :org_id, :role)`

	_, err := q.NamedExecContext(ctx, query, orgUser)

	return err
}

func (q *OrgQueries) RemoveUserFromOrg(ctx context.Context, orgID, userID string) error {
	query := `DELETE FROM org_users
	  WHERE org_id = $1 AND user_id = $2`

	_, err := q.ExecContext(ctx, query, orgID, userID)

	return err
}

func (q *OrgQueries) UpdateOrg(ctx context.Context, org *typesDB.OrgDB) error {
	query := `UPDATE org
	  SET name = :name, slug = :slug, avatar_url = :avatar_url
	  WHERE id = :id`

	_, err := q.NamedExecContext(ctx, query, org)

	return err
}

func (q *OrgQueries) DeleteOrg(ctx context.Context, orgID string) error {
	query := `DELETE FROM org
	  WHERE id = $1`

	_, err := q.ExecContext(ctx, query, orgID)

	return err
}

func (q *OrgQueries) UpdateOrgLicense(ctx context.Context, org *typesDB.OrgDB) error {
	query := `UPDATE org
	  SET license = :license
	  WHERE id = :id`

	_, err := q.NamedExecContext(ctx, query, org)

	return err
}
