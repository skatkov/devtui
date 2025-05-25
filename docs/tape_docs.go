package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type TapeModule struct {
	ID          string
	Title       string
	Description string
	TapeFile    string
	PngFile     string
	Actions     []TapeAction
}

type TapeAction struct {
	Action string
	Value  string
	Sleep  string
}

func extractTapeModules() ([]TapeModule, error) {
	var modules []TapeModule

	// Read the list.go file to extract module information
	listPath := "tui/root/list.go"
	if _, err := os.Stat("tui"); os.IsNotExist(err) {
		listPath = "../tui/root/list.go"
	}

	content, err := os.ReadFile(listPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read list.go: %w", err)
	}

	modules = extractModulesFromList(string(content))
	return modules, nil
}

func extractModulesFromList(content string) []TapeModule {
	var modules []TapeModule

	// Extract menu options from getMenuOptions function
	optionRegex := regexp.MustCompile(`{\s*id:\s*"([^"]+)",\s*title:\s*([^,]+),`)
	matches := optionRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			id := match[1]
			titleVar := strings.TrimSpace(match[2])
			
			// Extract actual title value
			title := extractTitleFromVariable(content, titleVar)
			if title == "" {
				title = strings.Title(strings.ReplaceAll(id, "-", " "))
			}

			modules = append(modules, TapeModule{
				ID:       id,
				Title:    title,
				TapeFile: fmt.Sprintf("demo-%s.tape", id),
				PngFile:  fmt.Sprintf("%s.png", id),
				Actions:  generateTapeActions(id),
			})
		}
	}

	// Sort modules alphabetically by title
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Title < modules[j].Title
	})

	return modules
}

func extractTitleFromVariable(content, variable string) string {
	// Remove package prefix if present (e.g., "js.Title" -> "Title")
	parts := strings.Split(variable, ".")
	titleVar := parts[len(parts)-1]
	
	// Look for const Title = "..." pattern
	titleRegex := regexp.MustCompile(fmt.Sprintf(`const\s+%s\s*=\s*"([^"]+)"`, titleVar))
	matches := titleRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	
	return ""
}

func generateTapeActions(moduleID string) []TapeAction {
	actions := []TapeAction{
		{Action: "Set", Value: "Shell zsh"},
		{Action: "Set", Value: "WindowBar ColorfulRight"},
		{Action: "Type", Value: "cd .. && go run .", Sleep: "500ms"},
		{Action: "Enter", Sleep: "2.5s"},
	}

	// Use search to find the module
	// Press / to start search
	actions = append(actions, TapeAction{Action: "Type", Value: "/", Sleep: "300ms"})
	
	// Type the module name to search
	searchText := getSearchText(moduleID)
	actions = append(actions, TapeAction{Action: "Type", Value: searchText, Sleep: "500ms"})
	
	// Select the top result and enter the module
	actions = append(actions, TapeAction{Action: "Enter", Sleep: "500ms"})
	actions = append(actions, TapeAction{Action: "Enter", Sleep: "1s"})

	// Take screenshot immediately after entering the tool
	actions = append(actions, TapeAction{Action: "Sleep", Value: "500ms"})
	actions = append(actions, TapeAction{Action: "Screenshot"})

	// Exit back to menu
	actions = append(actions, TapeAction{Action: "Escape", Sleep: "500ms"})

	return actions
}

func getSearchText(moduleID string) string {
	// Generate search text that uniquely identifies the module
	switch moduleID {
	case "base64-decoder":
		return "base64 decoder"
	case "base64-encoder":
		return "base64 encoder"
	case "cron":
		return "cron"
	case "css":
		return "css"
	case "csv2md":
		return "csv2md"
	case "csvjson":
		return "csvjson"
	case "graphql-query":
		return "graphql"
	case "html":
		return "html"
	case "json":
		return "json"
	case "jsonstruct":
		return "jsonstruct"
	case "jsontoml":
		return "jsontoml"
	case "markdown":
		return "markdown"
	case "numbers":
		return "numbers"
	case "toml":
		return "toml"
	case "tomljson":
		return "tomljson"
	case "tsv2md":
		return "tsv2md"
	case "uuiddecode":
		return "uuid decode"
	case "uuidgenerate":
		return "uuid generate"
	case "xml":
		return "xml"
	case "yaml":
		return "yaml"
	case "yamlstruct":
		return "yamlstruct"
	default:
		// Fallback: use the module ID
		return moduleID
	}
}



