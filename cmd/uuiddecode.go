package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/google/uuid"
	"github.com/skatkov/devtui/internal/input"
	"github.com/skatkov/devtui/internal/uuidutil"
	"github.com/spf13/cobra"
)

type uuidDecodeJSON struct {
	UUID   string           `json:"uuid"`
	Fields []uuidutil.Field `json:"fields"`
}

var uuiddecodeCmd = &cobra.Command{
	Use:   "uuiddecode [uuid]",
	Short: "Decode a UUID into its components",
	Long: `Decode a UUID and show its components, including version and variant.

Input can be provided as an argument or piped from stdin.`,
	Example: `  # Decode a UUID argument
  devtui uuiddecode 4326ff5f-774d-4506-a18c-4bc50c761863

  # Decode a UUID from stdin
  echo "4326ff5f-774d-4506-a18c-4bc50c761863" | devtui uuiddecode

  # Output as JSON
  devtui uuiddecode --json 4326ff5f-774d-4506-a18c-4bc50c761863`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputStr, err := input.ReadFromArgsOrStdin(cmd, args)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		uuidStr := strings.TrimSpace(inputStr)
		if uuidStr == "" {
			return errors.New("no input provided. pipe UUID input to this command")
		}

		parsed, err := uuid.Parse(uuidStr)
		if err != nil {
			return fmt.Errorf("invalid uuid: %w", err)
		}

		fields := uuidutil.Decode(parsed)
		if uuiddecodeJSONOutput {
			payload := uuidDecodeJSON{
				UUID:   parsed.String(),
				Fields: fields,
			}
			bytes, err := json.MarshalIndent(payload, "", "  ")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			return err
		}

		tableOutput := table.New().Border(lipgloss.NormalBorder())
		for _, row := range uuidutil.FieldsToRows(fields) {
			tableOutput.Row(row...)
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), tableOutput.String())
		return err
	},
}

var uuiddecodeJSONOutput bool

func init() {
	rootCmd.AddCommand(uuiddecodeCmd)
	uuiddecodeCmd.Flags().BoolVar(&uuiddecodeJSONOutput, "json", false, "output decoded fields as JSON")
}
