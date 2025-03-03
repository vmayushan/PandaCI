package middleware

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	ory "github.com/ory/client-go"
	"github.com/pandaci-com/pandaci/pkg/utils/env"
	"github.com/pandaci-com/pandaci/types"
)

type oryMiddleware struct {
	ory *ory.APIClient
}

func NewOryMiddleware() *oryMiddleware {
	url, err := env.GetOryURL()
	if err != nil {
		panic(err)
	}

	cfg := ory.NewConfiguration()
	cfg.Servers = ory.ServerConfigurations{
		{
			URL: *url,
		},
	}

	return &oryMiddleware{
		ory: ory.NewAPIClient(cfg),
	}
}

func GetUser(c echo.Context) types.User {
	user := c.Get("user").(*types.User)

	if user == nil {
		panic("this function isn't called in the context of the auth middleware")
	}

	return *user
}

func oryToUserModel(identity ory.Identity) types.User {
	user := types.User{
		ID: identity.GetId(),
	}

	traits := identity.GetTraits().(map[string]interface{})

	if name, ok := traits["name"].(string); ok {
		user.Name = &name
	}

	if email, ok := traits["email"].(string); ok {
		user.Email = email
	}

	return user
}

func (o *oryMiddleware) Session(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := o.validateSession(c.Request())
		if err != nil || !session.GetActive() {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		user := oryToUserModel(session.GetIdentity())

		c.Set("user", &user)

		return next(c)
	}
}

func (o *oryMiddleware) validateSession(r *http.Request) (*ory.Session, error) {
	authorization := r.Header.Get("Authorization")

	var req ory.FrontendAPIToSessionRequest

	if authorization != "" {

		if len(authorization) < 7 {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
		tokenType := authorization[:6]

		if tokenType != "Bearer" && tokenType != "bearer" {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		authorization = authorization[7:]

		req = o.ory.FrontendAPI.ToSession(r.Context()).XSessionToken(authorization)
	} else {
		decodedCookies, err := url.QueryUnescape(r.Header.Get("Cookie"))
		if err != nil {
			return nil, err
		}

		req = o.ory.FrontendAPI.ToSession(r.Context()).Cookie(decodedCookies)

	}

	resp, _, err := req.Execute()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
