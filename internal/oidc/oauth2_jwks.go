package oidc

import (
	"crypto/rsa"
	"encoding/json"

	"gopkg.in/square/go-jose.v2"

	"github.com/authelia/authelia/internal/middlewares"
)

// JWKsGet handler serving the jwks used to verify the JWT tokens.
func JWKsGet(publicKey *rsa.PublicKey) middlewares.RequestHandler {
	return func(req *middlewares.AutheliaCtx) {
		key := jose.JSONWebKey{}
		key.Key = publicKey
		key.KeyID = "main-key"
		key.Algorithm = "RS256"
		key.Use = "sig"

		keySet := new(jose.JSONWebKeySet)
		keySet.Keys = append(keySet.Keys, key)

		if err := json.NewEncoder(req).Encode(keySet); err != nil {
			req.Error(err, "failed to serve jwk set")
		}
	}
}