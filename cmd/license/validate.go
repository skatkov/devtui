package license

import (
	"fmt"
	"os"

	license "github.com/skatkov/devtui/internal/license"
	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate a license",
	Long:    "Reads a license and validates it",
	Example: "devtui license validate",
	Run: func(cmd *cobra.Command, args []string) {
		licenseData, err := license.LoadLicense()
		if err != nil {
			fmt.Printf("Error: no active license found")
			os.Exit(0)
		}

		if licenseData == nil {
			fmt.Println("Error: no active license found")
			os.Exit(0)
		}

		if err := licenseData.Validate(); err != nil {
			fmt.Println("Error: license file is not valid")
			os.Exit(0)
		}
		fmt.Println("License is valid")
	},
}
