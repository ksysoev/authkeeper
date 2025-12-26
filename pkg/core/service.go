package core

import (
	"context"
	"fmt"
)

// Service implements the core business logic for managing OIDC clients
type Service struct {
	repo Repository
	prov Provider
}

// NewService creates a new core service
func NewService(repo Repository, prov Provider) *Service {
	return &Service{
		repo: repo,
		prov: prov,
	}
}

// AddClient adds a new OIDC client to the repository
func (s *Service) AddClient(ctx context.Context, client Client) error {
	if client.Name == "" {
		return fmt.Errorf("client name is required")
	}
	if client.ClientID == "" {
		return fmt.Errorf("client ID is required")
	}
	if client.ClientSecret == "" {
		return fmt.Errorf("client secret is required")
	}
	if client.TokenURL == "" {
		return fmt.Errorf("token URL is required")
	}

	return s.repo.Save(ctx, client)
}

// GetClient retrieves a client by name
func (s *Service) GetClient(ctx context.Context, name string) (*Client, error) {
	return s.repo.Get(ctx, name)
}

// ListClients returns all client names
func (s *Service) ListClients(ctx context.Context) ([]string, error) {
	return s.repo.List(ctx)
}

// GetAllClients returns all clients with full details
func (s *Service) GetAllClients(ctx context.Context) ([]Client, error) {
	return s.repo.GetAll(ctx)
}

// DeleteClient removes a client
func (s *Service) DeleteClient(ctx context.Context, name string) error {
	return s.repo.Delete(ctx, name)
}

// IssueToken obtains an access token for the specified client
func (s *Service) IssueToken(ctx context.Context, clientName string) (*Token, error) {
	client, err := s.repo.Get(ctx, clientName)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	token, err := s.prov.GetToken(ctx, *client)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return token, nil
}

// IsRepositoryInitialized checks if the repository is initialized
func (s *Service) IsRepositoryInitialized() bool {
	return s.repo.Exists()
}

// CheckPassword verifies the provided password by attempting to load the repository
func (s *Service) CheckPassword(ctx context.Context, password string) error {
	return s.repo.Load(ctx, password)
}
