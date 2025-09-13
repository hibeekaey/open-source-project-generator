package template

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/open-source-template-generator/pkg/filesystem"
	"github.com/open-source-template-generator/pkg/models"
)

// ParallelProcessor provides parallel template processing capabilities
type ParallelProcessor struct {
	engine     *Engine
	maxWorkers int
	batchOps   *filesystem.BatchFileOperations
	workerPool chan struct{}
	resultChan chan ProcessResult
	errorChan  chan error
	wg         sync.WaitGroup
}

// ProcessResult represents the result of processing a template
type ProcessResult struct {
	SourcePath      string
	DestinationPath string
	ProcessingTime  time.Duration
	Size            int64
	Error           error
}

// ProcessingStats provides statistics about template processing
type ProcessingStats struct {
	TotalFiles      int
	ProcessedFiles  int
	FailedFiles     int
	TotalSize       int64
	ProcessingTime  time.Duration
	AverageFileTime time.Duration
	Throughput      float64 // Files per second
}

// NewParallelProcessor creates a new parallel template processor
func NewParallelProcessor(engine *Engine, maxWorkers int) *ParallelProcessor {
	if maxWorkers <= 0 {
		maxWorkers = runtime.NumCPU()
	}

	return &ParallelProcessor{
		engine:     engine,
		maxWorkers: maxWorkers,
		batchOps:   filesystem.NewBatchFileOperations(),
		workerPool: make(chan struct{}, maxWorkers),
		resultChan: make(chan ProcessResult, maxWorkers*2),
		errorChan:  make(chan error, maxWorkers),
	}
}

// ProcessDirectoryParallel processes a template directory using parallel workers
func (pp *ParallelProcessor) ProcessDirectoryParallel(ctx context.Context, templateDir string, outputDir string, config *models.ProjectConfig) (*ProcessingStats, error) {
	startTime := time.Now()

	// Enhance config with versions once for all templates
	enhancedConfig, err := pp.engine.enhanceConfigWithVersions(config)
	if err != nil {
		return nil, fmt.Errorf("failed to enhance config with versions: %w", err)
	}

	// Collect all template files first
	templateFiles, err := pp.collectTemplateFiles(templateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to collect template files: %w", err)
	}

	stats := &ProcessingStats{
		TotalFiles: len(templateFiles),
	}

	// Start result collector
	go pp.collectResults(stats)

	// Process files in parallel
	for _, file := range templateFiles {
		select {
		case <-ctx.Done():
			return stats, ctx.Err()
		default:
			pp.processFileAsync(ctx, file, templateDir, outputDir, enhancedConfig)
		}
	}

	// Wait for all workers to complete
	pp.wg.Wait()
	close(pp.resultChan)

	// Calculate final statistics
	stats.ProcessingTime = time.Since(startTime)
	if stats.ProcessedFiles > 0 {
		stats.AverageFileTime = stats.ProcessingTime / time.Duration(stats.ProcessedFiles)
		stats.Throughput = float64(stats.ProcessedFiles) / stats.ProcessingTime.Seconds()
	}

	return stats, nil
}

// processFileAsync processes a single file asynchronously
func (pp *ParallelProcessor) processFileAsync(ctx context.Context, filePath, templateDir, outputDir string, config *models.ProjectConfig) {
	pp.wg.Add(1)

	go func() {
		defer pp.wg.Done()

		// Acquire worker slot
		pp.workerPool <- struct{}{}
		defer func() { <-pp.workerPool }()

		result := pp.processFile(ctx, filePath, templateDir, outputDir, config)

		select {
		case pp.resultChan <- result:
		case <-ctx.Done():
		}
	}()
}

// processFile processes a single template file
func (pp *ParallelProcessor) processFile(ctx context.Context, srcPath, templateDir, outputDir string, config *models.ProjectConfig) ProcessResult {
	startTime := time.Now()

	result := ProcessResult{
		SourcePath: srcPath,
	}

	// Calculate relative path from template directory
	relPath, err := filepath.Rel(templateDir, srcPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to get relative path: %w", err)
		return result
	}

	// Calculate output path
	outputPath := filepath.Join(outputDir, relPath)
	result.DestinationPath = outputPath

	// Get file info for size
	info, err := os.Stat(srcPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to get file info: %w", err)
		return result
	}
	result.Size = info.Size()

	// Check for context cancellation
	select {
	case <-ctx.Done():
		result.Error = ctx.Err()
		return result
	default:
	}

	// Process the file
	if strings.HasSuffix(srcPath, ".tmpl") {
		// Remove .tmpl extension from destination
		outputPath = strings.TrimSuffix(outputPath, ".tmpl")
		result.DestinationPath = outputPath

		// Process as template
		content, err := pp.engine.ProcessTemplate(srcPath, config)
		if err != nil {
			result.Error = fmt.Errorf("failed to process template: %w", err)
			return result
		}

		// Use batch operations for better performance
		pp.batchOps.AddWrite(outputPath, content, 0644)
	} else {
		// Copy binary file
		pp.batchOps.AddCopy(srcPath, outputPath, info.Mode())
	}

	result.ProcessingTime = time.Since(startTime)
	return result
}

