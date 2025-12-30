package repo

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ksysoev/authkeeper/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVaultRepository(t *testing.T) {
	repo := NewVaultRepository("/tmp/vault.enc")

	assert.NotNil(t, repo)
	assert.Equal(t, "/tmp/vault.enc", repo.path)
}

func TestVaultRepository_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")

	repo := NewVaultRepository(vaultPath)

	assert.False(t, repo.Exists())

	err := os.WriteFile(vaultPath, []byte("test"), 0600)
	require.NoError(t, err)

	assert.True(t, repo.Exists())
}

func TestVaultRepository_SaveAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")
	password := "test-password"

	repo := NewVaultRepository(vaultPath)
	ctx := context.Background()

	err := repo.Load(ctx, password)
	require.NoError(t, err)

	client := core.Client{
		Name:         "test-client",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		TokenURL:     "https://example.com/token",
		Scopes:       []string{"read", "write"},
		CreatedAt:    time.Now().Truncate(time.Second),
	}

	err = repo.Save(ctx, client)
	require.NoError(t, err)

	retrieved, err := repo.Get(ctx, "test-client")
	require.NoError(t, err)
	assert.Equal(t, client.Name, retrieved.Name)
	assert.Equal(t, client.ClientID, retrieved.ClientID)
	assert.Equal(t, client.ClientSecret, retrieved.ClientSecret)
	assert.Equal(t, client.TokenURL, retrieved.TokenURL)
	assert.Equal(t, client.Scopes, retrieved.Scopes)
	assert.True(t, client.CreatedAt.Equal(retrieved.CreatedAt))
}

func TestVaultRepository_Save_DuplicateName(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")
	repo := NewVaultRepository(vaultPath)
	ctx := context.Background()

	err := repo.Load(ctx, "password")
	require.NoError(t, err)

	client1 := core.Client{
		Name:         "test-client",
		ClientID:     "client-id-1",
		ClientSecret: "secret-1",
		TokenURL:     "https://example.com/token",
	}

	err = repo.Save(ctx, client1)
	require.NoError(t, err)

	client2 := core.Client{
		Name:         "test-client",
		ClientID:     "client-id-2",
		ClientSecret: "secret-2",
		TokenURL:     "https://example.com/token",
	}

	err = repo.Save(ctx, client2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestVaultRepository_Save_AutoTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")
	repo := NewVaultRepository(vaultPath)
	ctx := context.Background()

	err := repo.Load(ctx, "password")
	require.NoError(t, err)

	client := core.Client{
		Name:         "test-client",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		TokenURL:     "https://example.com/token",
	}

	before := time.Now()
	err = repo.Save(ctx, client)
	after := time.Now()

	require.NoError(t, err)

	retrieved, err := repo.Get(ctx, "test-client")
	require.NoError(t, err)
	assert.True(t, retrieved.CreatedAt.After(before) || retrieved.CreatedAt.Equal(before))
	assert.True(t, retrieved.CreatedAt.Before(after) || retrieved.CreatedAt.Equal(after))
}

func TestVaultRepository_Get_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")
	repo := NewVaultRepository(vaultPath)
	ctx := context.Background()

	err := repo.Load(ctx, "password")
	require.NoError(t, err)

	client, err := repo.Get(ctx, "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "not found")
}

func TestVaultRepository_List(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")
	repo := NewVaultRepository(vaultPath)
	ctx := context.Background()

	err := repo.Load(ctx, "password")
	require.NoError(t, err)

	names, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, names)

	clients := []core.Client{
		{Name: "client1", ClientID: "id1", ClientSecret: "secret1", TokenURL: "url1"},
		{Name: "client2", ClientID: "id2", ClientSecret: "secret2", TokenURL: "url2"},
		{Name: "client3", ClientID: "id3", ClientSecret: "secret3", TokenURL: "url3"},
	}

	for _, client := range clients {
		err := repo.Save(ctx, client)
		require.NoError(t, err)
	}

	names, err = repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, names, 3)
	assert.Contains(t, names, "client1")
	assert.Contains(t, names, "client2")
	assert.Contains(t, names, "client3")
}

func TestVaultRepository_GetAll(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")
	repo := NewVaultRepository(vaultPath)
	ctx := context.Background()

	err := repo.Load(ctx, "password")
	require.NoError(t, err)

	clients, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Empty(t, clients)

	testClients := []core.Client{
		{Name: "client1", ClientID: "id1", ClientSecret: "secret1", TokenURL: "url1", Scopes: []string{"read"}},
		{Name: "client2", ClientID: "id2", ClientSecret: "secret2", TokenURL: "url2", Scopes: []string{"write"}},
	}

	for _, client := range testClients {
		err := repo.Save(ctx, client)
		require.NoError(t, err)
	}

	clients, err = repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, clients, 2)

	for i, expected := range testClients {
		assert.Equal(t, expected.Name, clients[i].Name)
		assert.Equal(t, expected.ClientID, clients[i].ClientID)
		assert.Equal(t, expected.ClientSecret, clients[i].ClientSecret)
		assert.Equal(t, expected.TokenURL, clients[i].TokenURL)
		assert.Equal(t, expected.Scopes, clients[i].Scopes)
	}
}

