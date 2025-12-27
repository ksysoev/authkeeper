package ui

import (
	"context"

	"github.com/ksysoev/authkeeper/pkg/core"
)

// CoreService defines what UI needs from core (interface on consumer side)
type CoreService interface {
	AddClient(ctx context.Context, client core.Client) error
	GetClient(ctx context.Context, name string) (*core.Client, error)
	ListClients(ctx context.Context) ([]string, error)
	GetAllClients(ctx context.Context) ([]core.Client, error)
	DeleteClient(ctx context.Context, name string) error
	IssueToken(ctx context.Context, clientName string) (*core.Token, error)
	IsRepositoryInitialized() bool
	CheckPassword(ctx context.Context, password string) error
}

// CLI implements the command-line interface
type CLI struct {
	service CoreService
}

// NewCLI creates a new CLI
func NewCLI(service CoreService) *CLI {
	return &CLI{
		service: service,
	}
}
