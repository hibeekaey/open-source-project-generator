package versions

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config represents the centralized version configuration
type Config struct {
	Frontend       FrontendVersions       `yaml:"frontend"`
	Backend        BackendVersions        `yaml:"backend"`
	Android        AndroidVersions        `yaml:"android"`
	IOS            IOSVersions            `yaml:"ios"`
	Docker         DockerVersions         `yaml:"docker"`
	Infrastructure InfrastructureVersions `yaml:"infrastructure"`
	Metadata       Metadata               `yaml:"metadata"`
}

type FrontendVersions struct {
	NextJS struct {
		Version string `yaml:"version"`
		Package string `yaml:"package"`
	} `yaml:"nextjs"`
	React struct {
		Version string `yaml:"version"`
	} `yaml:"react"`
	ReactDOM struct {
		Version string `yaml:"version"`
	} `yaml:"react_dom"`
	TypeScript struct {
		Version string `yaml:"version"`
	} `yaml:"typescript"`
}

type BackendVersions struct {
	Go struct {
		Version   string `yaml:"version"`
		DockerTag string `yaml:"docker_tag"`
	} `yaml:"go"`
	Frameworks struct {
		Gin struct {
			Version string `yaml:"version"`
			Package string `yaml:"package"`
		} `yaml:"gin"`
		GinCORS struct {
			Version string `yaml:"version"`
			Package string `yaml:"package"`
		} `yaml:"gin_cors"`
		Echo struct {
			Version string `yaml:"version"`
			Package string `yaml:"package"`
		} `yaml:"echo"`
		Fiber struct {
			Version string `yaml:"version"`
			Package string `yaml:"package"`
		} `yaml:"fiber"`
	} `yaml:"frameworks"`
}

type AndroidVersions struct {
	CompileSDK  string `yaml:"compile_sdk"`
	MinSDK      string `yaml:"min_sdk"`
	TargetSDK   string `yaml:"target_sdk"`
	JavaVersion string `yaml:"java_version"`
	Kotlin      struct {
		Version string `yaml:"version"`
	} `yaml:"kotlin"`
	Gradle struct {
		Version         string `yaml:"version"`
		DistributionURL string `yaml:"distribution_url"`
	} `yaml:"gradle"`
	GradlePlugin struct {
		Version string `yaml:"version"`
	} `yaml:"gradle_plugin"`
	AndroidX struct {
		CoreKTX struct {
			Version string `yaml:"version"`
			Package string `yaml:"package"`
		} `yaml:"core_ktx"`
		AppCompat struct {
			Version string `yaml:"version"`
			Package string `yaml:"package"`
		} `yaml:"appcompat"`
		Material struct {
			Version string `yaml:"version"`
			Package string `yaml:"package"`
		} `yaml:"material"`
		ConstraintLayout struct {
			Version string `yaml:"version"`
			Package string `yaml:"package"`
		} `yaml:"constraintlayout"`
	} `yaml:"androidx"`
	Testing struct {
		JUnit struct {
			Version string `yaml:"version"`
			Package string `yaml:"package"`
		} `yaml:"junit"`
		AndroidXJUnit struct {
			Version string `yaml:"version"`
			Package string `yaml:"package"`
		} `yaml:"androidx_junit"`
		Espresso struct {
			Version string `yaml:"version"`
			Package string `yaml:"package"`
		} `yaml:"espresso"`
	} `yaml:"testing"`
}

type IOSVersions struct {
	Swift struct {
		Version      string `yaml:"version"`
		ShortVersion string `yaml:"short_version"`
	} `yaml:"swift"`
	Xcode struct {
		Version string `yaml:"version"`
	} `yaml:"xcode"`
	DeploymentTarget string `yaml:"deployment_target"`
}

type DockerVersions struct {
	Alpine struct {
		Version string `yaml:"version"`
	} `yaml:"alpine"`
	Golang struct {
		Version string `yaml:"version"`
	} `yaml:"golang"`
	Ubuntu struct {
		Version string `yaml:"version"`
	} `yaml:"ubuntu"`
}

type InfrastructureVersions struct {
	Terraform struct {
		Version string `yaml:"version"`
	} `yaml:"terraform"`
	Kubernetes struct {
		Version string `yaml:"version"`
	} `yaml:"kubernetes"`
}

type Metadata struct {
	LastUpdated   string `yaml:"last_updated"`
	SchemaVersion string `yaml:"schema_version"`
}

var (
	globalConfig *Config
	configMutex  sync.RWMutex
	configPath   = "configs/versions.yaml"
)

// Load reads the versions configuration from the YAML file
func Load() (*Config, error) {
	// Try to find the config file if the default path doesn't exist
	path := configPath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		foundPath, findErr := FindConfigPath()
		if findErr == nil {
			path = foundPath
		}
	}
	return LoadFrom(path)
}

// LoadFrom reads the versions configuration from a specific path
func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read versions config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse versions config: %w", err)
	}

	return &config, nil
}

// Get returns the global version configuration, loading it if necessary
func Get() (*Config, error) {
	configMutex.RLock()
	if globalConfig != nil {
		defer configMutex.RUnlock()
		return globalConfig, nil
	}
	configMutex.RUnlock()

	configMutex.Lock()
	defer configMutex.Unlock()

	// Double-check after acquiring write lock
	if globalConfig != nil {
		return globalConfig, nil
	}

	config, err := Load()
	if err != nil {
		return nil, err
	}

	globalConfig = config
	return globalConfig, nil
}

// Reload forces a reload of the version configuration
func Reload() error {
	configMutex.Lock()
	defer configMutex.Unlock()

	config, err := Load()
	if err != nil {
		return err
	}

	globalConfig = config
	return nil
}

// SetConfigPath sets a custom path for the versions configuration file
func SetConfigPath(path string) {
	configMutex.Lock()
	defer configMutex.Unlock()
	configPath = path
	globalConfig = nil // Force reload on next Get()
}

// FindConfigPath searches for the versions.yaml file starting from the current directory
func FindConfigPath() (string, error) {
	// Try current directory first
	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	}

	// Try relative to executable
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		candidatePath := filepath.Join(execDir, configPath)
		if _, err := os.Stat(candidatePath); err == nil {
			return candidatePath, nil
		}
	}

	// Try walking up the directory tree
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	for {
		candidatePath := filepath.Join(currentDir, configPath)
		if _, err := os.Stat(candidatePath); err == nil {
			return candidatePath, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}

	return "", fmt.Errorf("versions.yaml not found")
}
