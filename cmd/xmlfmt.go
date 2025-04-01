package cmd

import (
	"io"
	"log"
	"os"

	"github.com/go-xmlfmt/xmlfmt"
	"github.com/spf13/cobra"
)

// ---
// Example usage:
// go run . xmlfmt < testdata/sample.xml
//
// Output to file:
// go run . xmlfmt < testdata/sample.xml > output.xml
// ---

var xmlfmtCmd = &cobra.Command{
	Use:   "xmlfmt",
	Short: "Format XML files",
	Long:  "Format XML files",
	Run: func(cmd *cobra.Command, args []string) {
		// Read all input data from stdin
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		// Process the XML
		result := xmlfmt.FormatXML(string(data),
			xmlPrefix, xmlIndent, xmlNested)

		// Write the result to stdout
		_, err = os.Stdout.WriteString(result)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
	},
}

var (
	xmlPrefix string
	xmlIndent string
	xmlNested bool
)

func init() {
	rootCmd.AddCommand(xmlfmtCmd)
	xmlfmtCmd.Flags().StringVarP(&xmlPrefix, "prefix", "p", "", "Each element begins on a new line and this prefix")
	xmlfmtCmd.Flags().StringVarP(&xmlIndent, "indent", "i", "  ", "Indent string for nested elements")
	xmlfmtCmd.Flags().BoolVarP(&xmlNested, "nested", "n", false, "Nested tags in comments")
}
