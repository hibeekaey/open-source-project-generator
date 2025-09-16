package models

import (
	"time"
)

// VersionInfo represents basic version information for a package or language
type VersionInfo struct {
	// Package identification
	Name     string `yaml:"name" json:"name"`
	Language string `yaml:"language" json:"language"`
	Type     string `yaml:"type" json:"type"`

	// Version tracking
	LatestVersion string `yaml:"latest_version" json:"latest_version"`

	// Update metadata
	UpdatedAt    time.Time `yaml:"updated_at" json:"updated_at"`
	CheckedAt    time.Time `yaml:"checked_at" json:"checked_at"`
	UpdateSource string    `yaml:"update_source" json:"update_source"`

	// Registry information
	RegistryURL string            `yaml:"registry_url,omitempty" json:"registry_url,omitempty"`
	Metadata    map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// VersionConfig represents the version configuration for a project
type VersionConfig struct {
	// Language versions
	Node string `yaml:"node" json:"node"`
	Go   string `yaml:"go" json:"go"`

	// Package versions
	Packages map[string]string `yaml:"packages" json:"packages"`
}
