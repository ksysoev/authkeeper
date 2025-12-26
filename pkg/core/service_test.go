package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewService(t *testing.T) {
	repo := NewMockRepository(t)
	prov := NewMockProvider(t)

	svc := NewService(repo, prov)

	assert.NotNil(t, svc)
	assert.Equal(t, repo, svc.repo)
	assert.Equal(t, prov, svc.prov)
}

func TestService_AddClient(t *testing.T) {
	tests := []struct {
		name        string
		client      Client
		setupMock   func(*MockRepository)
		expectedErr string
	}{
		{
			name: "successful add",
			client: Client{
				Name:         "test-client",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
				TokenURL:     "https://example.com/token",
				Scopes:       []string{"read", "write"},
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: "",
		},
		{
			name: "missing name",
			client: Client{
				ClientID:     "client-id",
				ClientSecret: "client-secret",
				TokenURL:     "https://example.com/token",
			},
			setupMock:   func(repo *MockRepository) {},
			expectedErr: "client name is required",
		},
		{
			name: "missing client ID",
			client: Client{
				Name:         "test-client",
				ClientSecret: "client-secret",
				TokenURL:     "https://example.com/token",
			},
			setupMock:   func(repo *MockRepository) {},
			expectedErr: "client ID is required",
		},
		{
			name: "missing client secret",
			client: Client{
				Name:     "test-client",
				ClientID: "client-id",
				TokenURL: "https://example.com/token",
			},
			setupMock:   func(repo *MockRepository) {},
			expectedErr: "client secret is required",
		},
		{
			name: "missing token URL",
			client: Client{
				Name:         "test-client",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
			},
			setupMock:   func(repo *MockRepository) {},
			expectedErr: "token URL is required",
		},
		{
			name: "repository error",
			client: Client{
				Name:         "test-client",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
				TokenURL:     "https://example.com/token",
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("repo error"))
			},
			expectedErr: "repo error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			prov := NewMockProvider(t)
			tt.setupMock(repo)

			svc := NewService(repo, prov)
			err := svc.AddClient(context.Background(), tt.client)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetClient(t *testing.T) {
	tests := []struct {
		name        string
		clientName  string
		setupMock   func(*MockRepository)
		expected    *Client
		expectedErr bool
	}{
		{
			name:       "successful get",
			clientName: "test-client",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, "test-client").Return(&Client{
					Name:     "test-client",
					ClientID: "client-id",
				}, nil)
			},
			expected: &Client{
				Name:     "test-client",
				ClientID: "client-id",
			},
			expectedErr: false,
		},
		{
			name:       "not found",
			clientName: "nonexistent",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, "nonexistent").Return(nil, errors.New("not found"))
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			prov := NewMockProvider(t)
			tt.setupMock(repo)

			svc := NewService(repo, prov)
			client, err := svc.GetClient(context.Background(), tt.clientName)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, client)
			}
		})
	}
}

func TestService_ListClients(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockRepository)
		expected    []string
		expectedErr bool
	}{
		{
			name: "successful list",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().List(mock.Anything).Return([]string{"client1", "client2"}, nil)
			},
			expected:    []string{"client1", "client2"},
			expectedErr: false,
		},
		{
			name: "empty list",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().List(mock.Anything).Return([]string{}, nil)
			},
			expected:    []string{},
			expectedErr: false,
		},
		{
			name: "repository error",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().List(mock.Anything).Return(nil, errors.New("repo error"))
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			prov := NewMockProvider(t)
			tt.setupMock(repo)

			svc := NewService(repo, prov)
			clients, err := svc.ListClients(context.Background())

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, clients)
			}
		})
	}
}

func TestService_GetAllClients(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		setupMock   func(*MockRepository)
		expected    []Client
		expectedErr bool
	}{
		{
			name: "successful get all",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().GetAll(mock.Anything).Return([]Client{
					{Name: "client1", ClientID: "id1", CreatedAt: now},
					{Name: "client2", ClientID: "id2", CreatedAt: now},
				}, nil)
			},
			expected: []Client{
				{Name: "client1", ClientID: "id1", CreatedAt: now},
				{Name: "client2", ClientID: "id2", CreatedAt: now},
			},
			expectedErr: false,
		},
		{
			name: "repository error",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().GetAll(mock.Anything).Return(nil, errors.New("repo error"))
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			prov := NewMockProvider(t)
			tt.setupMock(repo)

			svc := NewService(repo, prov)
			clients, err := svc.GetAllClients(context.Background())

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, clients)
			}
		})
	}
}

func TestService_DeleteClient(t *testing.T) {
	tests := []struct {
		name        string
		clientName  string
		setupMock   func(*MockRepository)
		expectedErr bool
	}{
		{
			name:       "successful delete",
			clientName: "test-client",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Delete(mock.Anything, "test-client").Return(nil)
			},
			expectedErr: false,
		},
		{
			name:       "repository error",
			clientName: "test-client",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Delete(mock.Anything, "test-client").Return(errors.New("repo error"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			prov := NewMockProvider(t)
			tt.setupMock(repo)

			svc := NewService(repo, prov)
			err := svc.DeleteClient(context.Background(), tt.clientName)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_IssueToken(t *testing.T) {
	tests := []struct {
		name        string
		clientName  string
		setupMock   func(*MockRepository, *MockProvider)
		expected    *Token
		expectedErr string
	}{
		{
			name:       "successful token issue",
			clientName: "test-client",
			setupMock: func(repo *MockRepository, prov *MockProvider) {
				client := &Client{
					Name:         "test-client",
					ClientID:     "client-id",
					ClientSecret: "client-secret",
					TokenURL:     "https://example.com/token",
				}
				repo.EXPECT().Get(mock.Anything, "test-client").Return(client, nil)
				prov.EXPECT().GetToken(mock.Anything, *client).Return(&Token{
					AccessToken: "access-token",
					TokenType:   "Bearer",
					ExpiresIn:   3600,
				}, nil)
			},
			expected: &Token{
				AccessToken: "access-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			},
			expectedErr: "",
		},
		{
			name:       "client not found",
			clientName: "nonexistent",
			setupMock: func(repo *MockRepository, prov *MockProvider) {
				repo.EXPECT().Get(mock.Anything, "nonexistent").Return(nil, errors.New("not found"))
			},
			expected:    nil,
			expectedErr: "failed to get client",
		},
		{
			name:       "provider error",
			clientName: "test-client",
			setupMock: func(repo *MockRepository, prov *MockProvider) {
				client := &Client{
					Name:         "test-client",
					ClientID:     "client-id",
					ClientSecret: "client-secret",
					TokenURL:     "https://example.com/token",
				}
				repo.EXPECT().Get(mock.Anything, "test-client").Return(client, nil)
				prov.EXPECT().GetToken(mock.Anything, *client).Return(nil, errors.New("token error"))
			},
			expected:    nil,
			expectedErr: "failed to get token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			prov := NewMockProvider(t)
			tt.setupMock(repo, prov)

			svc := NewService(repo, prov)
			token, err := svc.IssueToken(context.Background(), tt.clientName)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, token)
			}
		})
	}
}

func TestService_IsRepositoryInitialized(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*MockRepository)
		expected  bool
	}{
		{
			name: "repository exists",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Exists().Return(true)
			},
			expected: true,
		},
		{
			name: "repository does not exist",
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Exists().Return(false)
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			prov := NewMockProvider(t)
			tt.setupMock(repo)

			svc := NewService(repo, prov)
			result := svc.IsRepositoryInitialized()

			assert.Equal(t, tt.expected, result)
		})
	}
}
