package typesDB

import (
	"encoding/json"

	"github.com/alfiejones/panda-ci/types"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	"github.com/rs/zerolog/log"
)

type OrgDB struct {
	ID               string              `db:"id"`
	Slug             string              `db:"slug"`
	Name             string              `db:"name"`
	License          *sqlxTypes.JSONText `db:"license"`
	OwnerID          string              `db:"owner_user_id"`
	AvatarURL        string              `db:"avatar_url"`
	CurrentUsersRole *types.OrgRole      `db:"current_users_role"`
}

type OrgUsersDB struct {
	OrgID  string        `db:"org_id"`
	UserID string        `db:"user_id"`
	Role   types.OrgRole `db:"role"`
}

func (o *OrgDB) SetLicense(license types.CloudLicense) error {
	rawBytes, err := json.Marshal(license)
	if err != nil {
		log.Err(err).Msg("Failed to marshal license")
		return err
	}

	if o.License == nil {
		o.License = &sqlxTypes.JSONText{}
	}

	return o.License.Scan(rawBytes)
}

func (o *OrgDB) GetLicense() (types.CloudLicense, error) {

	if o.License == nil {
		return types.CloudLicense{
			Plan: types.CloudSubscriptionPlanPaused,
		}, nil
	}

	var license types.CloudLicense
	if err := o.License.Unmarshal(&license); err != nil {
		return types.CloudLicense{}, err
	}

	return license, nil
}

type PendingOrgInviteDB struct {
	OrgID     string        `db:"org_id"`
	Email     string        `db:"email"`
	Role      types.OrgRole `db:"role"`
	CreatedAt string        `db:"created_at"`
}
