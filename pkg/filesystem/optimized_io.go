package filesystem

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// BufferedFileWriter provides buffered file writing for better performance
type BufferedFileWriter struct {
	file   *os.File
	writer *bufio.Writer
	mutex  sync.Mutex
}

// NewBufferedFileWriter creates a new buffered file writer
func NewBufferedFileWriter(path string, perm os.FileMode) (*BufferedFileWriter, error) {
	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create parent directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return &BufferedFileWriter{
		file:   file,
		writer: bufio.NewWriterSize(file, 64*1024), // 64KB buffer
	}, nil
}

// Write writes data to the buffered writer
func (bfw *BufferedFileWriter) Write(data []byte) (int, error) {
	bfw.mutex.Lock()
	defer bfw.mutex.Unlock()
	return bfw.writer.Write(data)
}

// WriteString writes a string to the buffered writer
func (bfw *BufferedFileWriter) WriteString(s string) (int, error) {
	bfw.mutex.Lock()
	defer bfw.mutex.Unlock()
	return bfw.writer.WriteString(s)
}

// Flush flushes the buffer to disk
func (bfw *BufferedFileWriter) Flush() error {
	bfw.mutex.Lock()
	defer bfw.mutex.Unlock()
	return bfw.writer.Flush()
}

// Close flushes and closes the file
func (bfw *BufferedFileWriter) Close() error {
	bfw.mutex.Lock()
	defer bfw.mutex.Unlock()

	if err := bfw.writer.Flush(); err != nil {
		bfw.file.Close()
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	return bfw.file.Close()
}

// BatchFileOperations provides efficient batch file operations
type BatchFileOperations struct {
	operations []FileOperation
	mutex      sync.Mutex
}

// FileOperation represents a file operation to be batched
type FileOperation struct {
	Type        string // "write", "copy", "mkdir"
	Source      string
	Destination string
	Content     []byte
	Permissions os.FileMode
}

// NewBatchFileOperations creates a new batch file operations manager
func NewBatchFileOperations() *BatchFileOperations {
	return &BatchFileOperations{
		operations: make([]FileOperation, 0),
	}
}

// AddWrite adds a write operation to the batch
func (bfo *BatchFileOperations) AddWrite(path string, content []byte, perm os.FileMode) {
	bfo.mutex.Lock()
	defer bfo.mutex.Unlock()

	bfo.operations = append(bfo.operations, FileOperation{
		Type:        "write",
		Destination: path,
		Content:     content,
		Permissions: perm,
	})
}

// AddCopy adds a copy operation to the batch
func (bfo *BatchFileOperations) AddCopy(src, dest string, perm os.FileMode) {
	bfo.mutex.Lock()
	defer bfo.mutex.Unlock()

	bfo.operations = append(bfo.operations, FileOperation{
		Type:        "copy",
		Source:      src,
		Destination: dest,
		Permissions: perm,
	})
}

// AddMkdir adds a directory creation operation to the batch
func (bfo *BatchFileOperations) AddMkdir(path string, perm os.FileMode) {
	bfo.mutex.Lock()
	defer bfo.mutex.Unlock()

	bfo.operations = append(bfo.operations, FileOperation{
		Type:        "mkdir",
		Destination: path,
		Permissions: perm,
	})
}

// Execute executes all batched operations efficiently
func (bfo *BatchFileOperations) Execute() error {
	bfo.mutex.Lock()
	defer bfo.mutex.Unlock()

	// Group operations by type for better efficiency
	mkdirOps := make([]FileOperation, 0)
	writeOps := make([]FileOperation, 0)
	copyOps := make([]FileOperation, 0)

	for _, op := range bfo.operations {
		switch op.Type {
		case "mkdir":
			mkdirOps = append(mkdirOps, op)
		case "write":
			writeOps = append(writeOps, op)
		case "copy":
			copyOps = append(copyOps, op)
		}
	}

	// Execute mkdir operations first
	for _, op := range mkdirOps {
		if err := os.MkdirAll(op.Destination, op.Permissions); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", op.Destination, err)
		}
	}

	// Execute write operations
	for _, op := range writeOps {
		if err := bfo.executeWrite(op); err != nil {
			return err
		}
	}

	// Execute copy operations
	for _, op := range copyOps {
		if err := bfo.executeCopy(op); err != nil {
			return err
		}
	}

	// Clear operations after successful execution
	bfo.operations = bfo.operations[:0]

	return nil
}

