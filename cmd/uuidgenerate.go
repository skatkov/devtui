package cmd

import (
	"fmt"

	"github.com/skatkov/devtui/internal/uuidutil"
	"github.com/spf13/cobra"
)

var uuidgenerateCmd = &cobra.Command{
	Use:   "uuidgenerate",
	Short: "Generate a UUID",
	Long: `Generate a UUID of a specified version.

By default, generates a version 4 UUID. Versions 3 and 5 accept a namespace value.`,
	Example: `  # Generate a default UUID (v4)
  devtui uuidgenerate

  # Generate a UUID v7
  devtui uuidgenerate --uuid-version 7

  # Generate a UUID v3 with namespace
  devtui uuidgenerate --uuid-version 3 --namespace example.com`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		generated, err := uuidutil.Generate(uuidgenerateVersion, uuidgenerateNamespace)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(cmd.OutOrStdout(), generated.String())
		return err
	},
}

var (
	uuidgenerateVersion   int
	uuidgenerateNamespace string
)

func init() {
	rootCmd.AddCommand(uuidgenerateCmd)
	uuidgenerateCmd.Flags().IntVarP(&uuidgenerateVersion, "uuid-version", "v", 4, "UUID version to generate (1-7)")
	uuidgenerateCmd.Flags().StringVarP(&uuidgenerateNamespace, "namespace", "n", "", "namespace for UUID v3/v5 generation")
}
