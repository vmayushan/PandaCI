package email

import (
	"context"
	_ "embed"
	"strings"
	"text/template"
	"html"
)

type OrgInviteEmailNewUserData struct {
	To string

	TeamName       string
	TeamAvatarURL  string
	TeamSlug       string
	InvitedByName  string
	InvitedByEmail string
}

//go:embed templates/invite-new.tmpl.html
var orgInviteEmailNewUserTemplate string

// SendOrgInviteEmailToNewUser sends an email to a new user inviting them to join an organization
func (h *Handler) SendOrgInviteEmailToNewUser(ctx context.Context, data OrgInviteEmailNewUserData) error {
	data.To = html.EscapeString(data.To)
	data.TeamName = html.EscapeString(data.TeamName)
	data.TeamAvatarURL = html.EscapeString(data.TeamAvatarURL)
	data.TeamSlug = html.EscapeString(data.TeamSlug)
	data.InvitedByName = html.EscapeString(data.InvitedByName)
	data.InvitedByEmail = html.EscapeString(data.InvitedByEmail)

	tmpl, err := template.New("InviteEmail").Parse(orgInviteEmailNewUserTemplate)
	if err != nil {
		return err
	}

	var body strings.Builder

	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	return h.SendEmail(ctx, EmailData{
		To:      []string{data.To},
		Body:    body.String(),
		Subject: "You've been invited to join " + data.TeamName + " on PandaCI",
	})
}

type OrgInviteEmailExistingUserData struct {
	To string

	InviteeName string

	TeamName       string
	TeamAvatarURL  string
	TeamSlug       string
	InvitedByName  string
	InvitedByEmail string
}

//go:embed templates/invite-existing.tmpl.html
var orgInviteEmailExistingUserTemplate string

// Sends an email intended for an existing user who has been invited to an organization
func (h *Handler) SendOrgInviteEmailToExistingUser(ctx context.Context, data OrgInviteEmailExistingUserData) error {
	data.To = html.EscapeString(data.To)
	data.InviteeName = html.EscapeString(data.InviteeName)
	data.TeamName = html.EscapeString(data.TeamName)
	data.TeamAvatarURL = html.EscapeString(data.TeamAvatarURL)
	data.TeamSlug = html.EscapeString(data.TeamSlug)
	data.InvitedByName = html.EscapeString(data.InvitedByName)
	data.InvitedByEmail = html.EscapeString(data.InvitedByEmail)

	tmpl, err := template.New("InviteEmail").Parse(orgInviteEmailExistingUserTemplate)
	if err != nil {
		return err
	}

	var body strings.Builder

	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	return h.SendEmail(ctx, EmailData{
		To:      []string{data.To},
		Body:    body.String(),
		Subject: "You've been invited to join " + data.TeamName + " on PandaCI",
	})
}
