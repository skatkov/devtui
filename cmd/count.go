package cmd

import (
	"bufio"
	"fmt"
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
	count "test me please"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var text string
		var err error

		if len(args) > 0 {
			text = args[0]
		} else {
			reader := bufio.NewReader(os.Stdin)
			text, err = reader.ReadString('\n')
			if err != nil {
				return err
			}
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
