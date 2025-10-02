// Package cli provides adapter structs to make components compatible with interfaces.
package cli

// OutputAdapter adapts OutputManager to work with interactive and validation components
type OutputAdapter struct {
	outputManager *OutputManager
}

// NewOutputAdapter creates a new output adapter
func NewOutputAdapter(outputManager *OutputManager) *OutputAdapter {
	return &OutputAdapter{outputManager: outputManager}
}

// Methods for interactive.OutputInterface
func (oa *OutputAdapter) QuietOutput(format string, args ...interface{}) {
	oa.outputManager.QuietOutput(format, args...)
}
func (oa *OutputAdapter) VerboseOutput(format string, args ...interface{}) {
	oa.outputManager.VerboseOutput(format, args...)
}
func (oa *OutputAdapter) WarningOutput(format string, args ...interface{}) {
	oa.outputManager.WarningOutput(format, args...)
}
func (oa *OutputAdapter) Error(text string) string {
	return oa.outputManager.GetColorManager().Error(text)
}
func (oa *OutputAdapter) Info(text string) string {
	return oa.outputManager.GetColorManager().Info(text)
}
func (oa *OutputAdapter) Warning(text string) string {
	return oa.outputManager.GetColorManager().Warning(text)
}
func (oa *OutputAdapter) Success(text string) string {
	return oa.outputManager.GetColorManager().Success(text)
}
func (oa *OutputAdapter) Highlight(text string) string {
	return oa.outputManager.GetColorManager().Highlight(text)
}
func (oa *OutputAdapter) Dim(text string) string { return oa.outputManager.GetColorManager().Dim(text) }

// Methods for validation.OutputInterface
func (oa *OutputAdapter) DebugOutput(format string, args ...interface{}) {
	oa.outputManager.DebugOutput(format, args...)
}
