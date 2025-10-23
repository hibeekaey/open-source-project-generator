package orchestrator

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/internal/generator/bootstrap"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
)

func TestExecutorRegistry_Register(t *testing.T) {
	log := logger.NewLogger()
	registry := NewExecutorRegistry(log)

	// Register NextJS executor
	executor := bootstrap.NewExecutorAdapter(bootstrap.NewNextJSExecutor())
	registry.Register("nextjs", executor)

	// Verify it was registered
	if !registry.Has("nextjs") {
		t.Error("Expected registry to have 'nextjs' executor")
	}

	// Verify we can retrieve it
	retrieved, err := registry.Get("nextjs")
	if err != nil {
		t.Errorf("Expected to retrieve 'nextjs' executor, got error: %v", err)
	}
	if retrieved == nil {
		t.Error("Expected non-nil executor")
	}
}

func TestExecutorRegistry_Get_NotFound(t *testing.T) {
	log := logger.NewLogger()
	registry := NewExecutorRegistry(log)

	// Try to get non-existent executor
	_, err := registry.Get("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent executor")
	}
}

func TestExecutorRegistry_Has(t *testing.T) {
	log := logger.NewLogger()
	registry := NewExecutorRegistry(log)

	// Register NextJS executor
	executor := bootstrap.NewExecutorAdapter(bootstrap.NewNextJSExecutor())
	registry.Register("nextjs", executor)

	tests := []struct {
		name          string
		componentType string
		want          bool
	}{
		{
			name:          "registered component",
			componentType: "nextjs",
			want:          true,
		},
		{
			name:          "unregistered component",
			componentType: "android",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := registry.Has(tt.componentType)
			if got != tt.want {
				t.Errorf("Has(%s) = %v, want %v", tt.componentType, got, tt.want)
			}
		})
	}
}

func TestExecutorRegistry_GetDefaultFlags(t *testing.T) {
	log := logger.NewLogger()
	registry := NewExecutorRegistry(log)

	// Register executors
	registry.Register("nextjs", bootstrap.NewExecutorAdapter(bootstrap.NewNextJSExecutor()))
	registry.Register("go-backend", bootstrap.NewExecutorAdapter(bootstrap.NewGoExecutor()))

	tests := []struct {
		name          string
		componentType string
		wantNonEmpty  bool
	}{
		{
			name:          "nextjs returns flags",
			componentType: "nextjs",
			wantNonEmpty:  true,
		},
		{
			name:          "go-backend returns flags",
			componentType: "go-backend",
			wantNonEmpty:  true,
		},
		{
			name:          "unregistered returns empty",
			componentType: "android",
			wantNonEmpty:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := registry.GetDefaultFlags(tt.componentType)

			if tt.wantNonEmpty && len(flags) == 0 {
				t.Errorf("GetDefaultFlags(%s) returned empty, want non-empty", tt.componentType)
			}

			if !tt.wantNonEmpty && len(flags) != 0 {
				t.Errorf("GetDefaultFlags(%s) returned %d flags, want empty", tt.componentType, len(flags))
			}
		})
	}
}

func TestExecutorRegistry_GetSupportedTypes(t *testing.T) {
	log := logger.NewLogger()
	registry := NewExecutorRegistry(log)

	// Register multiple executors
	registry.Register("nextjs", bootstrap.NewExecutorAdapter(bootstrap.NewNextJSExecutor()))
	registry.Register("go-backend", bootstrap.NewExecutorAdapter(bootstrap.NewGoExecutor()))
	registry.Register("android", bootstrap.NewExecutorAdapter(bootstrap.NewAndroidExecutor(nil)))

	types := registry.GetSupportedTypes()

	// Should have 3 types
	if len(types) != 3 {
		t.Errorf("GetSupportedTypes() returned %d types, want 3", len(types))
	}

	// Check that all expected types are present
	expectedTypes := map[string]bool{
		"nextjs":     false,
		"go-backend": false,
		"android":    false,
	}

	for _, typ := range types {
		if _, exists := expectedTypes[typ]; exists {
			expectedTypes[typ] = true
		}
	}

	for typ, found := range expectedTypes {
		if !found {
			t.Errorf("Expected type %s not found in supported types", typ)
		}
	}
}

func TestExecutorRegistry_MultipleRegistrations(t *testing.T) {
	log := logger.NewLogger()
	registry := NewExecutorRegistry(log)

	// Register all executors
	registry.Register("nextjs", bootstrap.NewExecutorAdapter(bootstrap.NewNextJSExecutor()))
	registry.Register("go-backend", bootstrap.NewExecutorAdapter(bootstrap.NewGoExecutor()))
	registry.Register("android", bootstrap.NewExecutorAdapter(bootstrap.NewAndroidExecutor(nil)))
	registry.Register("ios", bootstrap.NewExecutorAdapter(bootstrap.NewiOSExecutor(nil)))

	// Verify all are registered
	componentTypes := []string{"nextjs", "go-backend", "android", "ios"}
	for _, typ := range componentTypes {
		if !registry.Has(typ) {
			t.Errorf("Expected registry to have '%s' executor", typ)
		}

		executor, err := registry.Get(typ)
		if err != nil {
			t.Errorf("Expected to retrieve '%s' executor, got error: %v", typ, err)
		}
		if executor == nil {
			t.Errorf("Expected non-nil executor for '%s'", typ)
		}

		// Verify SupportsComponent works
		if !executor.SupportsComponent(typ) {
			t.Errorf("Expected executor to support component type '%s'", typ)
		}
	}
}
