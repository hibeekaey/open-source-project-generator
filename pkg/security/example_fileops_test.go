package security

import (
	"fmt"
	"os"
	"path/filepath"
)

// ExampleWriteFileAtomic demonstrates the convenience function for atomic file writing
func ExampleWriteFileAtomic() {
	tempDir := os.TempDir()
	testFile := filepath.Join(tempDir, "convenience_example.txt")

	// Use the convenience function for atomic writing
	data := []byte("Using the convenience function for secure atomic writes")
	err := WriteFileAtomic(testFile, data, 0644)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Clean up
	defer os.Remove(testFile)

	fmt.Printf("File written atomically using convenience function\n")
	// Output: File written atomically using convenience function
}

// ExampleCreateSecureTempFile demonstrates the convenience function for secure temp files
func ExampleCreateSecureTempFile() {
	tempDir := os.TempDir()

	// Use the convenience function for secure temp file creation
	tempFile, err := CreateSecureTempFile(tempDir, "convenience.")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	fmt.Printf("Secure temp file created using convenience function\n")
	// Output: Secure temp file created using convenience function
}

// ExampleValidatePath demonstrates path validation for directory traversal
func ExampleValidatePath() {
	// Test directory traversal attempt - should be rejected
	err := ValidatePath("../../../etc/passwd", nil)
	if err != nil {
		fmt.Printf("Directory traversal: REJECTED\n")
	} else {
		fmt.Printf("Directory traversal: ALLOWED\n")
	}
	// Output: Directory traversal: REJECTED
}
