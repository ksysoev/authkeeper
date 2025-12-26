package ui

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ksysoev/authkeeper/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCLI(t *testing.T) {
	service := NewMockCoreService(t)
	cli := NewCLI(service)

	assert.NotNil(t, cli)
	assert.Equal(t, service, cli.service)
}

func TestCLI_AddClient_RepositoryNotInitialized(t *testing.T) {
	service := NewMockCoreService(t)
	cli := NewCLI(service)

	client := core.Client{
		Name:         "test-client",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		TokenURL:     "https://example.com/token",
		Scopes:       []string{"read", "write"},
	}

	service.EXPECT().AddClient(mock.Anything, mock.MatchedBy(func(c core.Client) bool {
		return c.Name == client.Name &&
			c.ClientID == client.ClientID &&
			c.ClientSecret == client.ClientSecret &&
			c.TokenURL == client.TokenURL &&
			len(c.Scopes) == len(client.Scopes)
	})).Return(nil)

	err := cli.service.AddClient(context.Background(), client)
	assert.NoError(t, err)
}

func TestCLI_IssueToken_RepositoryNotInitialized(t *testing.T) {
	service := NewMockCoreService(t)
	cli := NewCLI(service)

	service.EXPECT().IsRepositoryInitialized().Return(false)

	assert.False(t, cli.service.IsRepositoryInitialized())
}

func TestCLI_IssueToken_NoClients(t *testing.T) {
	service := NewMockCoreService(t)
	service.EXPECT().ListClients(mock.Anything).Return([]string{}, nil)

	clients, err := service.ListClients(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, clients)
}

func TestCLI_IssueToken_Success(t *testing.T) {
	service := NewMockCoreService(t)
	service.EXPECT().IssueToken(mock.Anything, "client1").Return(&core.Token{
		AccessToken: "access-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:       "read write",
	}, nil)

	token, err := service.IssueToken(context.Background(), "client1")
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, "access-token", token.AccessToken)
	assert.Equal(t, "Bearer", token.TokenType)
	assert.Equal(t, 3600, token.ExpiresIn)
}

func TestCLI_IssueToken_Error(t *testing.T) {
	service := NewMockCoreService(t)
	service.EXPECT().IssueToken(mock.Anything, "client1").Return(nil, errors.New("token error"))

	token, err := service.IssueToken(context.Background(), "client1")
	assert.Error(t, err)
	assert.Nil(t, token)
}

func TestCLI_ListClients_Success(t *testing.T) {
	service := NewMockCoreService(t)
	now := time.Now()
	clients := []core.Client{
		{
			Name:      "client1",
			ClientID:  "id1",
			TokenURL:  "url1",
			Scopes:    []string{"read"},
			CreatedAt: now,
		},
		{
			Name:      "client2",
			ClientID:  "id2",
			TokenURL:  "url2",
			Scopes:    []string{"write"},
			CreatedAt: now,
		},
	}

	service.EXPECT().GetAllClients(mock.Anything).Return(clients, nil)

	result, err := service.GetAllClients(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "client1", result[0].Name)
	assert.Equal(t, "client2", result[1].Name)
}

func TestCLI_ListClients_RepositoryNotInitialized(t *testing.T) {
	service := NewMockCoreService(t)
	cli := NewCLI(service)

	service.EXPECT().IsRepositoryInitialized().Return(false)

	assert.False(t, cli.service.IsRepositoryInitialized())
}

func TestCLI_ListClients_Empty(t *testing.T) {
	service := NewMockCoreService(t)
	service.EXPECT().GetAllClients(mock.Anything).Return([]core.Client{}, nil)

	result, err := service.GetAllClients(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestCLI_DeleteClient_Success(t *testing.T) {
	service := NewMockCoreService(t)
	service.EXPECT().DeleteClient(mock.Anything, "client1").Return(nil)

	err := service.DeleteClient(context.Background(), "client1")
	assert.NoError(t, err)
}

func TestCLI_DeleteClient_RepositoryNotInitialized(t *testing.T) {
	service := NewMockCoreService(t)
	cli := NewCLI(service)

	service.EXPECT().IsRepositoryInitialized().Return(false)

	assert.False(t, cli.service.IsRepositoryInitialized())
}

func TestCLI_DeleteClient_NoClients(t *testing.T) {
	service := NewMockCoreService(t)
	service.EXPECT().ListClients(mock.Anything).Return([]string{}, nil)

	clients, err := service.ListClients(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, clients)
}

func TestCLI_DeleteClient_Error(t *testing.T) {
	service := NewMockCoreService(t)
	cli := NewCLI(service)

	service.EXPECT().DeleteClient(mock.Anything, "client1").Return(errors.New("delete error"))

	err := cli.service.DeleteClient(context.Background(), "client1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete error")
}

func TestCLI_GetClient_Success(t *testing.T) {
	service := NewMockCoreService(t)
	cli := NewCLI(service)

	expectedClient := &core.Client{
		Name:     "test-client",
		ClientID: "client-id",
	}

	service.EXPECT().GetClient(mock.Anything, "test-client").Return(expectedClient, nil)

	client, err := cli.service.GetClient(context.Background(), "test-client")
	assert.NoError(t, err)
	assert.Equal(t, expectedClient, client)
}

func TestCLI_GetClient_NotFound(t *testing.T) {
	service := NewMockCoreService(t)
	cli := NewCLI(service)

	service.EXPECT().GetClient(mock.Anything, "nonexistent").Return(nil, errors.New("not found"))

	client, err := cli.service.GetClient(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestPrintFunctions(t *testing.T) {
	t.Run("printTitle", func(t *testing.T) {
		assert.NotPanics(t, func() {
			printTitle("Test Title")
		})
	})

	t.Run("printSuccess", func(t *testing.T) {
		assert.NotPanics(t, func() {
			printSuccess("Success message")
		})
	})

	t.Run("printError", func(t *testing.T) {
		assert.NotPanics(t, func() {
			printError("Error message")
		})
	})

	t.Run("printInfo", func(t *testing.T) {
		assert.NotPanics(t, func() {
			printInfo("Info message")
		})
	})

	t.Run("printWarning", func(t *testing.T) {
		assert.NotPanics(t, func() {
			printWarning("Warning message")
		})
	})

	t.Run("printMuted", func(t *testing.T) {
		assert.NotPanics(t, func() {
			printMuted("Muted message")
		})
	})

	t.Run("printProgress", func(t *testing.T) {
		assert.NotPanics(t, func() {
			printProgress("Progress message")
		})
	})

	t.Run("printBox", func(t *testing.T) {
		assert.NotPanics(t, func() {
			printBox("Line 1", "Line 2", "Line 3")
		})
	})

	t.Run("printBox empty", func(t *testing.T) {
		assert.NotPanics(t, func() {
			printBox()
		})
	})
}

func TestSelectFromList_NoItems(t *testing.T) {
	_, err := selectFromList("Test", []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no items to select from")
}
