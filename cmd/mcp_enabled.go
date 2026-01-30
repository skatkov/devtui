//go:build mcp

package cmd

import (
	mcp "github.com/skatkov/devtui-mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run DevTUI as an MCP stdio server",
	RunE: func(cmd *cobra.Command, args []string) error {
		tools := mcp.BuildTools(GetRootCmd())
		server := mcp.NewServer(mcp.ServerConfig{
			Tools: tools,
			ServerInfo: mcp.ServerInfo{
				Name:    "devtui",
				Version: GetVersion(),
			},
			Call: func(_ string, params mcp.CallParams) (string, error) {
				root := GetRootCmd()
				return mcp.ExecuteTool(root, params)
			},
		})

		return mcp.ServeStdio(server, cmd.InOrStdin(), cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
