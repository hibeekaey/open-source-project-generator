package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v3"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// FileStorage implements VersionStorage interface using file system
type FileStorage struct {
	filePath   string
	backupDir  string
	format     string // "yaml" or "json"
	store      *models.VersionStore
	lastLoaded time.Time
	mu         sync.RWMutex // Mutex for concurrent access
}

// NewFileStorage creates a new file-based version storage
func NewFileStorage(filePath string, format string) (*FileStorage, error) {
	if format != "yaml" && format != "json" {
		return nil, fmt.Errorf("unsupported format: %s, must be 'yaml' or 'json'", format)
	}

	backupDir := filepath.Join(filepath.Dir(filePath), "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	storage := &FileStorage{
		filePath:  filePath,
		backupDir: backupDir,
		format:    format,
	}

	// Initialize with empty store if file doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		storage.store = &models.VersionStore{
			LastUpdated: time.Now(),
			Version:     "1.0.0",
			Languages:   make(map[string]*models.VersionInfo),
			Frameworks:  make(map[string]*models.VersionInfo),
			Packages:    make(map[string]*models.VersionInfo),
			UpdatePolicy: models.UpdatePolicy{
				AutoUpdate:             true,
				SecurityPriority:       true,
				BreakingChangeApproval: true,
				UpdateSchedule:         "daily",
				MaxAge:                 24 * time.Hour,
			},
		}
		if err := storage.Save(storage.store); err != nil {
			return nil, fmt.Errorf("failed to initialize storage file: %w", err)
		}
	}

	return storage, nil
}

// Load loads the complete version store from storage
func (fs *FileStorage) Load() (*models.VersionStore, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read version store file: %w", err)
	}

	var store models.VersionStore
	switch fs.format {
	case "yaml":
		if err := yaml.Unmarshal(data, &store); err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
		}
	case "json":
		if err := json.Unmarshal(data, &store); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	// Initialize maps if they're nil
	if store.Languages == nil {
		store.Languages = make(map[string]*models.VersionInfo)
	}
	if store.Frameworks == nil {
		store.Frameworks = make(map[string]*models.VersionInfo)
	}
	if store.Packages == nil {
		store.Packages = make(map[string]*models.VersionInfo)
	}

	fs.store = &store
	fs.lastLoaded = time.Now()
	return &store, nil
}

// Save saves the complete version store to storage
func (fs *FileStorage) Save(store *models.VersionStore) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return fs.saveUnlocked(store)
}

// saveUnlocked saves the store without acquiring locks (internal use)
func (fs *FileStorage) saveUnlocked(store *models.VersionStore) error {
	// Update timestamp
	store.LastUpdated = time.Now()

	var data []byte
	var err error

	switch fs.format {
	case "yaml":
		data, err = yaml.Marshal(store)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}
	case "json":
		data, err = json.MarshalIndent(store, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	}

	// Write file with atomic operation using temporary file
	tempFile := fs.filePath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	if err := os.Rename(tempFile, fs.filePath); err != nil {
		os.Remove(tempFile) // Clean up on failure
		return fmt.Errorf("failed to atomically replace file: %w", err)
	}

	fs.store = store
	return nil
}

// GetVersionInfo retrieves version information for a specific package
func (fs *FileStorage) GetVersionInfo(name string) (*models.VersionInfo, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	if fs.store == nil {
		// Temporarily unlock to load, then relock
		fs.mu.RUnlock()
		if _, err := fs.Load(); err != nil {
			fs.mu.RLock()
			return nil, err
		}
		fs.mu.RLock()
	}

	// Search in all categories
	if info, exists := fs.store.Languages[name]; exists {
		return info, nil
	}
	if info, exists := fs.store.Frameworks[name]; exists {
		return info, nil
	}
	if info, exists := fs.store.Packages[name]; exists {
		return info, nil
	}

	return nil, fmt.Errorf("version info not found for: %s", name)
}

