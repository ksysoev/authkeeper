package prov

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ksysoev/authkeeper/pkg/core"
)

// OAuthProvider implements core.Provider interface for OAuth2 operations
type OAuthProvider struct {
	httpClient *http.Client
}

// NewOAuthProvider creates a new OAuth provider
func NewOAuthProvider() *OAuthProvider {
	return &OAuthProvider{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetToken obtains an access token using client credentials flow
func (p *OAuthProvider) GetToken(ctx context.Context, client core.Client) (*core.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", client.ClientID)
	data.Set("client_secret", client.ClientSecret)

	if len(client.Scopes) > 0 {
		data.Set("scope", strings.Join(client.Scopes, " "))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope,omitempty"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &core.Token{
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		ExpiresIn:   tokenResp.ExpiresIn,
		Scope:       tokenResp.Scope,
	}, nil
}
