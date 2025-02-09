package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skatkov/devtui/tui/root"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devtui",
	Short: "Swiss army knife for developers",
	Long: `devtui is a collection of small developer apps that help with day to day work.
	It includes tools like hash generator, unix timestamp converter, and number base converter and multiple others.`,
}

func Execute() {
	p := tea.NewProgram(root.RootScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting program: %v", err)
		os.Exit(1)
	}
}

func init() {
	// Add persistent flags here if needed
}
