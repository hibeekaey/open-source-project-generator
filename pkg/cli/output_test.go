package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simple tests without complex mocking
func TestColorManager_Basic(t *testing.T) {
	// Test with colors enabled
	cm := NewColorManager(true)
	assert.NotNil(t, cm)
	assert.True(t, cm.enabled)

	// Test colorize method
	result := cm.Colorize(ColorRed, "test")
	expected := ColorRed + "test" + ColorReset
	assert.Equal(t, expected, result)

	// Test with colors disabled
	cmDisabled := NewColorManager(false)
	assert.False(t, cmDisabled.enabled)

	result = cmDisabled.Colorize(ColorRed, "test")
	assert.Equal(t, "test", result)
}

func TestColorManager_Methods(t *testing.T) {
	cm := NewColorManager(true)

	// Test all color methods return non-empty strings
	assert.Contains(t, cm.Success("test"), "test")
	assert.Contains(t, cm.Info("test"), "test")
	assert.Contains(t, cm.Warning("test"), "test")
	assert.Contains(t, cm.Error("test"), "test")
	assert.Contains(t, cm.Highlight("test"), "test")
	assert.Contains(t, cm.Dim("test"), "test")
}

func TestOutputManager_Basic(t *testing.T) {
	// Test with nil logger (should not panic)
	om := NewOutputManager(false, false, false, nil)
	assert.NotNil(t, om)
	assert.NotNil(t, om.colorizer)

	// Test mode getters
	assert.False(t, om.IsVerboseMode())
	assert.False(t, om.IsQuietMode())
	assert.False(t, om.IsDebugMode())

	// Test mode setters
	om.SetVerboseMode(true)
	assert.True(t, om.IsVerboseMode())

	om.SetQuietMode(true)
	assert.True(t, om.IsQuietMode())

	om.SetDebugMode(true)
	assert.True(t, om.IsDebugMode())
}

func TestOutputManager_OutputMethods(t *testing.T) {
	// Test that output methods don't panic with nil logger
	om := NewOutputManager(true, false, true, nil)

	// These should not panic
	om.VerboseOutput("test")
	om.DebugOutput("test")
	om.QuietOutput("test")
	om.ErrorOutput("test")
	om.WarningOutput("test")
	om.SuccessOutput("test")

	// Test with metrics
	metrics := map[string]interface{}{"test": "value"}
	om.PerformanceOutput("test op", 0, metrics)
}

func TestOutputManager_ColorManager(t *testing.T) {
	om := NewOutputManager(false, false, false, nil)
	cm := om.GetColorManager()

	assert.NotNil(t, cm)
	assert.Equal(t, om.colorizer, cm)
}
