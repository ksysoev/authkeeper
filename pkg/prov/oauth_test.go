package prov

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ksysoev/authkeeper/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOAuthProvider(t *testing.T) {
	provider := NewOAuthProvider()

	assert.NotNil(t, provider)
	assert.NotNil(t, provider.httpClient)
	assert.Equal(t, 30*time.Second, provider.httpClient.Timeout)
}

func TestOAuthProvider_GetToken(t *testing.T) {
	tests := []struct {
		name           string
		client         core.Client
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedToken  *core.Token
		expectedErr    string
	}{
		{
			name: "successful token request",
			client: core.Client{
				Name:         "test-client",
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				Scopes:       []string{"read", "write"},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
				assert.Equal(t, "application/json", r.Header.Get("Accept"))

				err := r.ParseForm()
				require.NoError(t, err)
				assert.Equal(t, "client_credentials", r.FormValue("grant_type"))
				assert.Equal(t, "test-client-id", r.FormValue("client_id"))
				assert.Equal(t, "test-client-secret", r.FormValue("client_secret"))
				assert.Equal(t, "read write", r.FormValue("scope"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"access_token": "test-access-token",
					"token_type":   "Bearer",
					"expires_in":   3600,
					"scope":        "read write",
				})
			},
			expectedToken: &core.Token{
				AccessToken: "test-access-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
				Scope:       "read write",
			},
			expectedErr: "",
		},
		{
			name: "successful token request without scopes",
			client: core.Client{
				Name:         "test-client",
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				Scopes:       []string{},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				err := r.ParseForm()
				require.NoError(t, err)
				assert.Equal(t, "", r.FormValue("scope"))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"access_token": "test-access-token",
					"token_type":   "Bearer",
					"expires_in":   3600,
				})
			},
			expectedToken: &core.Token{
				AccessToken: "test-access-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
				Scope:       "",
			},
			expectedErr: "",
		},
		{
			name: "server returns 400 bad request",
			client: core.Client{
				Name:         "test-client",
				ClientID:     "invalid-client-id",
				ClientSecret: "invalid-secret",
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error":             "invalid_client",
					"error_description": "Invalid client credentials",
				})
			},
			expectedToken: nil,
			expectedErr:   "token request failed with status 400",
		},
		{
			name: "server returns 401 unauthorized",
			client: core.Client{
				Name:         "test-client",
				ClientID:     "test-client-id",
				ClientSecret: "wrong-secret",
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error": "unauthorized",
				})
			},
			expectedToken: nil,
			expectedErr:   "token request failed with status 401",
		},
		{
			name: "invalid json response",
			client: core.Client{
				Name:         "test-client",
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("invalid json"))
			},
			expectedToken: nil,
			expectedErr:   "failed to parse token response",
		},
		{
			name: "server returns 500 internal error",
			client: core.Client{
				Name:         "test-client",
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal Server Error"))
			},
			expectedToken: nil,
			expectedErr:   "token request failed with status 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			client := tt.client
			client.TokenURL = server.URL

			provider := NewOAuthProvider()
			token, err := provider.GetToken(context.Background(), client)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}
		})
	}
}

func TestOAuthProvider_GetToken_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := core.Client{
		Name:         "test-client",
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     server.URL,
	}

	provider := NewOAuthProvider()
	token, err := provider.GetToken(ctx, client)

	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestOAuthProvider_GetToken_InvalidURL(t *testing.T) {
	client := core.Client{
		Name:         "test-client",
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     "://invalid-url",
	}

	provider := NewOAuthProvider()
	token, err := provider.GetToken(context.Background(), client)

	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Error(), "failed to create request")
}

func TestOAuthProvider_GetToken_NetworkError(t *testing.T) {
	client := core.Client{
		Name:         "test-client",
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     "http://localhost:9999/nonexistent",
	}

	provider := NewOAuthProvider()
	token, err := provider.GetToken(context.Background(), client)

	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Contains(t, err.Error(), "failed to send request")
}
