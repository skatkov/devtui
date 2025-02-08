package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "dlink",
    Short: "A collection of developer utilities",
    Long: `dlink is a collection of small developer apps that help with day to day work.
It includes tools like hash generator, unix timestamp converter, and number base converter.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    // Add persistent flags here if needed
}
