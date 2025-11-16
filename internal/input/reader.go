package input

import (
	"io"

	"github.com/spf13/cobra"
)

// ReadFromArgsOrStdin reads input from either:
// 1. The first argument (if provided), or
// 2. Stdin (if no arguments)
//
// This is the standard pattern for all DevTUI formatter/converter commands.
// Commands should use Args: cobra.MaximumNArgs(1) to enforce 0 or 1 arguments.
//
// Returns string, which is suitable for text-based commands.
//
// Example usage:
//
//	var myCmd = &cobra.Command{
//	    Args: cobra.MaximumNArgs(1),
//	    RunE: func(cmd *cobra.Command, args []string) error {
//	        content, err := input.ReadFromArgsOrStdin(cmd, args)
//	        if err != nil {
//	            return err
//	        }
//	        // Process content...
//	        return nil
//	    },
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

// ReadBytesFromArgsOrStdin reads input from either args or stdin, returning bytes.
// This is identical to ReadFromArgsOrStdin but returns []byte instead of string.
//
// Use this when you need to preserve binary data or pass data to APIs that expect []byte.
// Commands should use Args: cobra.MaximumNArgs(1) to enforce 0 or 1 arguments.
//
// Example usage:
//
//	var myCmd = &cobra.Command{
//	    Args: cobra.MaximumNArgs(1),
//	    RunE: func(cmd *cobra.Command, args []string) error {
//	        data, err := input.ReadBytesFromArgsOrStdin(cmd, args)
//	        if err != nil {
//	            return err
//	        }
//	        // Process data...
//	        return nil
//	    },
//	}
func ReadBytesFromArgsOrStdin(cmd *cobra.Command, args []string) ([]byte, error) {
	if len(args) > 0 {
		// Use the first argument as input
		return []byte(args[0]), nil
	}

	// Read from stdin
	return io.ReadAll(cmd.InOrStdin())
}
