package orchestrator

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// ProgressIndicator provides real-time progress feedback for long operations
type ProgressIndicator struct {
	writer      io.Writer
	message     string
	spinner     []string
	spinnerIdx  int
	done        chan bool
	mu          sync.Mutex
	active      bool
	showSpinner bool
}

// NewProgressIndicator creates a new progress indicator
func NewProgressIndicator(writer io.Writer, message string, showSpinner bool) *ProgressIndicator {
	return &ProgressIndicator{
		writer:      writer,
		message:     message,
		spinner:     []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		spinnerIdx:  0,
		done:        make(chan bool),
		showSpinner: showSpinner,
	}
}

// Start begins showing the progress indicator
func (pi *ProgressIndicator) Start() {
	pi.mu.Lock()
	if pi.active {
		pi.mu.Unlock()
		return
	}
	pi.active = true
	pi.mu.Unlock()

	if !pi.showSpinner {
		// Just print the message once
		_, _ = fmt.Fprintf(pi.writer, "%s\n", pi.message) // Error intentionally ignored for output formatting
		return
	}

	// Start spinner animation
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-pi.done:
				return
			case <-ticker.C:
				pi.mu.Lock()
				if pi.active {
					// Clear line and print spinner
					_, _ = fmt.Fprintf(pi.writer, "\r%s %s", pi.spinner[pi.spinnerIdx], pi.message) // Error intentionally ignored for output formatting
					pi.spinnerIdx = (pi.spinnerIdx + 1) % len(pi.spinner)
				}
				pi.mu.Unlock()
			}
		}
	}()
}

// Stop stops the progress indicator
func (pi *ProgressIndicator) Stop() {
	pi.mu.Lock()
	defer pi.mu.Unlock()

	if !pi.active {
		return
	}

	pi.active = false
	if pi.showSpinner {
		close(pi.done)
		// Clear the spinner line
		_, _ = fmt.Fprintf(pi.writer, "\r%s\r", strings.Repeat(" ", len(pi.message)+3)) // Error intentionally ignored for output formatting
	}
}

// Update updates the progress message
func (pi *ProgressIndicator) Update(message string) {
	pi.mu.Lock()
	defer pi.mu.Unlock()

	pi.message = message
	if !pi.showSpinner && pi.active {
		_, _ = fmt.Fprintf(pi.writer, "%s\n", message) // Error intentionally ignored for output formatting
	}
}

// StreamingWriter wraps an io.Writer to provide line-by-line output with prefixes
type StreamingWriter struct {
	writer     io.Writer
	prefix     string
	buffer     []byte
	mu         sync.Mutex
	lineCount  int
	maxLines   int
	showPrefix bool
}

// NewStreamingWriter creates a new streaming writer
func NewStreamingWriter(writer io.Writer, prefix string, showPrefix bool) *StreamingWriter {
	return &StreamingWriter{
		writer:     writer,
		prefix:     prefix,
		buffer:     make([]byte, 0),
		maxLines:   100, // Limit output to prevent overwhelming the terminal
		showPrefix: showPrefix,
	}
}

// Write implements io.Writer interface
func (sw *StreamingWriter) Write(p []byte) (n int, err error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	// Add to buffer
	sw.buffer = append(sw.buffer, p...)

	// Process complete lines
	for {
		idx := -1
		for i, b := range sw.buffer {
			if b == '\n' {
				idx = i
				break
			}
		}

		if idx == -1 {
			// No complete line yet
			break
		}

		// Extract line
		line := sw.buffer[:idx]
		sw.buffer = sw.buffer[idx+1:]

		// Check line limit
		if sw.maxLines > 0 && sw.lineCount >= sw.maxLines {
			if sw.lineCount == sw.maxLines {
				_, _ = fmt.Fprintf(sw.writer, "%s[Output truncated - too many lines]\n", sw.getPrefix()) // Error intentionally ignored for output formatting
				sw.lineCount++
			}
			continue
		}

		// Write line with prefix
		if len(line) > 0 {
			_, _ = fmt.Fprintf(sw.writer, "%s%s\n", sw.getPrefix(), string(line)) // Error intentionally ignored for output formatting
			sw.lineCount++
		}
	}

	return len(p), nil
}

// Flush writes any remaining buffered data
func (sw *StreamingWriter) Flush() {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if len(sw.buffer) > 0 && sw.lineCount < sw.maxLines {
		_, _ = fmt.Fprintf(sw.writer, "%s%s\n", sw.getPrefix(), string(sw.buffer)) // Error intentionally ignored for output formatting
		sw.buffer = sw.buffer[:0]
		sw.lineCount++
	}
}

// GetLineCount returns the number of lines written
func (sw *StreamingWriter) GetLineCount() int {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.lineCount
}

// getPrefix returns the prefix string if enabled
func (sw *StreamingWriter) getPrefix() string {
	if sw.showPrefix && sw.prefix != "" {
		return fmt.Sprintf("[%s] ", sw.prefix)
	}
	return ""
}

// ProgressTracker tracks progress across multiple operations
type ProgressTracker struct {
	writer       io.Writer
	total        int
	completed    int
	mu           sync.Mutex
	startTime    time.Time
	showProgress bool
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(writer io.Writer, total int, showProgress bool) *ProgressTracker {
	return &ProgressTracker{
		writer:       writer,
		total:        total,
		completed:    0,
		startTime:    time.Now(),
		showProgress: showProgress,
	}
}

// Increment increments the completed count and updates progress
func (pt *ProgressTracker) Increment(itemName string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	pt.completed++

	if pt.showProgress {
		elapsed := time.Since(pt.startTime)
		percentage := float64(pt.completed) / float64(pt.total) * 100

		// Estimate remaining time
		var eta string
		if pt.completed > 0 {
			avgTime := elapsed / time.Duration(pt.completed)
			remaining := time.Duration(pt.total-pt.completed) * avgTime
			eta = fmt.Sprintf(" (ETA: %v)", remaining.Round(time.Second))
		}

		_, _ = fmt.Fprintf(pt.writer, "\r[%d/%d] %.0f%% - %s%s",
			pt.completed, pt.total, percentage, itemName, eta) // Error intentionally ignored for output formatting

		if pt.completed == pt.total {
			_, _ = fmt.Fprintf(pt.writer, "\n") // Error intentionally ignored for output formatting
		}
	}
}

// Complete marks all operations as complete
func (pt *ProgressTracker) Complete() {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if pt.showProgress && pt.completed < pt.total {
		pt.completed = pt.total
		elapsed := time.Since(pt.startTime)
		_, _ = fmt.Fprintf(pt.writer, "\r[%d/%d] 100%% - Completed in %v\n",
			pt.total, pt.total, elapsed.Round(time.Second)) // Error intentionally ignored for output formatting
	}
}
