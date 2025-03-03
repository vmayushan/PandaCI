package orgQueries

import (
	"context"

	"github.com/pandaci-com/pandaci/types"
	typesDB "github.com/pandaci-com/pandaci/types/database"
)

func (q *OrgQueries) GetOrgByURLNameAsUser(ctx context.Context, user types.User, slug string) (*typesDB.OrgDB, error) {
	var org typesDB.OrgDB

	query := `SELECT
      org.id,
      org.slug,
      org.name,
      org.license,
      org.owner_user_id,
      org.avatar_url,
      org_users.role as current_users_role
    FROM  org
    JOIN  org_users ON org_users.org_id = org.id
    WHERE org.slug = $1 AND org_users.user_id = $2`

	if err := q.GetContext(ctx, &org, query, slug, user.ID); err != nil {
		return nil, err
	}

	return &org, nil
}

// Unsafe since we don't have any user context
// This is used for internal operations
func (q *OrgQueries) Unsafe_GetOrgByID(ctx context.Context, id string) (*typesDB.OrgDB, error) {
	var org typesDB.OrgDB

	query := `SELECT
	  org.id,
	  org.slug,
	  org.name,
	  org.license,
	  org.owner_user_id,
	  org.avatar_url
	FROM  org
	WHERE org.id = $1`

	if err := q.GetContext(ctx, &org, query, id); err != nil {
		return nil, err
	}

	return &org, nil
}

func (q *OrgQueries) GetOrgByIDAsUser(ctx context.Context, user types.User, id string) (*typesDB.OrgDB, error) {
	var org typesDB.OrgDB

	query := `SELECT
      org.id,
      org.slug,
      org.name,
      org.license,
      org.owner_user_id,
      org.avatar_url,
      org_users.role as current_users_role
    FROM  org
    JOIN  org_users ON org_users.org_id = org.id
    WHERE org.id = $1 AND org_users.user_id = $2`

	if err := q.GetContext(ctx, &org, query, id, user.ID); err != nil {
		return nil, err
	}

	return &org, nil
}

func (q *OrgQueries) GetUsersOrgs(ctx context.Context, userID string) (*[]typesDB.OrgDB, error) {
	var orgs []typesDB.OrgDB

	query := `SELECT
      org.id,
      org.slug,
      org.name,
      org.license,
      org.owner_user_id,
      org.avatar_url,
      org.owner_user_id,
      org_users.role as current_users_role
    FROM  org
    JOIN  org_users on org.id = org_users.org_id
    WHERE  org_users.user_id = $1`

	if err := q.SelectContext(ctx, &orgs, query, userID); err != nil {
		return nil, err
	}

	return &orgs, nil
}

func (q *OrgQueries) GetOrgUsers(ctx context.Context, orgID string) (*[]typesDB.OrgUsersDB, error) {
	var orgUsers []typesDB.OrgUsersDB

	query := `SELECT
      user_id,
      org_id,
      role
    FROM  org_users
    WHERE org_id = $1`

	if err := q.SelectContext(ctx, &orgUsers, query, orgID); err != nil {
		return nil, err
	}

	return &orgUsers, nil
}

func (q *OrgQueries) GetUsersFreeOrg(ctx context.Context, userID string) (*typesDB.OrgDB, error) {
	var org typesDB.OrgDB

	query := `SELECT
	  org.id,
	  org.slug,
	  org.name,
	  org.license,
	  org.owner_user_id,
	  org.avatar_url,
	  org_users.role as current_users_role
	FROM  org
	JOIN  org_users on org.id = org_users.org_id
	WHERE  org_users.user_id = $1 AND org.license->>'plan' = 'free'`

	if err := q.GetContext(ctx, &org, query, userID); err != nil {
		return nil, err
	}

	return &org, nil
}
