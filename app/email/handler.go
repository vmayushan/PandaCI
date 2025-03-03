package email

import (
	"context"
	"net/smtp"
	"strings"

	nanoid "github.com/matoous/go-nanoid/v2"

	"github.com/pandaci-com/pandaci/pkg/utils/env"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

type EmailData struct {
	To      []string
	Subject string
	Body    string
}

func (h *Handler) SendEmail(ctx context.Context, data EmailData) error {

	randID, err := nanoid.New(8)
	if err != nil {
		return err
	}

	headers := []string{
		"From: PandaCI <support@pandaci.com>",
		"To: " + strings.Join(data.To, ","),
		"Subject: " + data.Subject,
		"Message-ID: " + randID,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
	}

	body := strings.Join(headers, "\r\n") + "\r\n\r\n" + data.Body

	host, err := env.GetSMTPHost()
	if err != nil {
		return err
	}

	port, err := env.GetSMTPPort()
	if err != nil {
		return err
	}

	username, err := env.GetSMTPUsername()
	if err != nil {
		return err
	}

	password, err := env.GetSMTPPassword()
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", *username, *password, *host)

	if err := smtp.SendMail(*host+":"+*port, auth, "support@pandaci.com", data.To, []byte(body)); err != nil {
		return err
	}

	return nil
}
