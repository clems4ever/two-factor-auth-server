package validator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authelia/authelia/internal/configuration/schema"
)

func TestShouldRaiseErrorWhenInvalidOIDCServerConfiguration(t *testing.T) {
	validator := schema.NewStructValidator()
	config := &schema.OAuthConfiguration{
		OIDCServer: &schema.OpenIDConnectServerConfiguration{
			HMACSecret:           "abc",
			IssuerPrivateKeyPath: "",
		},
	}

	ValidateOAuth(config, validator)

	require.Len(t, validator.Errors(), 3)

	assert.EqualError(t, validator.Errors()[0], "OIDC Server issuer private key path must be provided")
	assert.EqualError(t, validator.Errors()[1], fmt.Sprintf(errOAuthOIDCServerHMACLengthMustBe32Fmt, 3))
	assert.EqualError(t, validator.Errors()[2], "OIDC Server has no clients defined")
}

func TestShouldRaiseErrorWhenOIDCServerIssuerPrivateKeyPathInvalid(t *testing.T) {
	validator := schema.NewStructValidator()
	config := &schema.OAuthConfiguration{
		OIDCServer: &schema.OpenIDConnectServerConfiguration{
			HMACSecret:           "rLABDrx87et5KvRHVUgTm3pezWWd8LMN",
			IssuerPrivateKeyPath: "../abc",
		},
	}

	ValidateOAuth(config, validator)

	require.Len(t, validator.Errors(), 2)

	assert.EqualError(t, validator.Errors()[0], "OIDC Server issuer private key path doesn't exist")
	assert.EqualError(t, validator.Errors()[1], "OIDC Server has no clients defined")
}

func TestShouldRaiseErrorWhenOIDCServerClientBadValues(t *testing.T) {
	validator := schema.NewStructValidator()
	config := &schema.OAuthConfiguration{
		OIDCServer: &schema.OpenIDConnectServerConfiguration{
			HMACSecret:           "rLABDrx87et5KvRHVUgTm3pezWWd8LMN",
			IssuerPrivateKeyPath: "../../../README.md",
			Clients: []schema.OpenIDConnectClientConfiguration{
				{
					ID:     "",
					Secret: "",
					Policy: "",
					RedirectURIs: []string{
						"http://google.com",
					},
				},
				{
					ID:     "a-client",
					Secret: "a-secret",
					Policy: "a-policy",
					RedirectURIs: []string{
						"https://google.com",
					},
				},
				{
					ID:     "a-client",
					Secret: "a-secret",
					Policy: "a-policy",
					RedirectURIs: []string{
						"https://google.com",
					},
				},
			},
		},
	}

	ValidateOAuth(config, validator)

	require.Len(t, validator.Errors(), 5)

	assert.EqualError(t, validator.Errors()[0], fmt.Sprintf(errOAuthOIDCServerClientRedirectURIFmt, "http://google.com", "http"))
	assert.EqualError(t, validator.Errors()[1], "OIDC Server has one or more clients with an empty ID")
	assert.EqualError(t, validator.Errors()[2], "OIDC Server has one or more clients with an empty secret")
	assert.EqualError(t, validator.Errors()[3], "OIDC Server has one or more clients with an empty policy")
	assert.EqualError(t, validator.Errors()[4], "OIDC Server has clients with duplicate ID's")
}

func TestShouldNotRaiseErrorWhenOIDCServerConfiguredCorrectly(t *testing.T) {
	validator := schema.NewStructValidator()
	config := &schema.OAuthConfiguration{
		OIDCServer: &schema.OpenIDConnectServerConfiguration{
			HMACSecret:           "rLABDrx87et5KvRHVUgTm3pezWWd8LMN",
			IssuerPrivateKeyPath: "../../../README.md",
			Clients: []schema.OpenIDConnectClientConfiguration{
				{
					ID:     "a-client",
					Secret: "a-client-secret",
					Policy: "one_factor",
					RedirectURIs: []string{
						"https://google.com",
					},
				},
				{
					ID:     "b-client",
					Secret: "b-client-secret",
					Policy: "one_factor",
					RedirectURIs: []string{
						"https://google.com",
					},
					Scopes: []string{
						"groups",
					},
				},
			},
		},
	}

	ValidateOAuth(config, validator)

	assert.Len(t, validator.Errors(), 0)

	require.Len(t, config.OIDCServer.Clients[0].Scopes, 1)
	assert.Equal(t, "openid", config.OIDCServer.Clients[0].Scopes[0])

	require.Len(t, config.OIDCServer.Clients[1].Scopes, 1)
	assert.Equal(t, "groups", config.OIDCServer.Clients[1].Scopes[0])
}