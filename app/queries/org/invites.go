package orgQueries

import (
	"context"

	typesDB "github.com/pandaci-com/pandaci/types/database"
)

func (h *OrgQueries) CreateOrgInvite(ctx context.Context, invite typesDB.PendingOrgInviteDB) error {

	query := `INSERT INTO pending_org_invites (org_id, email, role) VALUES ($1, $2, $3)`

	_, err := h.ExecContext(ctx, query, invite.OrgID, invite.Email, invite.Role)
	if err != nil {
		return err
	}

	return nil
}

func (h *OrgQueries) GetOrgInvitesByEmail(ctx context.Context, email string) ([]typesDB.PendingOrgInviteDB, error) {

	var invites []typesDB.PendingOrgInviteDB

	query := `SELECT
		org_id,
		email,
		role
	FROM pending_org_invites
	WHERE email = $1 AND created_at > NOW() - INTERVAL '3 days'`

	if err := h.SelectContext(ctx, &invites, query, email); err != nil {
		return nil, err
	}

	return invites, nil
}

func (h *OrgQueries) DeleteOrgInvite(ctx context.Context, orgID, email string) error {

	query := `DELETE FROM pending_org_invites WHERE org_id = $1 AND email = $2`

	_, err := h.ExecContext(ctx, query, orgID, email)
	if err != nil {
		return err
	}

	return nil
}

func (h *OrgQueries) GetOrgInvitesByOrgID(ctx context.Context, orgID string) ([]typesDB.PendingOrgInviteDB, error) {

	var invites []typesDB.PendingOrgInviteDB

	query := `SELECT
		org_id,
		email,
		role
	FROM pending_org_invites
	WHERE org_id = $1 AND created_at > NOW() - INTERVAL '3 days'`

	if err := h.SelectContext(ctx, &invites, query, orgID); err != nil {
		return nil, err
	}

	return invites, nil
}

func (h *OrgQueries) CleanOldOrgInvites(ctx context.Context) error {

	query := `DELETE FROM pending_org_invites WHERE created_at < NOW() - INTERVAL '3 days'`

	_, err := h.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