func generateTapeFiles(modules []TapeModule) error {
	tapeDir := "tapes"
	if err := os.MkdirAll(tapeDir, 0755); err != nil {
		return fmt.Errorf("failed to create tapes directory: %w", err)
	}

	for _, module := range modules {
		if err := generateTapeFile(tapeDir, module); err != nil {
			return fmt.Errorf("failed to generate tape for %s: %w", module.ID, err)
		}
	}

	return nil
}

func generateTapeFile(tapeDir string, module TapeModule) error {
	var content strings.Builder

	// Add actions
	for _, action := range module.Actions {
		switch action.Action {
		case "Set":
			content.WriteString(fmt.Sprintf("Set %s\n", action.Value))
		case "Type":
			content.WriteString(fmt.Sprintf("Type \"%s\"\n", action.Value))
		case "Enter":
			content.WriteString("Enter\n")
		case "Down":
			content.WriteString("Down\n")
		case "Screenshot":
			content.WriteString(fmt.Sprintf("Screenshot ../site/assets/img/tui/%s\n", module.PngFile))
		case "Escape":
			content.WriteString("Escape\n")
		case "Sleep":
			content.WriteString(fmt.Sprintf("Sleep %s\n", action.Value))
		}

		if action.Sleep != "" {
			content.WriteString(fmt.Sprintf("Sleep %s\n", action.Sleep))
		}
	}

	filePath := filepath.Join(tapeDir, module.TapeFile)
	return os.WriteFile(filePath, []byte(content.String()), 0644)
}

func updateMarkdownFiles(modules []TapeModule) error {
	sitePath := "site/tui"
	if _, err := os.Stat("site"); os.IsNotExist(err) {
		sitePath = "../site/tui"
	}

	for _, module := range modules {
		if err := updateModuleMarkdownFile(sitePath, module); err != nil {
			return fmt.Errorf("failed to update markdown for %s: %w", module.ID, err)
		}
	}

	return nil
}

func getMarkdownFilename(moduleID string) string {
	// Map module IDs to correct markdown filenames
	switch moduleID {
	case "uuiddecode":
		return "uuid-decode.md"
	case "uuidgenerate":
		return "uuid-generate.md"
	default:
		return fmt.Sprintf("%s.md", moduleID)
	}
}

func updateModuleMarkdownFile(sitePath string, module TapeModule) error {
	// Map module IDs to correct markdown filenames
	filename := getMarkdownFilename(module.ID)
	filePath := filepath.Join(sitePath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Warning: Markdown file %s does not exist, skipping\n", filename)
		return nil
	}

	// Read existing file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read existing markdown file: %w", err)
	}

	contentStr := string(content)

	// Remove all existing screenshots first
	screenshotRegex := regexp.MustCompile(`!\[Screenshot\]\([^)]+\)\n?`)
	contentStr = screenshotRegex.ReplaceAllString(contentStr, "")

	// Find the end of front matter and the main title
	frontMatterRegex := regexp.MustCompile(`(---\n[\s\S]*?\n---\n\n)(# [^\n]+\n)`)
	screenshot := fmt.Sprintf("${1}${2}\n![Screenshot](/assets/img/tui/%s)\n", module.PngFile)
	
	if frontMatterRegex.MatchString(contentStr) {
		contentStr = frontMatterRegex.ReplaceAllString(contentStr, screenshot)
	} else {
		// Fallback to simple title matching if no front matter
		titleRegex := regexp.MustCompile(`^(# [^\n]+\n)`)
		screenshotFallback := fmt.Sprintf("${1}\n![Screenshot](/assets/img/tui/%s)\n", module.PngFile)
		contentStr = titleRegex.ReplaceAllString(contentStr, screenshotFallback)
	}

	return os.WriteFile(filePath, []byte(contentStr), 0644)
}

func createAssetsDirectory() error {
	assetsPath := "site/assets/img/tui"
	if _, err := os.Stat("site"); os.IsNotExist(err) {
		assetsPath = "../site/assets/img/tui"
	}

	return os.MkdirAll(assetsPath, 0755)
}

