package cmd

import (
	"errors"
	"fmt"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/xml"
	"github.com/spf13/cobra"
)

var xml2jsonCmd = &cobra.Command{
	Use:   "xml2json [string or file]",
	Short: "Convert XML to JSON format",
	Long: `Convert XML (Extensible Markup Language) to JSON (JavaScript Object Notation) format.

Input can be a string argument or piped from stdin.`,
	Example: `  # Convert XML from stdin
  devtui xml2json < data.xml
  cat feed.xml | devtui xml2json

  # Convert XML string argument
  devtui xml2json '<root><item>value</item></root>'

  # Output to file
  devtui xml2json < input.xml > output.json

  # Chain with other commands
  curl -s https://api.example.com/data.xml | devtui xml2json`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		if len(data) == 0 {
			return errors.New("no input provided. Pipe XML input to this command")
		}

		inputStr := string(data)
		result, err := xml.XMLToJSON(inputStr)
		if err != nil {
			return cmderror.FormatParseError("xml2json", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(xml2jsonCmd)
}
