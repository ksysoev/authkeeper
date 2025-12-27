package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectFromList_NoItems(t *testing.T) {
	_, err := selectFromList("Test", []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no items to select from")
}

// Note: readLine, readPassword, confirm, and selectFromList (with items)
// require stdin interaction and are tested through integration tests

