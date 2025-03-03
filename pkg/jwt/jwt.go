package jwt

import (
	"fmt"
	"time"

	"github.com/pandaci-com/pandaci/pkg/utils/env"
	"github.com/golang-jwt/jwt/v5"
)

func (h *JWTHandler) CreateWorkflowToken(claims WorkflowClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"workflowID": claims.WorkflowID,
			"orgID":      claims.OrgID,
			"projectID":  claims.ProjectID,
			"exp":        time.Now().Add(h.expiresIn).Unix(),
		})

	privateKeyStr, err := env.GetRunnerPrivateKey()
	if err != nil {
		return "", err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(*privateKeyStr))
	if err != nil {
		return "", fmt.Errorf("could not parse private key: %w", err)
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type WorkflowClaims struct {
	WorkflowID string `json:"workflowID"`
	OrgID      string `json:"orgID"`
	ProjectID  string `json:"projectID"`
}

func (h *JWTHandler) ValidateWorkflowToken(tokenString string) (WorkflowClaims, error) {
	publicKeyStr := h.publicKey

	if publicKeyStr == nil {
		var err error
		publicKeyStr, err = env.GetRunnerPublicKey()
		if err != nil {
			return WorkflowClaims{}, err
		}
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(*publicKeyStr))
	if err != nil {
		return WorkflowClaims{}, fmt.Errorf("could not parse public key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return WorkflowClaims{}, fmt.Errorf("could not parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return WorkflowClaims{}, jwt.ErrInvalidKey
	}

	workflowID, ok := claims["workflowID"].(string)
	if !ok {
		return WorkflowClaims{}, jwt.ErrInvalidKey
	}

	orgID, ok := claims["orgID"].(string)
	if !ok {
		return WorkflowClaims{}, jwt.ErrInvalidKey
	}

	projectID, ok := claims["projectID"].(string)
	if !ok {
		return WorkflowClaims{}, jwt.ErrInvalidKey
	}

	return WorkflowClaims{
		WorkflowID: workflowID,
		OrgID:      orgID,
		ProjectID:  projectID,
	}, nil
}
