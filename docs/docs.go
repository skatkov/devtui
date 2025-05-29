package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Generating documentation...")

	// Generate CLI documentation
	fmt.Println("Generating CLI documentation...")
	GenerateCLIDocumentation()

	// Generate TUI documentation
	fmt.Println("Generating TUI documentation...")
	if err := GenerateTUIDocumentation(); err != nil {
		log.Fatalf("Failed to generate TUI documentation: %v", err)
	}

	fmt.Println("Documentation generation completed successfully!")
}
