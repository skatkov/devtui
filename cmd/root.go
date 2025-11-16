package cmd

import (
	"context"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/fang"
	"github.com/skatkov/devtui/tui/root"
	"github.com/spf13/cobra"
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

var flagTUI bool

func Execute() {
	err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(GetVersionShort()),
	)
	if err != nil {
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// fang.Execute automatically adds --version flag
	// We configure it via fang.WithVersion() in Execute()
}
