package security_test

import (
	"fmt"
	"log"

	"github.com/open-source-template-generator/pkg/security"
)

func ExampleSecureRandom_GenerateRandomSuffix() {
	sr := security.NewSecureRandom()

	// Generate a 16-character hex suffix
	suffix, err := sr.GenerateRandomSuffix(16)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Random suffix length: %d\n", len(suffix))
	// Output: Random suffix length: 16
}

func ExampleSecureRandom_GenerateSecureID() {
	sr := security.NewSecureRandom()

	// Generate a secure ID with prefix
	id, err := sr.GenerateSecureID("audit")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Secure ID starts with 'audit_': %t\n", len(id) > 6 && id[:6] == "audit_")
	// Output: Secure ID starts with 'audit_': true
}

func ExampleGenerateRandomSuffix() {
	// Using the global convenience function
	suffix, err := security.GenerateRandomSuffix(12)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated suffix length: %d\n", len(suffix))
	// Output: Generated suffix length: 12
}

func ExampleGenerateSecureID() {
	// Using the global convenience function
	id, err := security.GenerateSecureID("temp")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated ID has prefix: %t\n", len(id) > 5 && id[:5] == "temp_")
	// Output: Generated ID has prefix: true
}

func ExampleSecureRandom_GenerateAlphanumeric() {
	sr := security.NewSecureRandomWithConfig(16, "alphanumeric")

	// Generate alphanumeric string
	alphaStr, err := sr.GenerateAlphanumeric(8)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Alphanumeric string length: %d\n", len(alphaStr))
	// Output: Alphanumeric string length: 8
}
