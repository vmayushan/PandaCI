package jwt

import (
	"time"
)

type JWTOpts struct {
	ExpiresIn time.Duration
	PublicKey *string
}

type JWTHandler struct {
	expiresIn time.Duration
	publicKey *string
}

func NewJWTHelper(opts JWTOpts) JWTHandler {
	return JWTHandler{
		expiresIn: opts.ExpiresIn,
		publicKey: opts.PublicKey,
	}
}

const DefaultExpireTime = time.Hour * 24 * 7
