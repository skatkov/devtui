package main

import (
	"log"

	"github.com/skatkov/devtui/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	// Get the root command from your cmd package
	rootCmd := cmd.GetRootCmd() // You'll need to export this

	// Generate markdown documentation
	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
