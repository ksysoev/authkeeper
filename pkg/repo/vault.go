package repo

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ksysoev/authkeeper/pkg/core"
	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 32
	iterations = 100000
	keySize    = 32
)

type vaultData struct {
	Clients []clientData `json:"clients"`
}

type clientData struct {
	Name         string    `json:"name"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	TokenURL     string    `json:"token_url"`
	Scopes       []string  `json:"scopes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// VaultRepository implements core.Repository interface using encrypted file storage
type VaultRepository struct {
	path     string
	password string
}

// NewVaultRepository creates a new vault repository
func NewVaultRepository(path string) *VaultRepository {
	return &VaultRepository{
		path: path,
	}
}

// Exists checks if the vault file exists
func (r *VaultRepository) Exists() bool {
	_, err := os.Stat(r.path)
	return err == nil
}

// Load unlocks the vault with the given password
func (r *VaultRepository) Load(_ context.Context, password string) error {
	r.password = password
	_, err := r.load()
	return err
}

// Save stores a client in the vault
func (r *VaultRepository) Save(ctx context.Context, client core.Client) error {
	data, err := r.load()
	if err != nil {
		return err
	}

	// Check for duplicate names
	for _, c := range data.Clients {
		if c.Name == client.Name {
			return fmt.Errorf("client with name %q already exists", client.Name)
		}
	}

	// Add timestamp if not set
	if client.CreatedAt.IsZero() {
		client.CreatedAt = time.Now()
	}

	data.Clients = append(data.Clients, toClientData(client))

	return r.save(data)
}

// Get retrieves a client by name
func (r *VaultRepository) Get(ctx context.Context, name string) (*core.Client, error) {
	data, err := r.load()
	if err != nil {
		return nil, err
	}

	for _, c := range data.Clients {
		if c.Name == name {
			client := toClient(c)
			return &client, nil
		}
	}

	return nil, fmt.Errorf("client %q not found", name)
}

// List returns all client names
func (r *VaultRepository) List(ctx context.Context) ([]string, error) {
	data, err := r.load()
	if err != nil {
		return nil, err
	}

	names := make([]string, len(data.Clients))
	for i, c := range data.Clients {
		names[i] = c.Name
	}

	return names, nil
}

// GetAll returns all clients with full details
func (r *VaultRepository) GetAll(ctx context.Context) ([]core.Client, error) {
	data, err := r.load()
	if err != nil {
		return nil, err
	}

	clients := make([]core.Client, len(data.Clients))
	for i, c := range data.Clients {
		clients[i] = toClient(c)
	}

	return clients, nil
}

// Delete removes a client by name
func (r *VaultRepository) Delete(ctx context.Context, name string) error {
	data, err := r.load()
	if err != nil {
		return err
	}

	for i, c := range data.Clients {
		if c.Name == name {
			data.Clients = append(data.Clients[:i], data.Clients[i+1:]...)
			return r.save(data)
		}
	}

	return fmt.Errorf("client %q not found", name)
}

// load loads and decrypts the vault
func (r *VaultRepository) load() (*vaultData, error) {
	fileData, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &vaultData{Clients: []clientData{}}, nil
		}
		return nil, fmt.Errorf("failed to read vault: %w", err)
	}

	if len(fileData) == 0 {
		return &vaultData{Clients: []clientData{}}, nil
	}

	plaintext, err := decrypt(fileData, r.password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault (wrong password?): %w", err)
	}

	var data vaultData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, fmt.Errorf("failed to parse vault data: %w", err)
	}

	return &data, nil
}

// save encrypts and saves the vault
func (r *VaultRepository) save(data *vaultData) error {
	plaintext, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vault data: %w", err)
	}

	ciphertext, err := encrypt(plaintext, r.password)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}

	dir := filepath.Dir(r.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	if err := os.WriteFile(r.path, ciphertext, 0600); err != nil {
		return fmt.Errorf("failed to write vault: %w", err)
	}

	return nil
}

// encrypt encrypts plaintext using AES-256-GCM with a password-derived key
func encrypt(plaintext []byte, password string) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	key := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Prepend salt to ciphertext
	result := make([]byte, len(salt)+len(ciphertext))
	copy(result, salt)
	copy(result[len(salt):], ciphertext)

	return result, nil
}

// decrypt decrypts ciphertext using AES-256-GCM with a password-derived key
func decrypt(ciphertext []byte, password string) ([]byte, error) {
	if len(ciphertext) < saltSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	salt := ciphertext[:saltSize]
	ciphertext = ciphertext[saltSize:]

	key := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// Helper functions to convert between core.Client and clientData
func toClientData(c core.Client) clientData {
	return clientData{
		Name:         c.Name,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     c.TokenURL,
		Scopes:       c.Scopes,
		CreatedAt:    c.CreatedAt,
	}
}

func toClient(c clientData) core.Client {
	return core.Client{
		Name:         c.Name,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     c.TokenURL,
		Scopes:       c.Scopes,
		CreatedAt:    c.CreatedAt,
	}
}
