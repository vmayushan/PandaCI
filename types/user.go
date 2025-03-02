package types

type User struct {
	ID     string  `json:"id"`
	Email  string  `json:"email"`
	Name   *string `json:"name"`
	Avatar *string `json:"avatar"`
}

type OrgRole string

const (
	ORG_USERS_ROLE_ADMIN  OrgRole = "admin"
	ORG_USERS_ROLE_MEMBER OrgRole = "member"
)
