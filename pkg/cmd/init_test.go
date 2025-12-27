package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCommands(t *testing.T) {
	version := "1.0.0"

	rootCmd, err := InitCommands(version)

	assert.NoError(t, err)
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "authkeeper", rootCmd.Use)
	assert.Equal(t, version, rootCmd.Version)
	assert.NotEmpty(t, rootCmd.Short)
	assert.NotEmpty(t, rootCmd.Long)

	subCommands := rootCmd.Commands()
	assert.Len(t, subCommands, 4)

	commandNames := make(map[string]bool)
	for _, cmd := range subCommands {
		// Extract the base command name (first word)
		cmdName := cmd.Name()
		commandNames[cmdName] = true
	}

	assert.True(t, commandNames["add"])
	assert.True(t, commandNames["token"])
	assert.True(t, commandNames["list"])
	assert.True(t, commandNames["delete"])
}

func TestAddCommand(t *testing.T) {
	args := &args{
		version:   "1.0.0",
		vaultPath: "/tmp/vault.enc",
	}

	cmd := AddCommand(args)

	assert.NotNil(t, cmd)
	assert.Equal(t, "add [flags]", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotNil(t, cmd.RunE)
}

func TestTokenCommand(t *testing.T) {
	args := &args{
		version:   "1.0.0",
		vaultPath: "/tmp/vault.enc",
	}

	cmd := TokenCommand(args)

	assert.NotNil(t, cmd)
	assert.Equal(t, "token [client-name]", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotNil(t, cmd.RunE)
}

func TestListCommand(t *testing.T) {
	args := &args{
		version:   "1.0.0",
		vaultPath: "/tmp/vault.enc",
	}

	cmd := ListCommand(args)

	assert.NotNil(t, cmd)
	assert.Equal(t, "list", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotNil(t, cmd.RunE)
}

func TestDeleteCommand(t *testing.T) {
	args := &args{
		version:   "1.0.0",
		vaultPath: "/tmp/vault.enc",
	}

	cmd := DeleteCommand(args)

	assert.NotNil(t, cmd)
	assert.Equal(t, "delete [client-name]", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotNil(t, cmd.RunE)
}
