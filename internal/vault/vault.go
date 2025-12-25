package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ksysoev/authkeeper/pkg/crypto"
)

type Client struct {
	Name         string    `json:"name"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	TokenURL     string    `json:"token_url"`
	Scopes       []string  `json:"scopes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type VaultData struct {
	Clients []Client `json:"clients"`
}

type Vault struct {
	path string
}

func New(path string) *Vault {
	return &Vault{path: path}
}

// GetDefaultVaultPath returns the default vault path in user's home directory
func GetDefaultVaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	vaultDir := filepath.Join(home, ".authkeeper")
	if err := os.MkdirAll(vaultDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create vault directory: %w", err)
	}

	return filepath.Join(vaultDir, "vault.enc"), nil
}

// Load loads and decrypts the vault
func (v *Vault) Load(password string) (*VaultData, error) {
	data, err := os.ReadFile(v.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &VaultData{Clients: []Client{}}, nil
		}
		return nil, fmt.Errorf("failed to read vault: %w", err)
	}

	if len(data) == 0 {
		return &VaultData{Clients: []Client{}}, nil
	}

	plaintext, err := crypto.Decrypt(data, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault (wrong password?): %w", err)
	}

	var vaultData VaultData
	if err := json.Unmarshal(plaintext, &vaultData); err != nil {
		return nil, fmt.Errorf("failed to parse vault data: %w", err)
	}

	return &vaultData, nil
}

// Save encrypts and saves the vault
func (v *Vault) Save(data *VaultData, password string) error {
	plaintext, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vault data: %w", err)
	}

	ciphertext, err := crypto.Encrypt(plaintext, password)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}

	if err := os.WriteFile(v.path, ciphertext, 0600); err != nil {
		return fmt.Errorf("failed to write vault: %w", err)
	}

	return nil
}

// AddClient adds a new client to the vault
func (v *Vault) AddClient(client Client, password string) error {
	data, err := v.Load(password)
	if err != nil {
		return err
	}

	// Check for duplicate names
	for _, c := range data.Clients {
		if c.Name == client.Name {
			return fmt.Errorf("client with name %q already exists", client.Name)
		}
	}

	client.CreatedAt = time.Now()
	data.Clients = append(data.Clients, client)

	return v.Save(data, password)
}

// DeleteClient removes a client from the vault
func (v *Vault) DeleteClient(name string, password string) error {
	data, err := v.Load(password)
	if err != nil {
		return err
	}

	for i, c := range data.Clients {
		if c.Name == name {
			data.Clients = append(data.Clients[:i], data.Clients[i+1:]...)
			return v.Save(data, password)
		}
	}

	return fmt.Errorf("client %q not found", name)
}

// GetClient retrieves a client by name
func (v *Vault) GetClient(name string, password string) (*Client, error) {
	data, err := v.Load(password)
	if err != nil {
		return nil, err
	}

	for _, c := range data.Clients {
		if c.Name == name {
			return &c, nil
		}
	}

	return nil, fmt.Errorf("client %q not found", name)
}

// ListClients returns all client names
func (v *Vault) ListClients(password string) ([]string, error) {
	data, err := v.Load(password)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(data.Clients))
	for i, c := range data.Clients {
		names[i] = c.Name
	}

	return names, nil
}
