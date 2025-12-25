package core

import "time"

// Client represents an OIDC/OAuth2 client configuration
type Client struct {
	Name         string
	ClientID     string
	ClientSecret string
	TokenURL     string
	Scopes       []string
	CreatedAt    time.Time
}

// Token represents an OAuth2 access token response
type Token struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int
	Scope       string
}
