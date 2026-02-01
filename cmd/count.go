package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/textanalyzer"
	"github.com/spf13/cobra"
)

var countCmd = &cobra.Command{
	Use:   "count [string or file]",
	Short: "Count characters, spaces, and words in text",
	Long: `Count characters, spaces, and words in text input.

Provides detailed statistics including character count, space count, and word count
in a formatted table. Input can be a string argument or piped from stdin.`,
	Example: `  # Count text from a string
  devtui count "test me please"
  devtui count "hello world"

  # Count text from stdin
  echo "hello world" | devtui count
  cat document.txt | devtui count

  # Count text from file
  devtui count < document.txt
  cat README.md | devtui count

  # Output to file
  devtui count "sample text" > stats.txt

  # Chain with other commands
  curl -s https://example.com | devtui count
  cat article.txt | devtui count`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text, err := input.ReadFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
		}

		stats, err := textanalyzer.Analyze(text)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), table.New().Border(lipgloss.NormalBorder()).
			Row("Characters", strconv.Itoa(stats.Characters)).
			Row("Spaces", strconv.Itoa(stats.Spaces)).
			Row("Words", strconv.Itoa(stats.Words)))
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(countCmd)
}