func generateRunScript(modules []TapeModule) error {
	var script strings.Builder
	
	script.WriteString("#!/bin/bash\n\n")
	script.WriteString("# Auto-generated script to run all tape demos\n")
	script.WriteString("# Make sure you have vhs installed: go install github.com/charmbracelet/vhs@latest\n\n")
	
	script.WriteString("set -e\n\n")
	script.WriteString("# Change to docs directory if not already there\n")
	script.WriteString("cd \"$(dirname \"$0\")\"\n\n")
	script.WriteString("echo \"Generating screenshots for all TUI modules...\"\n\n")
	
	script.WriteString("echo \"Generating main menu demo...\"\n")
	script.WriteString("vhs tapes/demo-main.tape\n")
	script.WriteString("sleep 1\n\n")

	for _, module := range modules {
		script.WriteString(fmt.Sprintf("echo \"Generating screenshot for %s...\"\n", module.Title))
		script.WriteString(fmt.Sprintf("vhs tapes/%s\n", module.TapeFile))
		script.WriteString("sleep 1\n\n")
	}

	script.WriteString("echo \"All screenshots generated successfully!\"\n")

	return os.WriteFile("run_demos.sh", []byte(script.String()), 0755)
}

func generateMasterDemoTape() error {
	tapeDir := "tapes"
	if err := os.MkdirAll(tapeDir, 0755); err != nil {
		return fmt.Errorf("failed to create tapes directory: %w", err)
	}

	var content strings.Builder
	
	// Master demo showing main menu navigation
	content.WriteString("Set Shell zsh\n")
	content.WriteString("Set WindowBar ColorfulRight\n")
	content.WriteString("Type \"cd .. && go run .\"\n")
	content.WriteString("Sleep 500ms\n")
	content.WriteString("Enter\n")
	content.WriteString("Sleep 2.5s\n")
	
	// Show navigation through a few items
	content.WriteString("Down\n")
	content.WriteString("Sleep 300ms\n")
	content.WriteString("Down\n")
	content.WriteString("Sleep 300ms\n")
	content.WriteString("Down\n")
	content.WriteString("Sleep 300ms\n")
	content.WriteString("Up\n")
	content.WriteString("Sleep 300ms\n")
	content.WriteString("Up\n")
	content.WriteString("Sleep 500ms\n")
	
	// Show search functionality
	content.WriteString("Type \"/\"\n")
	content.WriteString("Sleep 300ms\n")
	content.WriteString("Type \"json\"\n")
	content.WriteString("Sleep 800ms\n")
	content.WriteString("Screenshot ../site/assets/img/devtui-main.png\n")
	content.WriteString("Escape\n")
	content.WriteString("Sleep 500ms\n")
	content.WriteString("Escape\n")
	content.WriteString("Sleep 500ms\n")

	filePath := filepath.Join(tapeDir, "demo-main.tape")
	return os.WriteFile(filePath, []byte(content.String()), 0644)
}

func GenerateTapeDocumentation() error {
	// Create assets directory
	if err := createAssetsDirectory(); err != nil {
		return fmt.Errorf("failed to create assets directory: %w", err)
	}

	// Extract modules
	modules, err := extractTapeModules()
	if err != nil {
		return fmt.Errorf("failed to extract tape modules: %w", err)
	}

	if len(modules) == 0 {
		fmt.Println("No TUI modules found")
		return nil
	}

	// Generate master demo tape
	if err := generateMasterDemoTape(); err != nil {
		return fmt.Errorf("failed to generate master demo tape: %w", err)
	}

	// Generate tape files
	if err := generateTapeFiles(modules); err != nil {
		return fmt.Errorf("failed to generate tape files: %w", err)
	}

	// Update markdown files with screenshots
	if err := updateMarkdownFiles(modules); err != nil {
		return fmt.Errorf("failed to update markdown files: %w", err)
	}

	// Generate run script
	if err := generateRunScript(modules); err != nil {
		return fmt.Errorf("failed to generate run script: %w", err)
	}

	fmt.Printf("Generated tape documentation for %d TUI modules\n", len(modules))
	fmt.Println("Run './run_demos.sh' to generate all screenshots")
	
	return nil
}