// collectResults collects processing results and updates statistics
func (pp *ParallelProcessor) collectResults(stats *ProcessingStats) {
	for result := range pp.resultChan {
		if result.Error != nil {
			stats.FailedFiles++
		} else {
			stats.ProcessedFiles++
			stats.TotalSize += result.Size
		}
	}
}

// ExecuteBatchOperations executes all batched file operations
func (pp *ParallelProcessor) ExecuteBatchOperations() error {
	return pp.batchOps.Execute()
}

// collectTemplateFiles collects all files in a template directory
func (pp *ParallelProcessor) collectTemplateFiles(templateDir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// TemplatePreloader preloads and caches frequently used templates
type TemplatePreloader struct {
	engine        *Engine
	preloadedDirs map[string]bool
	mutex         sync.RWMutex
}

// NewTemplatePreloader creates a new template preloader
func NewTemplatePreloader(engine *Engine) *TemplatePreloader {
	return &TemplatePreloader{
		engine:        engine,
		preloadedDirs: make(map[string]bool),
	}
}

// PreloadDirectory preloads all templates in a directory
func (tp *TemplatePreloader) PreloadDirectory(templateDir string) error {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	// Check if already preloaded
	if tp.preloadedDirs[templateDir] {
		return nil
	}

	// Collect all template files
	var templateFiles []string
	err := filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".tmpl") {
			templateFiles = append(templateFiles, path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk template directory: %w", err)
	}

	// Preload templates in parallel
	var wg sync.WaitGroup
	errorChan := make(chan error, len(templateFiles))

	for _, templateFile := range templateFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			_, err := tp.engine.LoadTemplate(file)
			if err != nil {
				errorChan <- fmt.Errorf("failed to preload template %s: %w", file, err)
			}
		}(templateFile)
	}

	wg.Wait()
	close(errorChan)

	// Check for errors
	for err := range errorChan {
		return err
	}

	tp.preloadedDirs[templateDir] = true
	return nil
}

// IsPreloaded checks if a directory has been preloaded
func (tp *TemplatePreloader) IsPreloaded(templateDir string) bool {
	tp.mutex.RLock()
	defer tp.mutex.RUnlock()
	return tp.preloadedDirs[templateDir]
}

// ClearPreloaded clears the preloaded directory cache
func (tp *TemplatePreloader) ClearPreloaded() {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	tp.preloadedDirs = make(map[string]bool)
}

// OptimizedTemplateEngine wraps the regular engine with performance optimizations
type OptimizedTemplateEngine struct {
	*Engine
	parallelProcessor *ParallelProcessor
	preloader         *TemplatePreloader
	enableParallel    bool
	enablePreload     bool
}

// NewOptimizedTemplateEngine creates a new optimized template engine
func NewOptimizedTemplateEngine(engine *Engine) *OptimizedTemplateEngine {
	return &OptimizedTemplateEngine{
		Engine:            engine,
		parallelProcessor: NewParallelProcessor(engine, runtime.NumCPU()),
		preloader:         NewTemplatePreloader(engine),
		enableParallel:    true,
		enablePreload:     true,
	}
}

// ProcessDirectory processes a directory with optimizations
func (ote *OptimizedTemplateEngine) ProcessDirectory(templateDir string, outputDir string, config *models.ProjectConfig) error {
	// Preload templates if enabled
	if ote.enablePreload && !ote.preloader.IsPreloaded(templateDir) {
		if err := ote.preloader.PreloadDirectory(templateDir); err != nil {
			// Log warning but don't fail - preloading is an optimization
			fmt.Printf("Warning: Failed to preload templates: %v\n", err)
		}
	}

	// Use parallel processing if enabled and there are enough files
	if ote.enableParallel {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		stats, err := ote.parallelProcessor.ProcessDirectoryParallel(ctx, templateDir, outputDir, config)
		if err != nil {
			// Fall back to sequential processing
			fmt.Printf("Warning: Parallel processing failed, falling back to sequential: %v\n", err)
			return ote.Engine.ProcessDirectory(templateDir, outputDir, config)
		}

		// Execute batch operations
		if err := ote.parallelProcessor.ExecuteBatchOperations(); err != nil {
			return fmt.Errorf("failed to execute batch operations: %w", err)
		}

		fmt.Printf("Processed %d files in %v (%.2f files/sec)\n",
			stats.ProcessedFiles, stats.ProcessingTime, stats.Throughput)

		return nil
	}

	// Fall back to sequential processing
	return ote.Engine.ProcessDirectory(templateDir, outputDir, config)
}

// SetParallelProcessing enables or disables parallel processing
func (ote *OptimizedTemplateEngine) SetParallelProcessing(enabled bool) {
	ote.enableParallel = enabled
}

// SetPreloading enables or disables template preloading
func (ote *OptimizedTemplateEngine) SetPreloading(enabled bool) {
	ote.enablePreload = enabled
}

// GetProcessingStats returns processing statistics
func (ote *OptimizedTemplateEngine) GetProcessingStats() interface{} {
	// Return cache statistics
	return struct {
		TemplateCache interface{}
		RenderCache   interface{}
	}{
		TemplateCache: ote.Engine.templateCache.Stats(),
		RenderCache:   "render cache stats", // Placeholder
	}
}