// SetVersionInfo stores version information for a specific package
func (fs *FileStorage) SetVersionInfo(name string, info *models.VersionInfo) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if fs.store == nil {
		// Temporarily unlock to load, then relock
		fs.mu.Unlock()
		if _, err := fs.Load(); err != nil {
			fs.mu.Lock()
			return err
		}
		fs.mu.Lock()
	}

	// Determine which category to store in based on type
	switch strings.ToLower(info.Type) {
	case "language":
		fs.store.Languages[name] = info
	case "framework":
		fs.store.Frameworks[name] = info
	case "package":
		fs.store.Packages[name] = info
	default:
		return fmt.Errorf("unknown version info type: %s", info.Type)
	}

	// Call save without additional locking since we already hold the lock
	return fs.saveUnlocked(fs.store)
}

// DeleteVersionInfo removes version information for a specific package
func (fs *FileStorage) DeleteVersionInfo(name string) error {
	if fs.store == nil {
		if _, err := fs.Load(); err != nil {
			return err
		}
	}

	found := false
	if _, exists := fs.store.Languages[name]; exists {
		delete(fs.store.Languages, name)
		found = true
	}
	if _, exists := fs.store.Frameworks[name]; exists {
		delete(fs.store.Frameworks, name)
		found = true
	}
	if _, exists := fs.store.Packages[name]; exists {
		delete(fs.store.Packages, name)
		found = true
	}

	if !found {
		return fmt.Errorf("version info not found for: %s", name)
	}

	return fs.Save(fs.store)
}

// ListVersions returns all stored version information
func (fs *FileStorage) ListVersions() (map[string]*models.VersionInfo, error) {
	if fs.store == nil {
		if _, err := fs.Load(); err != nil {
			return nil, err
		}
	}

	result := make(map[string]*models.VersionInfo)

	// Combine all categories
	for name, info := range fs.store.Languages {
		result[name] = info
	}
	for name, info := range fs.store.Frameworks {
		result[name] = info
	}
	for name, info := range fs.store.Packages {
		result[name] = info
	}

	return result, nil
}

// Query searches for version information based on criteria
func (fs *FileStorage) Query(query *models.VersionQuery) (map[string]*models.VersionInfo, error) {
	allVersions, err := fs.ListVersions()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*models.VersionInfo)

	for name, info := range allVersions {
		// Apply filters
		if query.Name != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(query.Name)) {
			continue
		}
		if query.Language != "" && !strings.EqualFold(info.Language, query.Language) {
			continue
		}
		if query.Type != "" && !strings.EqualFold(info.Type, query.Type) {
			continue
		}
		if query.Outdated && info.CurrentVersion == info.LatestVersion {
			continue
		}
		if query.Insecure && info.IsSecure {
			continue
		}

		result[name] = info
	}

	return result, nil
}

// Backup creates a backup of the current version store
func (fs *FileStorage) Backup() error {
	if _, err := os.Stat(fs.filePath); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", fs.filePath)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("versions_backup_%s.%s", timestamp, fs.format)
	backupPath := filepath.Join(fs.backupDir, backupFileName)

	sourceData, err := os.ReadFile(fs.filePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	if err := os.WriteFile(backupPath, sourceData, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// Restore restores the version store from a backup
func (fs *FileStorage) Restore(backupPath string) error {
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupPath)
	}

	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	if err := os.WriteFile(fs.filePath, backupData, 0644); err != nil {
		return fmt.Errorf("failed to restore file: %w", err)
	}

	// Reload the store
	_, err = fs.Load()
	return err
}

// GetLastUpdated returns the timestamp of the last update
func (fs *FileStorage) GetLastUpdated() (time.Time, error) {
	if fs.store == nil {
		if _, err := fs.Load(); err != nil {
			return time.Time{}, err
		}
	}
	return fs.store.LastUpdated, nil
}

// SetLastUpdated updates the last updated timestamp
func (fs *FileStorage) SetLastUpdated(timestamp time.Time) error {
	if fs.store == nil {
		if _, err := fs.Load(); err != nil {
			return err
		}
	}
	fs.store.LastUpdated = timestamp
	return fs.Save(fs.store)
}

// Ensure FileStorage implements VersionStorage interface
var _ interfaces.VersionStorage = (*FileStorage)(nil)
