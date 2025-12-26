package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultVaultPath(t *testing.T) {
	path, err := getDefaultVaultPath()

	assert.NoError(t, err)
	assert.NotEmpty(t, path)
	assert.Contains(t, path, ".authkeeper")
	assert.Contains(t, path, "vault.enc")
}
