package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/fang/v2"
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
		return runProgram(root.RootScreen())
	},
}

// runProgram runs a TUI; tests replace it so they never open one.
var runProgram = func(model tea.Model) error {
	_, err := tea.NewProgram(model).Run()
	return err
}

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

	// The transform applies only to the first line of the first paragraph.
	lineStyle := styles.ErrorText.UnsetWidth()
	separator := ""

	// Split error message by double newlines to preserve paragraph breaks
	for _, para := range strings.Split(err.Error(), "\n\n") {
		// Add blank line between paragraphs (but not after the last one)
		_, _ = fmt.Fprint(w, separator)
		separator = "\n"

		for _, line := range strings.Split(para, "\n") {
			style := lineStyle
			lineStyle = baseStyle
			if line == "" {
				continue
			}
			_, _ = fmt.Fprintln(w, style.Render(line))
		}
	}
	_, _ = fmt.Fprintln(w)
}

// exit ends the process; tests replace it.
var exit = os.Exit

func Execute() {
	err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(GetVersionShort()),
		fang.WithErrorHandler(customErrorHandler),
	)
	if err != nil {
		exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// fang.Execute automatically adds --version flag
	// We configure it via fang.WithVersion() in Execute()
}
