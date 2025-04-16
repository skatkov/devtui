package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/skatkov/devtui/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	// Get the root command from your cmd package
	rootCmd := cmd.GetRootCmd() // You'll need to export this

	// Only generate docs for specific commands
	cmds := rootCmd.Commands()
	for _, cmd := range cmds {
		// Create a file
		file, err := os.Create(filepath.Join("./docs", cmd.Name()+".md"))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Use the file as io.Writer
		err = doc.GenMarkdown(cmd, file)
		if err != nil {
			log.Fatal(err)
		}
	}
}
