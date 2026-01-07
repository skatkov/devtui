package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/term"
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

// customErrorHandler handles errors while preserving newlines for multiline messages.
// The default fang error handler applies a fixed width that collapses newlines,
// and transforms text (capitalizing first word). We want to preserve the structure
// of multiline error messages, especially hints.
func customErrorHandler(w io.Writer, styles fang.Styles, err error) {
	if f, ok := w.(term.File); ok {
		if !term.IsTerminal(f.Fd()) {
			_, _ = fmt.Fprintln(w, err.Error())
			return
		}
	}

	_, _ = fmt.Fprintln(w, styles.ErrorHeader.String())

	// Get base style without width constraint and without transform
	// (UnsetTransform prevents auto-capitalizing "devtui" to "Devtui")
	baseStyle := styles.ErrorText.UnsetWidth().UnsetTransform()

	// Split error message by double newlines to preserve paragraph breaks
	errMsg := err.Error()
	paragraphs := strings.Split(errMsg, "\n\n")

	for i, para := range paragraphs {
		// Render each paragraph separately to preserve structure
		lines := strings.Split(para, "\n")
		for j, line := range lines {
			if line == "" {
				continue
			}
			// Apply transform only to first line of first paragraph
			var styled string
			if i == 0 && j == 0 {
				styled = styles.ErrorText.UnsetWidth().Render(line)
			} else {
				styled = baseStyle.Render(line)
			}
			_, _ = fmt.Fprintln(w, styled)
		}
		// Add blank line between paragraphs (but not after the last one)
		if i < len(paragraphs)-1 {
			_, _ = fmt.Fprintln(w)
		}
	}
	_, _ = fmt.Fprintln(w)
}

func Execute() {
	err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(GetVersionShort()),
		fang.WithErrorHandler(customErrorHandler),
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
