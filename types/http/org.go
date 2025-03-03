package typesHTTP

import "github.com/pandaci-com/pandaci/types"

type UpdateOrgHTTP struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	AvatarURL string `json:"avatarURL"`
}

type OrgHTTP struct {
	ID               string             `json:"id"`
	Slug             string             `json:"slug"`
	Name             string             `json:"name"`
	License          types.CloudLicense `json:"license,omitempty"`
	OwnerID          string             `json:"ownerID"`
	AvatarURL        string             `json:"avatarURL"`
	CurrentUsersRole *types.OrgRole     `json:"currentUsersRole,omitempty"`
}

type OrgUserHTTP struct {
	types.User
	Role types.OrgRole `json:"role"`
}
