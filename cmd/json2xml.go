package cmd

import (
	"errors"
	"fmt"

	"github.com/skatkov/devtui/internal/cmderror"
	"github.com/skatkov/devtui/internal/converter"
	"github.com/skatkov/devtui/internal/input"
	"github.com/spf13/cobra"
)

var json2xmlCmd = &cobra.Command{
	Use:   "json2xml [string or file]",
	Short: "Convert JSON to XML format",
	Long: `Convert JSON to XML format.

Input can be a string argument or piped from stdin.`,
	Example: `  # Convert JSON from stdin
  devtui json2xml < data.json
  cat feed.json | devtui json2xml

  # Convert JSON string argument
  devtui json2xml '{"item": "value"}'

  # Output to file
  devtui json2xml < input.json > output.xml

  # Chain with other commands
  curl -s https://api.example.com/data.json | devtui json2xml`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}

		if len(data) == 0 {
			return errors.New("no input provided. Pipe JSON input to this command")
		}

		inputStr := string(data)
		result, err := converter.JSONToXML(inputStr)
		if err != nil {
			return cmderror.FormatParseError("json2xml", inputStr, err)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(json2xmlCmd)
}
