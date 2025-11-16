package input

import (
	"io"

	"github.com/spf13/cobra"
)

// ReadFromStdin reads all input from the command's stdin (or actual stdin if not overridden).
// This is the standard pattern for formatter/transformer commands that only accept stdin.
//
// Example usage:
//
//	data, err := input.ReadFromStdin(cmd)
//	if err != nil {
//	    return err
//	}
func ReadFromStdin(cmd *cobra.Command) ([]byte, error) {
	return io.ReadAll(cmd.InOrStdin())
}

// ReadFromArgsOrStdin reads input from either:
// 1. The first argument (if provided), or
// 2. Stdin (if no arguments)
//
// This is the standard pattern for commands that accept quick string arguments
// but can also work as filters with piped input.
//
// Example usage:
//
//	content, err := input.ReadFromArgsOrStdin(cmd, args)
//	if err != nil {
//	    return err
//	}
func ReadFromArgsOrStdin(cmd *cobra.Command, args []string) (string, error) {
	if len(args) > 0 {
		// Use the first argument as input
		return args[0], nil
	}

	// Read from stdin
	data, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ReadBytesFromArgsOrStdin is like ReadFromArgsOrStdin but returns bytes.
// Useful for commands that need to preserve binary data.
//
// Example usage:
//
//	data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
//	if err != nil {
//	    return err
//	}
func ReadBytesFromArgsOrStdin(cmd *cobra.Command, args []string) ([]byte, error) {
	if len(args) > 0 {
		// Use the first argument as input
		return []byte(args[0]), nil
	}

	// Read from stdin
	return io.ReadAll(cmd.InOrStdin())
}
