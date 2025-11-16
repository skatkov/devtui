package cmd

import (
	"fmt"
	"io"

	"github.com/skatkov/devtui/internal/base64"
	"github.com/spf13/cobra"
)

var base64Cmd = &cobra.Command{
	Use:   "base64 [string or file]",
	Short: "Encode or decode base64 strings and files",
	Long: `Encode or decode base64 strings and files.

By default, input is encoded to base64. Use the --decode flag to decode base64 input.
Input can be a string argument or piped from stdin.`,
	Example: `  # Encode a string
  devtui base64 "hello world"

  # Decode a base64 string
  devtui base64 "aGVsbG8gd29ybGQ=" --decode
  devtui base64 "aGVsbG8gd29ybGQ=" -d

  # Output to file
  devtui base64 "hello world" > encoded.txt
  devtui base64 "aGVsbG8gd29ybGQ=" --decode > decoded.txt

  # Pipe input from other commands
  echo -n "hello world" | devtui base64
  echo -n "aGVsbG8gd29ybGQ=" | devtui base64 --decode
  cat file.txt | devtui base64

  # Chain with other commands
  cat file.txt | devtui base64 | devtui base64 --decode`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var input []byte

		if len(args) > 0 {
			// Use string argument
			input = []byte(args[0])
		} else {
			// Read from stdin
			data, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return err
			}
			input = data
		}

		// Perform encoding or decoding
		if base64Decode {
			// Decode base64
			decoded, err := base64.DecodeToString(string(input))
			if err != nil {
				return err
			}
			fmt.Print(decoded)
		} else {
			// Encode to base64
			encoded := base64.Encode(input)
			fmt.Print(encoded)
		}

		return nil
	},
}

var base64Decode bool

func init() {
	rootCmd.AddCommand(base64Cmd)

	base64Cmd.Flags().BoolVarP(&base64Decode, "decode", "d", false, "decode base64 input instead of encoding")
}
