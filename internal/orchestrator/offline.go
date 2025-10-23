package orchestrator

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
)

// OfflineDetector manages offline mode detection and status
type OfflineDetector struct {
	isOffline     bool
	lastCheck     time.Time
	checkInterval time.Duration
	mu            sync.RWMutex
	logger        *logger.Logger
	testURLs      []string
	timeout       time.Duration
	forceOffline  bool
}

// OfflineDetectorConfig holds configuration for offline detection
type OfflineDetectorConfig struct {
	CheckInterval time.Duration // How often to check connectivity
	Timeout       time.Duration // Timeout for connectivity checks
	TestURLs      []string      // URLs to test for connectivity
	ForceOffline  bool          // Force offline mode regardless of connectivity
}

// DefaultOfflineDetectorConfig returns default offline detector configuration
func DefaultOfflineDetectorConfig() *OfflineDetectorConfig {
	return &OfflineDetectorConfig{
		CheckInterval: 30 * time.Second,
		Timeout:       5 * time.Second,
		TestURLs: []string{
			"https://www.google.com",
			"https://www.cloudflare.com",
			"https://1.1.1.1",
		},
		ForceOffline: false,
	}
}

// NewOfflineDetector creates a new offline detector
func NewOfflineDetector(config *OfflineDetectorConfig, log *logger.Logger) *OfflineDetector {
	if config == nil {
		config = DefaultOfflineDetectorConfig()
	}

	od := &OfflineDetector{
		isOffline:     false,
		lastCheck:     time.Time{},
		checkInterval: config.CheckInterval,
		logger:        log,
		testURLs:      config.TestURLs,
		timeout:       config.Timeout,
		forceOffline:  config.ForceOffline,
	}

	// Perform initial check
	od.Check()

	return od
}

// IsOffline returns whether the system is currently offline
func (od *OfflineDetector) IsOffline() bool {
	od.mu.RLock()
	defer od.mu.RUnlock()

	// If forced offline, always return true
	if od.forceOffline {
		return true
	}

	// Check if we need to refresh the status
	if time.Since(od.lastCheck) > od.checkInterval {
		// Release read lock and acquire write lock
		od.mu.RUnlock()
		od.Check()
		od.mu.RLock()
	}

	return od.isOffline
}

// Check performs a connectivity check
func (od *OfflineDetector) Check() bool {
	od.mu.Lock()
	defer od.mu.Unlock()

	// If forced offline, don't check
	if od.forceOffline {
		od.isOffline = true
		od.lastCheck = time.Now()
		return true
	}

	// Try multiple methods to detect connectivity
	online := od.checkHTTP() || od.checkDNS()

	od.isOffline = !online
	od.lastCheck = time.Now()

	if od.logger != nil {
		if od.isOffline {
			od.logger.Info("Offline mode detected - will use fallback generators and cached tools")
		} else {
			od.logger.Debug("Online connectivity confirmed")
		}
	}

	return od.isOffline
}

// checkHTTP attempts to connect to test URLs via HTTP
func (od *OfflineDetector) checkHTTP() bool {
	client := &http.Client{
		Timeout: od.timeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	for _, url := range od.testURLs {
		ctx, cancel := context.WithTimeout(context.Background(), od.timeout)
		req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
		if err != nil {
			cancel()
			continue
		}

		resp, err := client.Do(req)
		cancel()

		if err == nil {
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 500 {
				return true // Online
			}
		}
	}

	return false // All HTTP checks failed
}

// checkDNS attempts to resolve a known domain
func (od *OfflineDetector) checkDNS() bool {
	resolver := &net.Resolver{
		PreferGo: true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), od.timeout)
	defer cancel()

	// Try to resolve a well-known domain
	_, err := resolver.LookupHost(ctx, "www.google.com")
	return err == nil
}

// ForceOffline forces the detector to report offline status
func (od *OfflineDetector) ForceOffline(force bool) {
	od.mu.Lock()
	defer od.mu.Unlock()

	od.forceOffline = force
	if force {
		od.isOffline = true
		if od.logger != nil {
			od.logger.Info("Offline mode forced - will use fallback generators")
		}
	} else {
		// Trigger a fresh check
		od.mu.Unlock()
		od.Check()
		od.mu.Lock()
	}
}

// GetStatus returns detailed offline status information
func (od *OfflineDetector) GetStatus() map[string]interface{} {
	od.mu.RLock()
	defer od.mu.RUnlock()

	return map[string]interface{}{
		"offline":          od.isOffline,
		"forced":           od.forceOffline,
		"last_check":       od.lastCheck,
		"check_interval":   od.checkInterval.String(),
		"time_since_check": time.Since(od.lastCheck).String(),
	}
}

// SetCheckInterval updates the check interval
func (od *OfflineDetector) SetCheckInterval(interval time.Duration) {
	od.mu.Lock()
	defer od.mu.Unlock()

	od.checkInterval = interval
}

// GetOfflineMessage returns a user-friendly message about offline status
func (od *OfflineDetector) GetOfflineMessage() string {
	if !od.IsOffline() {
		return ""
	}

	if od.forceOffline {
		return "⚠️  Offline mode (forced) - Using fallback generators and cached tools"
	}

	return "⚠️  No internet connection detected - Using fallback generators and cached tools"
}

// OfflineMode adds offline detection to ToolDiscovery
func (td *ToolDiscovery) SetOfflineDetector(detector *OfflineDetector) {
	td.isOffline = detector.IsOffline()
}

// IsOfflineMode returns whether tool discovery is in offline mode
func (td *ToolDiscovery) IsOfflineMode() bool {
	return td.isOffline
}

// SetOfflineMode manually sets offline mode
func (td *ToolDiscovery) SetOfflineMode(offline bool) {
	td.isOffline = offline
	if td.logger != nil {
		if offline {
			td.logger.Info("Tool discovery set to offline mode")
		} else {
			td.logger.Info("Tool discovery set to online mode")
		}
	}
}

// ShouldUseFallback determines if fallback generation should be used
func (td *ToolDiscovery) ShouldUseFallback(componentType string) (bool, string) {
	// Check if we're in offline mode
	if td.isOffline {
		return true, "offline mode active"
	}

	// Check if required tools are available
	tools := td.GetToolsForComponent(componentType)
	if len(tools) == 0 {
		return true, "no tools registered for component type"
	}

	// Check each required tool
	for _, toolName := range tools {
		available, err := td.IsAvailable(toolName)
		if err != nil || !available {
			// Check if fallback is available
			if td.HasFallback(componentType) {
				return true, fmt.Sprintf("tool '%s' not available", toolName)
			}
			return false, fmt.Sprintf("tool '%s' required but not available", toolName)
		}
	}

	return false, "all required tools available"
}
