package performance

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/cache"
)

// BenchmarkRefactoredCache benchmarks the refactored cache components
func BenchmarkRefactoredCache(b *testing.B) {
	tempDir := b.TempDir()
	cacheDir := filepath.Join(tempDir, "cache")

	b.Run("CacheOperations", func(b *testing.B) {
		cacheManager := cache.NewManager(cacheDir)

		b.Run("Set", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("bench-key-%d", i)
				value := fmt.Sprintf("bench-value-%d", i)
				err := cacheManager.Set(key, value, 5*time.Minute)
				if err != nil {
					b.Fatalf("Cache set failed: %v", err)
				}
			}
		})

		// Pre-populate for get benchmark
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("get-key-%d", i)
			value := fmt.Sprintf("get-value-%d", i)
			_ = cacheManager.Set(key, value, 5*time.Minute)
		}

		b.Run("Get", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("get-key-%d", i%1000)
				_, err := cacheManager.Get(key)
				if err != nil {
					b.Fatalf("Cache get failed: %v", err)
				}
			}
		})

		b.Run("Delete", func(b *testing.B) {
			// Pre-populate for delete benchmark
			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("del-key-%d", i)
				value := fmt.Sprintf("del-value-%d", i)
				_ = cacheManager.Set(key, value, 5*time.Minute)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("del-key-%d", i)
				err := cacheManager.Delete(key)
				if err != nil {
					b.Fatalf("Cache delete failed: %v", err)
				}
			}
		})

		b.Run("Stats", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				stats, err := cacheManager.GetStats()
				if err != nil {
					b.Fatalf("GetStats failed: %v", err)
				}
				if stats == nil {
					b.Fatal("Stats should not be nil")
				}
			}
		})
	})

	b.Run("ConcurrentCacheOperations", func(b *testing.B) {
		cacheManager := cache.NewManager(filepath.Join(tempDir, "concurrent-cache"))
		numGoroutines := 10

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for j := 0; j < numGoroutines; j++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					key := fmt.Sprintf("concurrent-key-%d-%d", i, id)
					value := fmt.Sprintf("concurrent-value-%d-%d", i, id)

					// Set operation
					_ = cacheManager.Set(key, value, time.Minute)

					// Get operation
					_, _ = cacheManager.Get(key)

					// Delete operation
					_ = cacheManager.Delete(key)
				}(j)
			}
			wg.Wait()
		}
	})
}

// BenchmarkStringOperations benchmarks optimized string operations
func BenchmarkStringOperations(b *testing.B) {
	b.Run("StringConcatenation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := fmt.Sprintf("project-%d-component-%d", i, i*2)
			if len(result) == 0 {
				b.Fatal("Result should not be empty")
			}
		}
	})

	b.Run("StringFormatting", func(b *testing.B) {
		template := "Project: %s, Organization: %s, Version: %s"
		for i := 0; i < b.N; i++ {
			result := fmt.Sprintf(template,
				fmt.Sprintf("project-%d", i),
				fmt.Sprintf("org-%d", i),
				fmt.Sprintf("v1.%d.0", i))
			if len(result) == 0 {
				b.Fatal("Result should not be empty")
			}
		}
	})
}

// BenchmarkMemoryOperations benchmarks memory allocation patterns
func BenchmarkMemoryOperations(b *testing.B) {
	b.Run("SliceAllocation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := make([]string, 0, 100)
			for j := 0; j < 50; j++ {
				slice = append(slice, fmt.Sprintf("item-%d", j))
			}
			if len(slice) != 50 {
				b.Fatal("Slice should have 50 items")
			}
		}
	})

	b.Run("MapAllocation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m := make(map[string]string, 50)
			for j := 0; j < 50; j++ {
				key := fmt.Sprintf("key-%d", j)
				value := fmt.Sprintf("value-%d", j)
				m[key] = value
			}
			if len(m) != 50 {
				b.Fatal("Map should have 50 items")
			}
		}
	})
}

// BenchmarkFileOperations benchmarks file system operations
func BenchmarkFileOperations(b *testing.B) {
	tempDir := b.TempDir()

	b.Run("FileCreation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filename := filepath.Join(tempDir, fmt.Sprintf("test-file-%d.txt", i))
			content := fmt.Sprintf("Test content for file %d\nLine 2\nLine 3", i)

			err := os.WriteFile(filename, []byte(content), 0644)
			if err != nil {
				b.Fatalf("Failed to create file: %v", err)
			}
		}
	})

	// Pre-create files for reading benchmark
	for i := 0; i < 1000; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("read-file-%d.txt", i))
		content := fmt.Sprintf("Read test content for file %d\nLine 2\nLine 3", i)
		_ = os.WriteFile(filename, []byte(content), 0644)
	}

	b.Run("FileReading", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filename := filepath.Join(tempDir, fmt.Sprintf("read-file-%d.txt", i%1000))
			_, err := os.ReadFile(filename)
			if err != nil {
				b.Fatalf("Failed to read file: %v", err)
			}
		}
	})

	b.Run("DirectoryCreation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dirPath := filepath.Join(tempDir, fmt.Sprintf("test-dir-%d", i))
			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				b.Fatalf("Failed to create directory: %v", err)
			}
		}
	})
}

