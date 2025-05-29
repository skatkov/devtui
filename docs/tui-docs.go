package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type TUIModule struct {
	Name        string
	Title       string
	Description string
	KeyBindings []KeyBinding
	FilePath    string
}

type KeyBinding struct {
	Key         string
	Description string
}

func extractTUIModules() ([]TUIModule, error) {
	var modules []TUIModule

	tuiDir := "tui"
	if _, err := os.Stat("tui"); os.IsNotExist(err) {
		tuiDir = "../tui"
	}

	entries, err := os.ReadDir(tuiDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modulePath := filepath.Join(tuiDir, entry.Name())
		module, err := extractModuleInfo(modulePath, entry.Name())
		if err != nil {
			log.Printf("Warning: Failed to extract info for %s: %v", entry.Name(), err)
			continue
		}

		if module != nil {
			modules = append(modules, *module)
		}
	}

	// Sort modules alphabetically by title
	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Title < modules[j].Title
	})

	return modules, nil
}

func extractModuleInfo(modulePath, moduleName string) (*TUIModule, error) {
	files, err := os.ReadDir(modulePath)
	if err != nil {
		return nil, err
	}

	var module *TUIModule

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(modulePath, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		fileContent := string(content)

		// Extract title
		titleRegex := regexp.MustCompile(`const Title = "([^"]+)"`)
		if matches := titleRegex.FindStringSubmatch(fileContent); len(matches) > 1 {
			if module == nil {
				module = &TUIModule{
					Name:     moduleName,
					FilePath: filePath,
				}
			}
			module.Title = matches[1]
		}

		// Extract key bindings from helpView function
		keyBindings := extractKeyBindings(fileContent)
		if len(keyBindings) > 0 {
			if module == nil {
				module = &TUIModule{
					Name:     moduleName,
					FilePath: filePath,
				}
			}
			module.KeyBindings = append(module.KeyBindings, keyBindings...)
		}

		// Extract description from comments or docstrings
		description := extractDescription(module)
		if description != "" {
			if module == nil {
				module = &TUIModule{
					Name:     moduleName,
					FilePath: filePath,
				}
			}
			module.Description = description
		}
	}

	return module, nil
}

func extractKeyBindings(content string) []KeyBinding {
	var bindings []KeyBinding

	// Look for col1 array definitions in helpView functions
	col1Regex := regexp.MustCompile(`col1\s*:=\s*\[\]string\{([^}]+)\}`)
	matches := col1Regex.FindStringSubmatch(content)
	if len(matches) < 2 {
		return bindings
	}

	col1Content := matches[1]

	// Extract individual string entries like "c              copy text"
	entryRegex := regexp.MustCompile(`"([^"]+)"`)
	entryMatches := entryRegex.FindAllStringSubmatch(col1Content, -1)

	for _, match := range entryMatches {
		if len(match) >= 2 {
			entry := strings.TrimSpace(match[1])

			// Split on multiple spaces to separate key from description
			parts := regexp.MustCompile(`\s{2,}`).Split(entry, 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				desc := strings.TrimSpace(parts[1])

				// Skip navigation keys and empty descriptions
				if len(desc) > 3 && !strings.Contains(key, "↑") && !strings.Contains(key, "↓") && key != "k" && key != "j" {
					bindings = append(bindings, KeyBinding{
						Key:         key,
						Description: desc,
					})
				}
			}
		}
	}

	return bindings
}

func extractDescription(module *TUIModule) string {
	// Generate description based on title
	if module != nil && module.Title != "" {
		return generateDescriptionFromTitle(module.Title)
	}

	return ""
}

func generateDescriptionFromTitle(title string) string {
	title = strings.ToLower(title)

	if strings.Contains(title, "formatter") {
		return "A text formatting tool that prettifies and standardizes code or data format."
	}
	if strings.Contains(title, "converter") {
		return "A data conversion tool that transforms content between different formats."
	}
	if strings.Contains(title, "encoder") {
		return "An encoding tool that transforms text into encoded format."
	}
	if strings.Contains(title, "decoder") {
		return "A decoding tool that transforms encoded content back to readable format."
	}
	if strings.Contains(title, "generator") {
		return "A utility tool that generates new content based on specified parameters."
	}
	if strings.Contains(title, "parser") {
		return "A parsing tool that analyzes and interprets structured content."
	}
	if strings.Contains(title, "renderer") {
		return "A rendering tool that displays content in a formatted view."
	}

	return "An interactive TUI tool for text processing and manipulation."
}

func generateTUIDocumentation(modules []TUIModule) error {
	sitePath := "site/tui"
	if _, err := os.Stat("site"); os.IsNotExist(err) {
		sitePath = "../site/tui"
	}

	// Create tui directory if it doesn't exist
	if err := os.MkdirAll(sitePath, 0o755); err != nil {
		return err
	}

	// Generate index file
	if err := generateIndexFile(sitePath); err != nil {
		return err
	}

	// Generate individual module files
	for _, module := range modules {
		if err := generateModuleFile(sitePath, module); err != nil {
			return err
		}
	}

	return nil
}

func generateIndexFile(sitePath string) error {
	content := `---
title: TUI
nav_order: 4
---

# TUI (Terminal User Interface)

## Getting Started

To access the TUI interface, simply run:

` + "```bash\ndevtui\n```" + `

This will open the main menu where you can select any of the available tools using arrow keys and Enter.

## Common Key Bindings

Most TUI tools share these common key bindings:

- **q/Ctrl+C** - Quit and return to main menu
- **?** - Toggle help view
- **c** - Copy output to clipboard
- **v** - Paste content from clipboard
- **e** - Edit content in external editor
- **↑/k** - Navigate up
- **↓/j** - Navigate down
`

	return os.WriteFile(filepath.Join(sitePath, "index.md"), []byte(content), 0o644)
}

func generateModuleFile(sitePath string, module TUIModule) error {
	content := fmt.Sprintf(`---
title: %s
parent: TUI
---

# %s

## Usage

1. Run `+"`devtui`"+` to open the main menu
2. Select "`+module.Title+`" from the list
3. Use the key bindings below to interact with the tool
4. Press `+"`q`"+` or `+"`Ctrl+C`"+` to return to the main menu

## Key Bindings

`, module.Title, module.Title)

	if len(module.KeyBindings) > 0 {
		content += "| Key | Action |\n|-----|--------|\n"
		for _, binding := range module.KeyBindings {
			content += fmt.Sprintf("| `%s` | %s |\n", binding.Key, binding.Description)
		}
	} else {
		content += "Standard key bindings apply (see main TUI documentation).\n"
	}

	content += `


`

	filename := module.Name + ".md"
	return os.WriteFile(filepath.Join(sitePath, filename), []byte(content), 0o644)
}

func GenerateTUIDocumentation() error {
	modules, err := extractTUIModules()
	if err != nil {
		return fmt.Errorf("failed to extract TUI modules: %w", err)
	}

	if len(modules) == 0 {
		fmt.Println("No TUI modules found")
		return nil
	}

	if err := generateTUIDocumentation(modules); err != nil {
		return fmt.Errorf("failed to generate TUI documentation: %w", err)
	}

	fmt.Printf("Generated documentation for %d TUI modules\n", len(modules))
	return nil
}
