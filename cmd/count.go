package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/skatkov/devtui/internal/textanalyzer"
	"github.com/spf13/cobra"
)

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Character, spaces and word counter",
	Long:  "Count characters, spaces and words in a string",
	Example: `  # Count text from a string
  devtui count "test me please"

  # Count text from stdin
  cat testdata/example.csv | devtui count
  echo "hello world" | devtui count`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var text string

		if len(args) > 0 {
			// Use string argument
			text = args[0]
		} else {
			// Read from stdin
			data, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return err
			}
			text = string(data)
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
