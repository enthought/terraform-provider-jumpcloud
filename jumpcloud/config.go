package jumpcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

const headerAccept = "application/json"

// tokenEndpoint is JumpCloud's OAuth2 token endpoint for service-account
// client_credentials grants. Service accounts are OAuth-only: the minted
// access token must be sent as `Authorization: Bearer`, NOT `x-api-key`
// (JumpCloud returns 401 for an OAuth token presented in x-api-key).
// It is a var rather than a const so tests can point it at a stub server.
var tokenEndpoint = "https://admin-oauth.id.jumpcloud.com/oauth2/token"

// Config holds the JC configuration
type Config struct {
	APIKey       string // Legacy x-api-key auth token (admin-user key)
	ClientID     string // Service-account OAuth client id
	ClientSecret string // Service-account OAuth client secret
	OrgID        string // Organization ID
}

// Client instantiates a jcapiv2.Configuration struct that is passed
// to every Resource operation.
//
// Two auth modes are supported:
//   - Service-account OAuth (client_id + client_secret): a client_credentials
//     grant mints a short-lived (~1h) Bearer access token sent as
//     `Authorization: Bearer <token>`. Preferred — not tied to a person.
//   - Legacy x-api-key (api_key): an admin-user API key sent as `x-api-key`.
//
// OAuth takes precedence when both are configured.
func (c *Config) Client() (interface{}, error) {
	config := jcapiv2.NewConfiguration()

	if c.ClientID != "" && c.ClientSecret != "" {
		token, err := fetchAccessToken(c.ClientID, c.ClientSecret)
		if err != nil {
			return nil, err
		}
		config.AddDefaultHeader("Authorization", "Bearer "+token)
	} else {
		config.AddDefaultHeader("x-api-key", c.APIKey)
	}

	if c.OrgID != "" {
		config.AddDefaultHeader("x-org-id", c.OrgID)
	}
	// Instantiate the API client
	return config, nil
}

// fetchAccessToken performs a client_credentials OAuth2 grant against
// JumpCloud's service-account token endpoint and returns the access token.
//
// The token is minted once per provider configuration (~1h TTL) — enough for a
// single plan/apply. An apply that runs longer than the token lifetime is an
// accepted edge case (the provider would need to be re-configured / re-run).
func fetchAccessToken(clientID, clientSecret string) (string, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("scope", "api")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("building JumpCloud token request: %w", err)
	}
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", headerAccept)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("requesting JumpCloud access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return "", fmt.Errorf("JumpCloud token endpoint returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decoding JumpCloud token response: %w", err)
	}
	if payload.AccessToken == "" {
		return "", fmt.Errorf("JumpCloud token endpoint returned an empty access_token")
	}
	return payload.AccessToken, nil
}
