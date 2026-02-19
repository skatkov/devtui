package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestJSONOutput_CountCommand(t *testing.T) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "devtui_test_json")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build devtui: %v", err)
	}
	defer os.Remove("../devtui_test_json")

	binary := "../devtui_test_json"

	tests := []struct {
		name     string
		args     []string
		input    string
		validate func(t *testing.T, output string)
	}{
		{
			name: "count with JSON output - simple text",
			args: []string{"count", "--json", "hello world"},
			validate: func(t *testing.T, output string) {
				var result map[string]int
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Fatalf("Failed to parse JSON output: %v", err)
				}
				if result["characters"] != 11 {
					t.Errorf("Expected characters=11, got %d", result["characters"])
				}
				if result["spaces"] != 1 {
					t.Errorf("Expected spaces=1, got %d", result["spaces"])
				}
				if result["words"] != 2 {
					t.Errorf("Expected words=2, got %d", result["words"])
				}
			},
		},
		{
			name: "count with JSON short flag",
			args: []string{"count", "-j", "test"},
			validate: func(t *testing.T, output string) {
				var result map[string]int
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Fatalf("Failed to parse JSON output: %v", err)
				}
				if result["characters"] != 4 {
					t.Errorf("Expected characters=4, got %d", result["characters"])
				}
				if result["words"] != 1 {
					t.Errorf("Expected words=1, got %d", result["words"])
				}
			},
		},
		{
			name: "count with JSON output from stdin",
			args: []string{"count", "--json"},
			input: "hello world",
			validate: func(t *testing.T, output string) {
				var result map[string]int
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Fatalf("Failed to parse JSON output: %v", err)
				}
				if result["characters"] != 11 {
					t.Errorf("Expected characters=11, got %d", result["characters"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if tt.input != "" {
				cmd.Stdin = strings.NewReader(tt.input)
			}

			err := cmd.Run()
			if err != nil {
				t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
			}

			output := strings.TrimSpace(stdout.String())
			tt.validate(t, output)
		})
	}
}

func TestJSONOutput_VersionCommand(t *testing.T) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "devtui_test_json")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build devtui: %v", err)
	}
	defer os.Remove("../devtui_test_json")

	binary := "../devtui_test_json"

	cmd := exec.Command(binary, "version", "--json")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())

	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
	}

	// Check that expected keys exist
	expectedKeys := []string{"version", "commit", "date"}
	for _, key := range expectedKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("Expected key '%s' not found in JSON output", key)
		}
	}
}

func TestJSONOutput_IBANCommand(t *testing.T) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "devtui_test_json")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build devtui: %v", err)
	}
	defer os.Remove("../devtui_test_json")

	binary := "../devtui_test_json"

	tests := []struct {
		name   string
		args   []string
		hasErr bool
	}{
		{
			name: "iban with JSON output - GB",
			args: []string{"iban", "GB", "--json"},
		},
		{
			name: "iban with JSON output - DE",
			args: []string{"iban", "DE", "--json"},
		},
		{
			name: "iban with JSON and formatted output - FR",
			args: []string{"iban", "FR", "--json", "--format"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if tt.hasErr {
				if err == nil {
					t.Error("Expected error but command succeeded")
				}
				return
			}

			if err != nil {
				t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
			}

			output := strings.TrimSpace(stdout.String())

			var result map[string]string
			if err := json.Unmarshal([]byte(output), &result); err != nil {
				t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
			}

			// Check expected keys
			if _, ok := result["country_code"]; !ok {
				t.Error("Expected key 'country_code' not found in JSON output")
			}
			if _, ok := result["iban"]; !ok {
				t.Error("Expected key 'iban' not found in JSON output")
			}

			// Check if formatted key exists when --format flag is used
			for _, arg := range tt.args {
				if arg == "--format" || arg == "-f" {
					if _, ok := result["formatted"]; !ok {
						t.Error("Expected key 'formatted' not found in JSON output when --format flag is used")
					}
				}
			}
		})
	}
}

func TestJSONOutput_URLsCommand(t *testing.T) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "devtui_test_json")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build devtui: %v", err)
	}
	defer os.Remove("../devtui_test_json")

	binary := "../devtui_test_json"

	tests := []struct {
		name       string
		args       []string
		input      string
		expectURLs []string
	}{
		{
			name:       "urls with JSON output",
			args:       []string{"urls", "--json"},
			input:      "Visit https://google.com and http://example.com",
			expectURLs: []string{"https://google.com", "http://example.com"},
		},
		{
			name:       "urls with JSON output - single URL",
			args:       []string{"urls", "--json"},
			input:      "Check out https://github.com/skatkov/devtui",
			expectURLs: []string{"https://github.com/skatkov/devtui"},
		},
		{
			name:       "urls with JSON output - no URLs",
			args:       []string{"urls", "--json"},
			input:      "No URLs here",
			expectURLs: []string{},
		},
		{
			name:       "urls with JSON and strict mode",
			args:       []string{"urls", "--json", "--strict"},
			input:      "Visit https://google.com and just google.com",
			expectURLs: []string{"https://google.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			cmd.Stdin = strings.NewReader(tt.input)

			err := cmd.Run()
			if err != nil {
				t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
			}

			output := strings.TrimSpace(stdout.String())

			var result []string
			if err := json.Unmarshal([]byte(output), &result); err != nil {
				t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
			}

			if len(result) != len(tt.expectURLs) {
				t.Errorf("Expected %d URLs, got %d: %v", len(tt.expectURLs), len(result), result)
			}

			for i, url := range tt.expectURLs {
				if i >= len(result) || result[i] != url {
					t.Errorf("Expected URL %d to be %q, got %q", i, url, result[i])
				}
			}
		})
	}
}

func TestJSONOutput_FlagAvailableGlobally(t *testing.T) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "devtui_test_json")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build devtui: %v", err)
	}
	defer os.Remove("../devtui_test_json")

	binary := "../devtui_test_json"

	// Test that the --json flag is available on various commands
	commands := [][]string{
		{"count", "--help"},
		{"version", "--help"},
		{"iban", "--help"},
		{"urls", "--help"},
	}

	for _, args := range commands {
		t.Run("help_"+args[0], func(t *testing.T) {
			cmd := exec.Command(binary, args...)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
			}

			output := stdout.String()
			if !strings.Contains(output, "-j, --json") && !strings.Contains(output, "--json") {
				t.Errorf("Command %s should show --json flag in help", args[0])
			}
		})
	}
}

func TestJSONOutput_InvalidInput(t *testing.T) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "devtui_test_json")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build devtui: %v", err)
	}
	defer os.Remove("../devtui_test_json")

	binary := "../devtui_test_json"

	// Test count with empty input
	cmd := exec.Command(binary, "count", "--json", "")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())

	var result map[string]int
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
	}

	if result["characters"] != 0 {
		t.Errorf("Expected characters=0 for empty input, got %d", result["characters"])
	}
}