func TestVaultRepository_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")
	repo := NewVaultRepository(vaultPath)
	ctx := context.Background()

	err := repo.Load(ctx, "password")
	require.NoError(t, err)

	clients := []core.Client{
		{Name: "client1", ClientID: "id1", ClientSecret: "secret1", TokenURL: "url1"},
		{Name: "client2", ClientID: "id2", ClientSecret: "secret2", TokenURL: "url2"},
	}

	for _, client := range clients {
		err := repo.Save(ctx, client)
		require.NoError(t, err)
	}

	err = repo.Delete(ctx, "client1")
	require.NoError(t, err)

	names, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, names, 1)
	assert.Equal(t, "client2", names[0])

	_, err = repo.Get(ctx, "client1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestVaultRepository_Delete_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")
	repo := NewVaultRepository(vaultPath)
	ctx := context.Background()

	err := repo.Load(ctx, "password")
	require.NoError(t, err)

	err = repo.Delete(ctx, "nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestVaultRepository_WrongPassword(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")
	ctx := context.Background()

	repo1 := NewVaultRepository(vaultPath)
	err := repo1.Load(ctx, "correct-password")
	require.NoError(t, err)

	client := core.Client{
		Name:         "test-client",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		TokenURL:     "https://example.com/token",
	}

	err = repo1.Save(ctx, client)
	require.NoError(t, err)

	repo2 := NewVaultRepository(vaultPath)
	err = repo2.Load(ctx, "wrong-password")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decrypt vault")
}

func TestEncryptDecrypt(t *testing.T) {
	plaintext := []byte("test data")
	password := "test-password"

	ciphertext, err := encrypt(plaintext, password)
	require.NoError(t, err)
	assert.NotEmpty(t, ciphertext)
	assert.NotEqual(t, plaintext, ciphertext)

	decrypted, err := decrypt(ciphertext, password)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestDecrypt_WrongPassword(t *testing.T) {
	plaintext := []byte("test data")
	password := "correct-password"

	ciphertext, err := encrypt(plaintext, password)
	require.NoError(t, err)

	_, err = decrypt(ciphertext, "wrong-password")
	assert.Error(t, err)
}

func TestDecrypt_ShortCiphertext(t *testing.T) {
	shortData := []byte("short")
	_, err := decrypt(shortData, "password")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ciphertext too short")
}

func TestDecrypt_CorruptedData(t *testing.T) {
	plaintext := []byte("test data")
	password := "test-password"

	ciphertext, err := encrypt(plaintext, password)
	require.NoError(t, err)

	ciphertext[len(ciphertext)-1] ^= 0xFF

	_, err = decrypt(ciphertext, password)
	assert.Error(t, err)
}

func TestToClientData_ToClient(t *testing.T) {
	now := time.Now()
	client := core.Client{
		Name:         "test-client",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		TokenURL:     "https://example.com/token",
		Scopes:       []string{"read", "write"},
		CreatedAt:    now,
	}

	data := toClientData(client)

	assert.Equal(t, client.Name, data.Name)
	assert.Equal(t, client.ClientID, data.ClientID)
	assert.Equal(t, client.ClientSecret, data.ClientSecret)
	assert.Equal(t, client.TokenURL, data.TokenURL)
	assert.Equal(t, client.Scopes, data.Scopes)
	assert.Equal(t, client.CreatedAt, data.CreatedAt)

	converted := toClient(data)

	assert.Equal(t, client, converted)
}

func TestVaultRepository_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")
	repo := NewVaultRepository(vaultPath)
	ctx := context.Background()

	err := repo.Load(ctx, "password")
	require.NoError(t, err)

	err = os.WriteFile(vaultPath, []byte(""), 0600)
	require.NoError(t, err)

	clients, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, clients)
}

func TestVaultRepository_MultipleOperations(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "vault.enc")
	repo := NewVaultRepository(vaultPath)
	ctx := context.Background()

	err := repo.Load(ctx, "password")
	require.NoError(t, err)

	err = repo.Save(ctx, core.Client{Name: "client1", ClientID: "id1", ClientSecret: "s1", TokenURL: "u1"})
	require.NoError(t, err)

	err = repo.Save(ctx, core.Client{Name: "client2", ClientID: "id2", ClientSecret: "s2", TokenURL: "u2"})
	require.NoError(t, err)

	err = repo.Delete(ctx, "client1")
	require.NoError(t, err)

	err = repo.Save(ctx, core.Client{Name: "client3", ClientID: "id3", ClientSecret: "s3", TokenURL: "u3"})
	require.NoError(t, err)

	names, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Contains(t, names, "client2")
	assert.Contains(t, names, "client3")
	assert.NotContains(t, names, "client1")
}

func TestVaultRepository_CreatesDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	vaultPath := filepath.Join(tmpDir, "nested", "deep", "directory", "vault.enc")
	repo := NewVaultRepository(vaultPath)
	ctx := context.Background()

	dirPath := filepath.Join(tmpDir, "nested", "deep", "directory")
	_, err := os.Stat(dirPath)
	assert.True(t, os.IsNotExist(err), "directory should not exist before save")

	err = repo.Load(ctx, "password")
	require.NoError(t, err)

	client := core.Client{
		Name:         "test-client",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		TokenURL:     "https://example.com/token",
	}

	err = repo.Save(ctx, client)
	require.NoError(t, err)

	_, err = os.Stat(dirPath)
	assert.NoError(t, err, "directory should exist after save")

	info, err := os.Stat(vaultPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm(), "vault file should have 0600 permissions")

	dirInfo, err := os.Stat(dirPath)
	require.NoError(t, err)
	assert.True(t, dirInfo.IsDir(), "should be a directory")
	assert.Equal(t, os.FileMode(0700), dirInfo.Mode().Perm(), "directory should have 0700 permissions")

	retrieved, err := repo.Get(ctx, "test-client")
	require.NoError(t, err)
	assert.Equal(t, client.Name, retrieved.Name)
}
