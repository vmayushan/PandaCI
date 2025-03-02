package middleware

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	gitShared "github.com/alfiejones/panda-ci/app/git/shared"
)

func TranslateErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(404, "not found")
		}

		var gitOAuthError *gitShared.GitOAuthError
		if errors.As(err, &gitOAuthError) {

			log.Info().Msgf("GitOAuthError: %v", gitOAuthError)

			type message struct {
				Message     string `json:"message"`
				RedirectURL string `json:"redirectURL,omitempty"`
			}

			return c.JSON(401, message{
				Message:     gitOAuthError.Message,
				RedirectURL: gitOAuthError.RedirectURL,
			})
		}

		return err
	}
}
