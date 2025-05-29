package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"mvdan.cc/xurls/v2"
)

var urlsCmd = &cobra.Command{
	Use:   "urls [string or file]",
	Short: "Extract URLs from text, files, or stdin",
	Long: `Extract URLs from text, files, or stdin.

By default, uses relaxed mode which finds URLs without requiring a scheme.
Use the --strict flag to only find URLs with valid schemes (http, https, ftp, etc.).
Input can be a string, file path, or piped from stdin.`,
	Example: `  # Extract URLs from a string
  devtui urls "Visit https://google.com and http://example.com"

  # Extract URLs in strict mode (requires valid schemes)
  devtui urls "Visit google.com and https://example.com" --strict

  # Extract URLs from a file
  devtui urls /path/to/file.html

  # Extract URLs from stdin
  cat file.html | devtui urls
  echo "Check out google.com" | devtui urls

  # Chain with other commands
  curl -s https://example.com | devtui urls
  devtui urls file.txt > extracted_urls.txt`,
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
			return errors.New("no input provided. Use a string argument, file path, or pipe input")
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

		content := string(input)

		// Extract URLs based on mode
		var urls []string
		if urlsStrict {
			urls = xurls.Strict().FindAllString(content, -1)
		} else {
			urls = xurls.Relaxed().FindAllString(content, -1)
		}

		// Remove duplicates while preserving order
		seen := make(map[string]bool)
		uniqueURLs := make([]string, 0, len(urls))
		for _, url := range urls {
			if !seen[url] {
				seen[url] = true
				uniqueURLs = append(uniqueURLs, url)
			}
		}

		// Output results
		if len(uniqueURLs) > 0 {
			fmt.Print(strings.Join(uniqueURLs, "\n"))
			fmt.Print("\n")
		}

		return nil
	},
}

var urlsStrict bool

func init() {
	rootCmd.AddCommand(urlsCmd)

	urlsCmd.Flags().BoolVarP(&urlsStrict, "strict", "s", false, "use strict mode (require valid URL schemes)")
}
