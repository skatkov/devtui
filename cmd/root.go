package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/fang"
	"github.com/skatkov/devtui/tui/root"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var rootCmd = &cobra.Command{
	Use:   "devtui",
	Short: "A Swiss Army knife for developers",
	Long: `devtui is a collection of small developer apps that help with day to day work.
It includes tools like hash generator, unix timestamp converter, and number base converter and multiple others.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(root.RootScreen(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

var (
	flagTUI   bool
	flagJSON  bool
)

// Execute runs the root command
func Execute() {
	err := fang.Execute(context.Background(), rootCmd)
	if err != nil {
		os.Exit(1)
	}
}

// GetRootCmd returns the root command instance
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// Add persistent flags available to all subcommands
	rootCmd.PersistentFlags().BoolVarP(&flagJSON, "json", "j", false, "Output results in JSON format")
}
