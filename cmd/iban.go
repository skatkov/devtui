package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jacoelho/banking/iban"
	"github.com/spf13/cobra"
)

var ibanCmd = &cobra.Command{
	Use:   "iban <country-code>",
	Short: "Generate test IBAN numbers",
	Long: `Generate test IBAN numbers for a given country code using the banking library.

The country code is required and should be a valid ISO 3166-1 alpha-2 country code.
Use the --formatted flag to output the IBAN in paper format with spaces.`,
	Example: `  # Generate IBAN for Great Britain
  devtui iban GB

  # Generate formatted IBAN for Germany
  devtui iban DE --format
  devtui iban DE -f

  # Generate IBAN for France
  devtui iban FR`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		countryCode := strings.ToUpper(strings.TrimSpace(args[0]))

		if len(countryCode) != 2 {
			return errors.New("country code must be exactly 2 characters")
		}

		// Generate IBAN
		generatedIban, err := iban.Generate(countryCode)
		if err != nil {
			return fmt.Errorf("failed to generate IBAN for country code '%s': %v", countryCode, err)
		}

		// Format output based on flag
		output := generatedIban
		if ibanFormatted {
			output = iban.PaperFormat(generatedIban)
		}
		if outputJSON {
			return writeJSONValue(cmd.OutOrStdout(), output)
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), output)
		return err
	},
}

var ibanFormatted bool

func init() {
	rootCmd.AddCommand(ibanCmd)

	ibanCmd.Flags().BoolVarP(&ibanFormatted, "format", "f", false, "output IBAN in paper format with spaces")
}
