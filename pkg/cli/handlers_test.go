package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommandHandlers(t *testing.T) {
	// Create a minimal CLI instance for testing
	cli := &CLI{}
	handlers := NewCommandHandlers(cli)

	assert.NotNil(t, handlers)
	assert.Equal(t, cli, handlers.cli)
}

func TestCommandHandlers_Structure(t *testing.T) {
	// Test that CommandHandlers has the expected structure
	cli := &CLI{}
	handlers := NewCommandHandlers(cli)

	// Verify the handlers struct is properly initialized
	assert.NotNil(t, handlers)
	assert.NotNil(t, handlers.cli)
}