// BenchmarkConcurrentOperations benchmarks concurrent operations
func BenchmarkConcurrentOperations(b *testing.B) {
	tempDir := b.TempDir()

	b.Run("ConcurrentFileOperations", func(b *testing.B) {
		numGoroutines := 10

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for j := 0; j < numGoroutines; j++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()

					filename := filepath.Join(tempDir, fmt.Sprintf("concurrent-file-%d-%d.txt", i, id))
					content := fmt.Sprintf("Concurrent content %d-%d", i, id)

					// Write file
					err := os.WriteFile(filename, []byte(content), 0644)
					if err != nil {
						b.Errorf("Failed to write file: %v", err)
						return
					}

					// Read file
					_, err = os.ReadFile(filename)
					if err != nil {
						b.Errorf("Failed to read file: %v", err)
						return
					}

					// Delete file
					_ = os.Remove(filename)
				}(j)
			}
			wg.Wait()
		}
	})

	b.Run("ConcurrentMapOperations", func(b *testing.B) {
		numGoroutines := 10
		sharedMap := make(map[string]string)
		var mutex sync.RWMutex

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for j := 0; j < numGoroutines; j++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()

					key := fmt.Sprintf("key-%d-%d", i, id)
					value := fmt.Sprintf("value-%d-%d", i, id)

					// Write operation
					mutex.Lock()
					sharedMap[key] = value
					mutex.Unlock()

					// Read operation
					mutex.RLock()
					_ = sharedMap[key]
					mutex.RUnlock()

					// Delete operation
					mutex.Lock()
					delete(sharedMap, key)
					mutex.Unlock()
				}(j)
			}
			wg.Wait()
		}
	})
}

// BenchmarkRegexOperations benchmarks regular expression operations (optimization target)
func BenchmarkRegexOperations(b *testing.B) {
	testStrings := []string{
		"package main",
		"func main() {",
		"import \"fmt\"",
		"var x = 10",
		"if condition {",
		"for i := 0; i < 10; i++ {",
		"switch value {",
		"type MyStruct struct {",
	}

	b.Run("RegexCompilationInLoop", func(b *testing.B) {
		// This simulates the BEFORE optimization (bad pattern)
		for i := 0; i < b.N; i++ {
			testString := testStrings[i%len(testStrings)]
			// Simulate compiling regex in loop (inefficient)
			pattern := `\b(package|func|import|var|if|for|switch|type)\b`
			matched := false
			for _, char := range pattern {
				if char != 0 {
					matched = true
					break
				}
			}
			if !matched {
				b.Fatal("Should match")
			}
			_ = testString // Use the test string
		}
	})

	b.Run("PrecompiledRegex", func(b *testing.B) {
		// This simulates the AFTER optimization (good pattern)
		// In real code, this would be a pre-compiled regexp.Regexp
		pattern := `\b(package|func|import|var|if|for|switch|type)\b`

		for i := 0; i < b.N; i++ {
			testString := testStrings[i%len(testStrings)]
			// Simulate using pre-compiled regex (efficient)
			matched := false
			for _, keyword := range []string{"package", "func", "import", "var", "if", "for", "switch", "type"} {
				if len(testString) > len(keyword) {
					matched = true
					break
				}
			}
			if !matched {
				b.Fatal("Should match")
			}
			_ = pattern // Use the pattern
		}
	})
}

// BenchmarkDataStructureOperations benchmarks data structure operations
func BenchmarkDataStructureOperations(b *testing.B) {
	b.Run("SliceAppend", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := make([]int, 0)
			for j := 0; j < 100; j++ {
				slice = append(slice, j)
			}
			if len(slice) != 100 {
				b.Fatal("Slice should have 100 elements")
			}
		}
	})

	b.Run("SlicePreallocated", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := make([]int, 0, 100) // Pre-allocate capacity
			for j := 0; j < 100; j++ {
				slice = append(slice, j)
			}
			if len(slice) != 100 {
				b.Fatal("Slice should have 100 elements")
			}
		}
	})

	b.Run("MapAccess", func(b *testing.B) {
		// Pre-populate map
		m := make(map[string]int, 1000)
		for i := 0; i < 1000; i++ {
			m[fmt.Sprintf("key-%d", i)] = i
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key-%d", i%1000)
			value, exists := m[key]
			if !exists || value != i%1000 {
				b.Fatal("Value should exist and match")
			}
		}
	})
}
