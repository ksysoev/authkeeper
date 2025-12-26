package cmd

import (
	"testing"

	"github.com/ksysoev/authkeeper/pkg/ui"
	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	cli := &ui.CLI{}
	app := NewApp(cli)

	assert.NotNil(t, app)
	assert.Equal(t, cli, app.cli)
}

func TestApp_BuildRootCommand(t *testing.T) {
	cli := &ui.CLI{}
	app := NewApp(cli)
	version := "1.0.0"

	rootCmd := app.BuildRootCommand(version)

	assert.NotNil(t, rootCmd)
	assert.Equal(t, "authkeeper", rootCmd.Use)
	assert.Equal(t, version, rootCmd.Version)
	assert.NotEmpty(t, rootCmd.Short)
	assert.NotEmpty(t, rootCmd.Long)

	subCommands := rootCmd.Commands()
	assert.Len(t, subCommands, 4)

	commandNames := make(map[string]bool)
	for _, cmd := range subCommands {
		commandNames[cmd.Use] = true
	}

	assert.True(t, commandNames["add"])
	assert.True(t, commandNames["token"])
	assert.True(t, commandNames["list"])
	assert.True(t, commandNames["delete"])
}

func TestApp_AddCommand(t *testing.T) {
	cli := &ui.CLI{}
	app := NewApp(cli)

	cmd := app.addCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "add", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotNil(t, cmd.RunE)
}

func TestApp_TokenCommand(t *testing.T) {
	cli := &ui.CLI{}
	app := NewApp(cli)

	cmd := app.tokenCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "token", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotNil(t, cmd.RunE)
}

func TestApp_ListCommand(t *testing.T) {
	cli := &ui.CLI{}
	app := NewApp(cli)

	cmd := app.listCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "list", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotNil(t, cmd.RunE)
}

func TestApp_DeleteCommand(t *testing.T) {
	cli := &ui.CLI{}
	app := NewApp(cli)

	cmd := app.deleteCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "delete", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotNil(t, cmd.RunE)
}
