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

func TestCoreService_AddClient(t *testing.T) {
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

func TestCoreService_IsRepositoryInitialized(t *testing.T) {
	t.Run("initialized", func(t *testing.T) {
		service := NewMockCoreService(t)
		cli := NewCLI(service)

		service.EXPECT().IsRepositoryInitialized().Return(true)

		assert.True(t, cli.service.IsRepositoryInitialized())
	})

	t.Run("not initialized", func(t *testing.T) {
		service := NewMockCoreService(t)
		cli := NewCLI(service)

		service.EXPECT().IsRepositoryInitialized().Return(false)

		assert.False(t, cli.service.IsRepositoryInitialized())
	})
}

func TestCoreService_ListClients_NoClients(t *testing.T) {
	service := NewMockCoreService(t)
	service.EXPECT().ListClients(mock.Anything).Return([]string{}, nil)

	clients, err := service.ListClients(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, clients)
}

func TestCoreService_IssueToken_Success(t *testing.T) {
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

func TestCoreService_IssueToken_Error(t *testing.T) {
	service := NewMockCoreService(t)
	service.EXPECT().IssueToken(mock.Anything, "client1").Return(nil, errors.New("token error"))

	token, err := service.IssueToken(context.Background(), "client1")
	assert.Error(t, err)
	assert.Nil(t, token)
}

func TestCoreService_GetAllClients_Success(t *testing.T) {
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

func TestCoreService_GetAllClients_Empty(t *testing.T) {
	service := NewMockCoreService(t)
	service.EXPECT().GetAllClients(mock.Anything).Return([]core.Client{}, nil)

	result, err := service.GetAllClients(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestCoreService_DeleteClient_Success(t *testing.T) {
	service := NewMockCoreService(t)
	service.EXPECT().DeleteClient(mock.Anything, "client1").Return(nil)

	err := service.DeleteClient(context.Background(), "client1")
	assert.NoError(t, err)
}

func TestCoreService_DeleteClient_Error(t *testing.T) {
	service := NewMockCoreService(t)
	cli := NewCLI(service)

	service.EXPECT().DeleteClient(mock.Anything, "client1").Return(errors.New("delete error"))

	err := cli.service.DeleteClient(context.Background(), "client1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete error")
}

func TestCoreService_GetClient_Success(t *testing.T) {
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

func TestCoreService_GetClient_NotFound(t *testing.T) {
	service := NewMockCoreService(t)
	cli := NewCLI(service)

	service.EXPECT().GetClient(mock.Anything, "nonexistent").Return(nil, errors.New("not found"))

	client, err := cli.service.GetClient(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestCoreService_CheckPassword(t *testing.T) {
	t.Run("valid password", func(t *testing.T) {
		service := NewMockCoreService(t)
		cli := NewCLI(service)

		service.EXPECT().CheckPassword(mock.Anything, "validpassword").Return(nil)

		err := cli.service.CheckPassword(context.Background(), "validpassword")
		assert.NoError(t, err)
	})

	t.Run("invalid password", func(t *testing.T) {
		service := NewMockCoreService(t)
		cli := NewCLI(service)

		service.EXPECT().CheckPassword(mock.Anything, "wrong").Return(errors.New("invalid password"))

		err := cli.service.CheckPassword(context.Background(), "wrong")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid password")
	})
}
