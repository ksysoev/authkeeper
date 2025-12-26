package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// getDefaultVaultPath returns the default path for the vault file.
// It returns a string representing the path and an error if the path cannot be determined.
func getDefaultVaultPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".authkeeper", "vault.enc"), nil
}