// executeWrite executes a write operation
func (bfo *BatchFileOperations) executeWrite(op FileOperation) error {
	// Ensure parent directory exists
	dir := filepath.Dir(op.Destination)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory for %s: %w", op.Destination, err)
	}

	return os.WriteFile(op.Destination, op.Content, op.Permissions)
}

// executeCopy executes a copy operation
func (bfo *BatchFileOperations) executeCopy(op FileOperation) error {
	srcFile, err := os.Open(op.Source)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", op.Source, err)
	}
	defer srcFile.Close()

	// Ensure parent directory exists
	dir := filepath.Dir(op.Destination)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory for %s: %w", op.Destination, err)
	}

	destFile, err := os.OpenFile(op.Destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, op.Permissions)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", op.Destination, err)
	}
	defer destFile.Close()

	// Use buffered copy for better performance
	_, err = io.CopyBuffer(destFile, srcFile, make([]byte, 64*1024))
	if err != nil {
		return fmt.Errorf("failed to copy file from %s to %s: %w", op.Source, op.Destination, err)
	}

	return nil
}

// Clear clears all pending operations
func (bfo *BatchFileOperations) Clear() {
	bfo.mutex.Lock()
	defer bfo.mutex.Unlock()
	bfo.operations = bfo.operations[:0]
}

// Count returns the number of pending operations
func (bfo *BatchFileOperations) Count() int {
	bfo.mutex.Lock()
	defer bfo.mutex.Unlock()
	return len(bfo.operations)
}

// OptimizedCopy provides optimized file copying with larger buffers
func OptimizedCopy(src, dest string, perm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Get source file info for size optimization
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	// Ensure parent directory exists
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Use larger buffer for better performance on large files
	bufferSize := 64 * 1024         // 64KB default
	if srcInfo.Size() > 1024*1024 { // If file > 1MB, use 256KB buffer
		bufferSize = 256 * 1024
	}

	buffer := make([]byte, bufferSize)
	_, err = io.CopyBuffer(destFile, srcFile, buffer)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// ParallelDirectoryWalk provides parallel directory walking for better performance
type ParallelDirectoryWalk struct {
	maxWorkers int
	workChan   chan string
	resultChan chan WalkResult
	wg         sync.WaitGroup
}

// WalkResult represents the result of walking a directory entry
type WalkResult struct {
	Path  string
	Info  os.FileInfo
	Error error
}

// NewParallelDirectoryWalk creates a new parallel directory walker
func NewParallelDirectoryWalk(maxWorkers int) *ParallelDirectoryWalk {
	return &ParallelDirectoryWalk{
		maxWorkers: maxWorkers,
		workChan:   make(chan string, maxWorkers*2),
		resultChan: make(chan WalkResult, maxWorkers*2),
	}
}

// Walk walks a directory in parallel and calls the provided function for each entry
func (pdw *ParallelDirectoryWalk) Walk(root string, walkFn func(path string, info os.FileInfo, err error) error) error {
	// Start workers
	for i := 0; i < pdw.maxWorkers; i++ {
		pdw.wg.Add(1)
		go pdw.worker()
	}

	// Start result processor
	done := make(chan error, 1)
	go func() {
		defer close(done)
		for result := range pdw.resultChan {
			if err := walkFn(result.Path, result.Info, result.Error); err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()

	// Start directory traversal
	pdw.workChan <- root
	pdw.wg.Wait()
	close(pdw.resultChan)

	return <-done
}

// worker processes directory entries
func (pdw *ParallelDirectoryWalk) worker() {
	defer pdw.wg.Done()

	for path := range pdw.workChan {
		info, err := os.Stat(path)
		if err != nil {
			pdw.resultChan <- WalkResult{Path: path, Error: err}
			continue
		}

		pdw.resultChan <- WalkResult{Path: path, Info: info}

		if info.IsDir() {
			entries, err := os.ReadDir(path)
			if err != nil {
				pdw.resultChan <- WalkResult{Path: path, Error: err}
				continue
			}

			for _, entry := range entries {
				entryPath := filepath.Join(path, entry.Name())
				select {
				case pdw.workChan <- entryPath:
					pdw.wg.Add(1)
				default:
					// Channel full, add to waitgroup before goroutine
					pdw.wg.Add(1)
					go func(p string) {
						pdw.workChan <- p
					}(entryPath)
				}
			}
		}
	}
}
