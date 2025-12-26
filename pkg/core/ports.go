package core

import "context"

// Repository defines the interface for storing and retrieving clients
// Interface is defined on consumer side (core) following hexagonal architecture
type Repository interface {
	// Load initializes the repository with the given password
	Load(ctx context.Context, password string) error
	// Save stores a client
	Save(ctx context.Context, client Client) error

	// Get retrieves a client by name
	Get(ctx context.Context, name string) (*Client, error)

	// List returns all client names
	List(ctx context.Context) ([]string, error)

	// GetAll returns all clients with full details
	GetAll(ctx context.Context) ([]Client, error)

	// Delete removes a client by name
	Delete(ctx context.Context, name string) error

	// Exists checks if repository is initialized
	Exists() bool
}

// Provider defines the interface for OAuth2 token operations
// Interface is defined on consumer side (core) following hexagonal architecture
type Provider interface {
	// GetToken obtains an access token using client credentials
	GetToken(ctx context.Context, client Client) (*Token, error)
}
