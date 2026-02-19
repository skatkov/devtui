package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	Example: `count < testdata/example.csv
	count "test me please"
	count --json "test me please"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var text string

		if len(args) > 0 {
			text = args[0]
		} else {
			// Read all from stdin, not just until newline
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			text = string(data)
		}

		stats, err := textanalyzer.Analyze(text)
		if err != nil {
			return err
		}

		// Output in JSON format if flag is set
		if flagJSON {
			output := map[string]int{
				"characters": stats.Characters,
				"spaces":     stats.Spaces,
				"words":      stats.Words,
			}
			jsonBytes, err := json.Marshal(output)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(jsonBytes))
			if err != nil {
				return err
			}
			return nil
		}

		// Default table output
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
