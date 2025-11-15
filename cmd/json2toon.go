package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/hannes-sistemica/toon"
	"github.com/skatkov/devtui/tui/json2toon"
	"github.com/spf13/cobra"
)

var json2toonCmd = &cobra.Command{
	Use:   "json2toon",
	Short: "Convert JSON to TOON (Token-Oriented Object Notation)",
	Long: `Convert JSON to TOON (Token-Oriented Object Notation) - a compact, human-readable
format designed for passing structured data to Large Language Models with significantly
reduced token usage (typically 30-60% fewer tokens than JSON).`,
	Example: `  devtui json2toon < example.json                    # Convert with defaults
  devtui json2toon -i 4 < example.json               # Use 4-space indent
  devtui json2toon -d tab < example.json             # Use tab delimiter
  devtui json2toon -d pipe < example.json            # Use pipe delimiter
  devtui json2toon -l '#' < example.json             # Add length marker prefix
  cat example.json | devtui json2toon > output.toon  # Pipe and save to file`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		// Parse delimiter flag
		var delimiter string
		switch json2toonDelimiter {
		case "comma":
			delimiter = ","
		case "tab":
			delimiter = "\t"
		case "pipe":
			delimiter = "|"
		default:
			return fmt.Errorf("invalid delimiter: %s (must be comma, tab, or pipe)", json2toonDelimiter)
		}

		opts := toon.EncodeOptions{
			Indent:       json2toonIndent,
			Delimiter:    delimiter,
			LengthMarker: json2toonLengthMarker,
		}

		result, err := json2toon.ConvertWithOptions(string(data), opts)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

var (
	json2toonIndent       int    // Number of spaces per indentation level
	json2toonDelimiter    string // Delimiter for arrays: comma, tab, or pipe
	json2toonLengthMarker string // Optional marker to prefix array lengths
)

func init() {
	rootCmd.AddCommand(json2toonCmd)

	json2toonCmd.Flags().IntVarP(&json2toonIndent, "indent", "i", 2, "Number of spaces per indentation level")
	json2toonCmd.Flags().StringVarP(&json2toonDelimiter, "delimiter", "d", "comma", "Delimiter for arrays: comma, tab, or pipe")
	json2toonCmd.Flags().StringVarP(&json2toonLengthMarker, "length-marker", "l", "", "Optional marker to prefix array lengths (e.g., '#')")
}
