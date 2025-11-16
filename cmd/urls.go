package cmd

import (
	"fmt"
	"strings"

	"github.com/skatkov/devtui/internal/input"
	"github.com/spf13/cobra"
	"mvdan.cc/xurls/v2"
)

var urlsCmd = &cobra.Command{
	Use:   "urls [string or file]",
	Short: "Extract URLs from text, files, or stdin",
	Long: `Extract URLs from text, files, or stdin.

By default, uses relaxed mode which finds URLs without requiring a scheme.
Use the --strict flag to only find URLs with valid schemes (http, https, ftp, etc.).
Input can be a string argument or piped from stdin.`,
	Example: `  # Extract URLs from a string
  devtui urls "Visit https://google.com and http://example.com"

  # Extract URLs in strict mode (requires valid schemes)
  devtui urls "Visit google.com and https://example.com" --strict

  # Extract URLs from stdin
  cat file.html | devtui urls
  echo "Check out google.com" | devtui urls

  # Chain with other commands
  curl -s https://example.com | devtui urls
  cat file.txt | devtui urls > extracted_urls.txt`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := input.ReadFromArgsOrStdin(cmd, args)
		if err != nil {
			return err
		}

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
