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
	Run: func(cmd *cobra.Command, args []string) {
		var text string
		var err error

		if len(args) > 0 {
			text = args[0]
		} else {
			reader := bufio.NewReader(os.Stdin)
			text, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		stats, err := textanalyzer.Analyze(text)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(table.New().Border(lipgloss.NormalBorder()).
			Row("Characters", strconv.Itoa(stats.Characters)).
			Row("Spaces", strconv.Itoa(stats.Spaces)).
			Row("Words", strconv.Itoa(stats.Words)))
	},
}

func init() {
	rootCmd.AddCommand(countCmd)
}
