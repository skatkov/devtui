//go:build !mcp

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run DevTUI as an MCP stdio server",
	RunE: func(cmd *cobra.Command, args []string) error {
		message := "mcp disabled; rebuild with -tags mcp"
		fmt.Fprintln(cmd.ErrOrStderr(), message)
		return fmt.Errorf("%s", message)
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
