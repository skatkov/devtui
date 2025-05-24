package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/skatkov/devtui/internal/base64"
	"github.com/spf13/cobra"
)

var base64Cmd = &cobra.Command{
	Use:   "base64 [string or file]",
	Short: "Encode or decode base64 strings and files",
	Long: `Encode or decode base64 strings and files. 

By default, input is encoded to base64. Use the --decode flag to decode base64 input.
Input can be a string, file path, or piped from stdin.`,
	Example: `  # Encode a string
  devtui base64 "hello world"
  
  # Decode a base64 string
  devtui base64 "aGVsbG8gd29ybGQ=" --decode
  devtui base64 "aGVsbG8gd29ybGQ=" -d
  
  # Encode a file
  devtui base64 /path/to/file.txt
  
  # Decode a file containing base64
  devtui base64 /path/to/base64file.txt --decode
  
  # Output to file
  devtui base64 "hello world" > encoded.txt
  devtui base64 "aGVsbG8gd29ybGQ=" --decode > decoded.txt
  
  # Pipe input from other commands
  echo -n "hello world" | devtui base64
  echo -n "aGVsbG8gd29ybGQ=" | devtui base64 --decode
  
  # Chain with other commands
  cat file.txt | devtui base64 | devtui base64 --decode`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var input []byte
		var err error

		// Check if we have stdin input
		stat, _ := os.Stdin.Stat()
		hasStdin := (stat.Mode() & os.ModeCharDevice) == 0

		if hasStdin {
			// Read from stdin
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("error reading from stdin: %v", err)
			}
		} else if len(args) == 0 {
			return fmt.Errorf("no input provided. Use a string argument, file path, or pipe input")
		} else {
			inputArg := args[0]

			// Check if input is a file path
			if fileInfo, err := os.Stat(inputArg); err == nil && !fileInfo.IsDir() {
				// It's a file, read its contents
				input, err = os.ReadFile(inputArg)
				if err != nil {
					return fmt.Errorf("error reading file '%s': %v", inputArg, err)
				}
			} else {
				// It's a string
				input = []byte(inputArg)
			}
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