package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintTitle(t *testing.T) {
	assert.NotPanics(t, func() {
		printTitle("Test Title")
	})
}

func TestPrintSuccess(t *testing.T) {
	assert.NotPanics(t, func() {
		printSuccess("Success message")
	})
}

func TestPrintError(t *testing.T) {
	assert.NotPanics(t, func() {
		printError("Error message")
	})
}

func TestPrintInfo(t *testing.T) {
	assert.NotPanics(t, func() {
		printInfo("Info message")
	})
}

func TestPrintWarning(t *testing.T) {
	assert.NotPanics(t, func() {
		printWarning("Warning message")
	})
}

func TestPrintMuted(t *testing.T) {
	assert.NotPanics(t, func() {
		printMuted("Muted message")
	})
}

func TestPrintProgress(t *testing.T) {
	assert.NotPanics(t, func() {
		printProgress("Progress message")
	})
}
