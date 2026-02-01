//go:build !mcp

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run DevTUI as an MCP stdio server",
	RunE: func(cmd *cobra.Command, args []string) error {
		message := "mcp disabled; rebuild with -tags mcp"
		if _, err := fmt.Fprintln(cmd.ErrOrStderr(), message); err != nil {
			return err
		}
		return errors.New(message)
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
