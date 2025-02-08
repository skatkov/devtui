package cmd

import (
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Hash generator utility",
	Long:  `Generate various types of hashes (MD5, SHA1, SHA256, etc.)`,
	Run: func(cmd *cobra.Command, args []string) {
		hash.StartUI()
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
}
