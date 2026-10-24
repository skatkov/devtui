package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Generate writes the CLI and TUI documentation into the site directory.
func Generate(rootCmd *cobra.Command) {
	fmt.Println("Generating documentation...")

	// Generate CLI documentation
	fmt.Println("Generating CLI documentation...")
	GenerateCLIDocumentation(rootCmd, filepath.Join(projectPath("site"), "cli"))

	// Generate TUI documentation
	fmt.Println("Generating TUI documentation...")
	if err := GenerateTUIDocumentation(projectPath("tui"), filepath.Join(projectPath("site"), "tui")); err != nil {
		log.Fatalf("Failed to generate TUI documentation: %v", err)
	}

	fmt.Println("Documentation generation completed successfully!")
}

// projectPath resolves name from the repository root or from the docs directory.
func projectPath(name string) string {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return filepath.Join("..", name)
	}
	return name
}
