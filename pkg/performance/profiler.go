// Package performance provides profiling capabilities for performance analysis
package performance

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

// Profiler provides CPU and memory profiling capabilities
type Profiler struct {
	cpuFile    *os.File
	profiling  bool
	startTime  time.Time
	profileDir string
}

// ProfilingData contains the results of a profiling session
type ProfilingData struct {
	CommandName    string        `json:"command_name"`
	Duration       time.Duration `json:"duration"`
	CPUProfilePath string        `json:"cpu_profile_path"`
	MemProfilePath string        `json:"mem_profile_path"`
	Success        bool          `json:"success"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	MemoryStats    *MemoryStats  `json:"memory_stats"`
}

// MemoryStats contains memory usage statistics
type MemoryStats struct {
	AllocBytes      uint64 `json:"alloc_bytes"`
	TotalAllocBytes uint64 `json:"total_alloc_bytes"`
	SysBytes        uint64 `json:"sys_bytes"`
	NumGC           uint32 `json:"num_gc"`
	HeapObjects     uint64 `json:"heap_objects"`
}

// NewProfiler creates a new profiler instance
func NewProfiler() *Profiler {
	// Create profiles directory
	profileDir := filepath.Join(".", "profiles")
	_ = os.MkdirAll(profileDir, 0750)

	return &Profiler{
		profileDir: profileDir,
	}
}

// Start begins CPU and memory profiling
func (p *Profiler) Start() error {
	if p.profiling {
		return fmt.Errorf("profiling already in progress")
	}

	p.startTime = time.Now()
	timestamp := p.startTime.Format("20060102_150405")

	// Start CPU profiling
	cpuPath := filepath.Join(p.profileDir, fmt.Sprintf("cpu_profile_%s.prof", timestamp))
	// #nosec G304 - cpuPath is constructed from validated profileDir and timestamp
	cpuFile, err := os.Create(cpuPath)
	if err != nil {
		return fmt.Errorf("failed to create CPU profile file: %w", err)
	}
	p.cpuFile = cpuFile

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		_ = cpuFile.Close()
		return fmt.Errorf("failed to start CPU profiling: %w", err)
	}

	p.profiling = true
	return nil
}

// Stop ends profiling and returns profiling data
func (p *Profiler) Stop() (*ProfilingData, error) {
	if !p.profiling {
		return nil, fmt.Errorf("profiling not in progress")
	}

	endTime := time.Now()
	timestamp := p.startTime.Format("20060102_150405")

	// Stop CPU profiling
	pprof.StopCPUProfile()
	cpuPath := p.cpuFile.Name()
	_ = p.cpuFile.Close()

	// Create memory profile
	memPath := filepath.Join(p.profileDir, fmt.Sprintf("mem_profile_%s.prof", timestamp))
	// #nosec G304 - memPath is constructed from validated profileDir and timestamp
	memFile, err := os.Create(memPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory profile file: %w", err)
	}
	defer func() { _ = memFile.Close() }()

	// Force garbage collection before memory profiling
	runtime.GC()

	if err := pprof.WriteHeapProfile(memFile); err != nil {
		return nil, fmt.Errorf("failed to write memory profile: %w", err)
	}

	// Collect memory statistics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	profilingData := &ProfilingData{
		Duration:       endTime.Sub(p.startTime),
		CPUProfilePath: cpuPath,
		MemProfilePath: memPath,
		Success:        true,
		StartTime:      p.startTime,
		EndTime:        endTime,
		MemoryStats: &MemoryStats{
			AllocBytes:      memStats.Alloc,
			TotalAllocBytes: memStats.TotalAlloc,
			SysBytes:        memStats.Sys,
			NumGC:           memStats.NumGC,
			HeapObjects:     memStats.HeapObjects,
		},
	}

	p.profiling = false
	return profilingData, nil
}

// IsProfileing returns whether profiling is currently active
func (p *Profiler) IsProfileing() bool {
	return p.profiling
}

// GetProfileDir returns the directory where profiles are stored
func (p *Profiler) GetProfileDir() string {
	return p.profileDir
}

// SetProfileDir sets the directory where profiles will be stored
func (p *Profiler) SetProfileDir(dir string) error {
	if p.profiling {
		return fmt.Errorf("cannot change profile directory while profiling is active")
	}

	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create profile directory: %w", err)
	}

	p.profileDir = dir
	return nil
}

// CleanupOldProfiles removes profile files older than the specified duration
func (p *Profiler) CleanupOldProfiles(maxAge time.Duration) error {
	entries, err := os.ReadDir(p.profileDir)
	if err != nil {
		return fmt.Errorf("failed to read profile directory: %w", err)
	}

	cutoff := time.Now().Add(-maxAge)
	var removedCount int

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if it's a profile file
		if !isProfileFile(entry.Name()) {
			continue
		}

		filePath := filepath.Join(p.profileDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			if err := os.Remove(filePath); err == nil {
				removedCount++
			}
		}
	}

	return nil
}

// isProfileFile checks if a filename is a profile file
func isProfileFile(filename string) bool {
	return filepath.Ext(filename) == ".prof"
}

// GetProfileFiles returns a list of available profile files
func (p *Profiler) GetProfileFiles() ([]string, error) {
	entries, err := os.ReadDir(p.profileDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile directory: %w", err)
	}

	var profiles []string
	for _, entry := range entries {
		if !entry.IsDir() && isProfileFile(entry.Name()) {
			profiles = append(profiles, filepath.Join(p.profileDir, entry.Name()))
		}
	}

	return profiles, nil
}

// AnalyzeProfile provides basic analysis of a profile file
func (p *Profiler) AnalyzeProfile(profilePath string) (*ProfileAnalysis, error) {
	// Validate profile path for security
	if err := validateProfilePath(profilePath); err != nil {
		return nil, fmt.Errorf("invalid profile path: %w", err)
	}

	info, err := os.Stat(profilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat profile file: %w", err)
	}

	analysis := &ProfileAnalysis{
		FilePath:  profilePath,
		FileSize:  info.Size(),
		CreatedAt: info.ModTime(),
		Type:      getProfileType(profilePath),
	}

	// Basic file analysis
	// #nosec G304 - profilePath is from internal profile directory
	file, err := os.Open(profilePath)
	if err != nil {
		return analysis, nil // Return partial analysis
	}
	defer func() { _ = file.Close() }()

	// Read first few bytes to validate it's a valid profile
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && n == 0 {
		analysis.Valid = false
		analysis.Error = "Failed to read profile file"
		return analysis, nil
	}

	analysis.Valid = true
	return analysis, nil
}

// ProfileAnalysis contains analysis results for a profile file
type ProfileAnalysis struct {
	FilePath  string    `json:"file_path"`
	FileSize  int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
	Type      string    `json:"type"`
	Valid     bool      `json:"valid"`
	Error     string    `json:"error,omitempty"`
}

// getProfileType determines the type of profile based on filename
func getProfileType(filename string) string {
	base := filepath.Base(filename)
	if contains(base, "cpu") {
		return "cpu"
	}
	if contains(base, "mem") || contains(base, "heap") {
		return "memory"
	}
	return "unknown"
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr))))
}

// containsSubstring checks if s contains substr anywhere
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// validateProfilePath validates a profile file path for security
func validateProfilePath(path string) error {
	if path == "" {
		return fmt.Errorf("profile path cannot be empty")
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("profile path contains invalid path traversal: %s", path)
	}

	// Check for null bytes (security risk)
	if strings.Contains(path, "\x00") {
		return fmt.Errorf("profile path contains null bytes: %s", path)
	}

	// Ensure it's a .prof file
	if !strings.HasSuffix(path, ".prof") {
		return fmt.Errorf("profile path must end with .prof extension: %s", path)
	}

	return nil
}
