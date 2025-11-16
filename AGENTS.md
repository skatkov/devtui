# AGENTS.md - DevTUI Development Guide

This file helps AI agents (and developers) understand how to work effectively in the DevTUI codebase.

## Project Overview

**DevTUI** is a Swiss-army knife terminal app for developers, providing both CLI and TUI (Terminal User Interface) modes for common developer utilities like JSON formatting, base64 encoding/decoding, UUID generation, and more.

- **Language**: Go 1.25.3
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **TUI Framework**: Bubble Tea (github.com/charmbracelet/bubbletea)
- **Build Tool**: goreleaser
- **Testing**: Standard Go testing (`testing` package)
- **Linting**: golangci-lint with formatters gofumpt and goimports

## Essential Commands

### Build & Run
```bash
# Build the project
go build -v ./...

# Run the TUI (main mode)
go run main.go

# Run a specific CLI command
go run main.go jsonfmt < testdata/example.json
go run main.go base64 "hello world"
go run main.go --help
```

### Testing
```bash
# Run all tests with race detection (this is what CI uses)
go test -v -failfast -race ./...

# Run specific package tests
go test -v ./cmd/...
go test -v ./internal/ui/...

# Run specific test
go test -v -run TestJson2toonCmd ./cmd/
```

### Linting & Formatting
```bash
# Run golangci-lint (what CI uses)
golangci-lint run

# The linters will auto-format with gofumpt and goimports
# Enabled linters: bodyclose, goconst, gomoddirectives, goprintffuncname, 
# misspell, nakedret, nilerr, noctx, nolintlint, prealloc, rowserrcheck, 
# sqlclosecheck, tparallel, unconvert, unparam, whitespace, perfsprint
```

### Go Module Management
```bash
# Tidy dependencies (goreleaser runs this before builds)
go mod tidy

# Generate code (goreleaser runs this before builds)
go generate ./...
```

### Documentation Generation
```bash
# Generate both CLI and TUI documentation (for website)
cd docs && go run *.go

# Generate only CLI docs
cd docs && go run cli-docs.go docs.go

# Generate only TUI docs  
cd docs && go run tui-docs.go docs.go
```

### Release
```bash
# Build release with goreleaser (creates binaries for linux/darwin/windows)
goreleaser release --snapshot --clean

# The project auto-releases via GitHub Actions on tags
```

## Project Structure

```
devtui/
├── main.go                    # Entry point - just calls cmd.Execute()
├── cmd/                       # CLI command implementations (Cobra)
│   ├── root.go               # Root command, launches TUI by default
│   ├── base64.go             # CLI commands like base64, jsonfmt, etc.
│   ├── jsonfmt.go
│   ├── json2toon.go
│   ├── jsonrepair.go
│   ├── *_test.go             # CLI command tests
│   └── ...                   # ~12 CLI commands total
├── tui/                      # TUI module implementations (Bubble Tea)
│   ├── root/                 # Root TUI screen with menu
│   │   ├── root.go          # Root model managing screen switching
│   │   └── list.go          # Menu list with usage tracking
│   ├── json/                # Each TUI module in its own package
│   ├── yaml/
│   ├── base64-encoder/
│   ├── json2toon/
│   └── ...                  # ~27 TUI modules total
├── internal/                # Shared internal packages
│   ├── ui/                  # TUI UI components & common types
│   │   ├── ui.go           # CommonModel, ReturnToListMsg, constants
│   │   ├── base_pager_model.go  # BasePagerModel (base for most TUI screens)
│   │   └── style.go        # Lipgloss styles
│   ├── editor/             # External editor integration
│   ├── base64/             # Base64 encode/decode logic
│   ├── csv2md/             # CSV to markdown conversion
│   └── textanalyzer/       # Text analysis utilities
├── docs/                   # Documentation generators
├── testdata/               # Test fixtures (JSON, CSV, TOML, etc.)
├── site/                   # Jekyll website (generated docs)
├── .golangci.yml          # Linter config
├── .goreleaser.yaml       # Release config
└── .github/workflows/     # CI/CD
    ├── build.yml          # Build & test on push/PR
    ├── lint.yml           # Linting with golangci-lint
    └── jekyll-gh-pages.yml # Deploy docs to GitHub Pages
```

## Code Architecture

### Dual Interface Pattern

DevTUI provides **both CLI and TUI** for most utilities:

1. **CLI Commands** (`cmd/`): For piping/scripting
   - Use Cobra framework
   - Read from stdin or args
   - Write to stdout
   - Example: `devtui jsonfmt < file.json > output.json`

2. **TUI Modules** (`tui/`): For interactive use
   - Use Bubble Tea framework
   - Full-screen interactive UIs
   - Example: Run `devtui` (no args) to see menu

**Not all features have both interfaces yet** - some are CLI-only or TUI-only.

### TUI Architecture

**Root Screen** (`tui/root/`):
- `root.go`: Manages current view and screen switching
- `list.go`: Main menu with fuzzy search and usage tracking
- Menu items are sorted by usage frequency (stored in `~/.config/devtui/usage_stats.json`)

**TUI Modules** (`tui/*/`):
- Each tool is a separate package (e.g., `tui/json`, `tui/yaml`)
- Most inherit from `BasePagerModel` (`internal/ui/base_pager_model.go`)
- Common pattern: Title constant, Model struct, Init/Update/View methods

**BasePagerModel** provides:
- Common key bindings (q/ctrl+c to quit, esc to return to menu, ? for help)
- Copy to clipboard (press 'c')
- Paste from clipboard (press 'v')
- Edit in external editor (press 'e')
- Status message handling
- Viewport scrolling

**Common Types** (`internal/ui/ui.go`):
- `CommonModel`: Shared state (width, height, styles, last selected item)
- `ReturnToListMsg`: Message to return to main menu
- `PagerState`: Browse, StatusMessage, ErrorMessage states

### CLI Architecture

**Root Command** (`cmd/root.go`):
- No args = launch TUI
- Subcommands = run CLI utilities

**CLI Command Pattern**:
```go
var myCmd = &cobra.Command{
    Use:   "mycommand [args]",
    Short: "Short description",
    Long:  `Long description with examples`,
    Args:  cobra.NoArgs,  // Or cobra.ExactArgs(1), etc.
    PreRunE: func(cmd *cobra.Command, args []string) error {
        // Validate flags/args
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        // Read from cmd.InOrStdin()
        // Process data
        // Write to cmd.OutOrStdout()
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
    myCmd.Flags().StringVarP(&flagVar, "flag", "f", "default", "description")
}
```

## Naming Conventions

### Files
- CLI commands: `cmd/<name>.go` (e.g., `cmd/jsonfmt.go`)
- CLI tests: `cmd/<name>_test.go` (e.g., `cmd/jsonfmt_test.go`)
- TUI modules: `tui/<name>/main.go` or `tui/<name>/<name>.go`
- TUI tests: `tui/<name>/main_test.go`
- Internal packages: `internal/<name>/<name>.go`

### Variables
- Cobra commands: `<name>Cmd` (e.g., `jsonfmtCmd`, `base64Cmd`)
- Command flags: `<command><FlagName>` (e.g., `base64Decode`, `json2toonIndent`)
- Models: `<Name>Model` (e.g., `JsonModel`, `JsonToonModel`)
- Bubble Tea models implement: `Init()`, `Update(tea.Msg)`, `View() string`

### Constants
- TUI module titles: `const Title = "Display Name"` (exported from each TUI package)
- Help text convention: Use column layout with keyboard shortcuts

### Package Names
- CLI commands: `package cmd`
- TUI modules: `package <modulename>` (e.g., `package json`, `package json2toon`)
- Hyphenated directories become package names without hyphens where sensible

## Common Patterns

### Reading Input in CLI
```go
// Read from stdin or cmd.InOrStdin()
data, err := io.ReadAll(cmd.InOrStdin())
if err != nil {
    return fmt.Errorf("error reading input: %w", err)
}

// For commands that also accept args
if len(args) > 0 {
    data = []byte(args[0])
} else {
    data, err = io.ReadAll(cmd.InOrStdin())
}
```

### Writing Output in CLI
```go
// Always write to cmd.OutOrStdout() for testability
_, err = fmt.Fprintln(cmd.OutOrStdout(), result)
if err != nil {
    return err
}
```

### TUI Module Structure
```go
package mymodule

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/skatkov/devtui/internal/ui"
)

const Title = "My Module Name"

type MyModel struct {
    ui.BasePagerModel
    // Additional fields
}

func NewMyModel(common *ui.CommonModel) MyModel {
    return MyModel{
        BasePagerModel: ui.NewBasePagerModel(common, Title),
    }
}

func (m MyModel) Init() tea.Cmd {
    return m.BasePagerModel.Init()
}

func (m MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // First check common keys
        if cmd, handled := m.HandleCommonKeys(msg); handled {
            return m, cmd
        }
        // Then handle module-specific keys
        switch msg.String() {
        case "e":
            return m, editor.OpenEditor(m.Content, "json")
        case "v":
            // Paste from clipboard logic
        }
    case ui.StatusMessageTimeoutMsg:
        m.State = ui.PagerStateBrowse
    case editor.EditorFinishedMsg:
        // Handle editor close
    case tea.WindowSizeMsg:
        cmd = m.HandleWindowSizeMsg(msg)
        cmds = append(cmds, cmd)
    }
    
    m.Viewport, cmd = m.Viewport.Update(msg)
    cmds = append(cmds, cmd)
    return m, tea.Batch(cmds...)
}

func (m MyModel) View() string {
    var b strings.Builder
    fmt.Fprint(&b, m.Viewport.View()+"\n")
    fmt.Fprint(&b, m.StatusBarView())
    if m.ShowHelp {
        fmt.Fprint(&b, "\n"+m.helpView())
    }
    return b.String()
}

func (m *MyModel) SetContent(content string) error {
    m.Content = content
    m.FormattedContent = FormatContent(content)
    m.Viewport.SetContent(m.FormattedContent)
    return nil
}

func (m MyModel) helpView() string {
    // Return formatted help text
}
```

### Adding a TUI Module to Menu
Register in `tui/root/list.go` in `getMenuOptions()`:
```go
{
    id:    "mymodule",  // Used for usage tracking
    title: mymodule.Title,
    model: func() tea.Model { return mymodule.NewMyModel(common) },
},
```

### Test Patterns

**CLI Tests**:
```go
func TestMyCmd(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        args        []string
        wantContain string
        wantErr     bool
        description string
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := GetRootCmd()
            buf := new(bytes.Buffer)
            cmd.SetOut(buf)
            cmd.SetErr(buf)
            cmd.SetIn(strings.NewReader(tt.input))
            
            args := []string{"mycommand"}
            args = append(args, tt.args...)
            cmd.SetArgs(args)
            
            err := cmd.Execute()
            // Assert error and output
        })
    }
}
```

**TUI Tests**: Usually test the conversion logic separately (e.g., `tui/json2toon/main_test.go`)

### Error Handling
- Wrap errors with context: `fmt.Errorf("failed to parse JSON: %w", err)`
- CLI commands return errors (RunE)
- TUI shows errors via `m.ShowErrorMessage(err.Error())`

## Important Details

### Usage Statistics
- TUI menu tracks usage frequency in `~/.config/devtui/usage_stats.json`
- Uses `github.com/adrg/xdg` for proper XDG config directories
- Menu items auto-sort by frequency on each return to menu

### External Editor Integration
- Press 'e' in TUI modules to edit in external editor
- Uses `$EDITOR` environment variable
- Creates temp file with appropriate extension (`.json`, `.yaml`, etc.)
- Handled by `internal/editor/editor.go`

### Clipboard Integration
- Press 'c' to copy formatted output
- Press 'v' to paste input
- Uses `github.com/tiagomelo/go-clipboard`
- Works cross-platform

### Key Bindings (Standard Across TUI)
- `q` or `ctrl+c`: Quit
- `esc`: Return to main menu
- `?`: Toggle help
- `c`: Copy formatted content
- `v`: Paste content from clipboard
- `e`: Edit in external editor (where applicable)
- `k`/`↑` or `j`/`↓`: Scroll up/down
- `u`/`d`: Half page up/down
- `b`/`f` or `pgup`/`pgdn`: Full page up/down

Additional module-specific keys (examples):
- `i`: Toggle indent (json2toon)
- `l`: Toggle length marker (json2toon)

### Syntax Highlighting
- Uses `github.com/alecthomas/chroma` for code highlighting
- Theme: "nord"
- Applied in TUI modules via `quick.Highlight()`

### Dependencies
Key external libraries:
- **TUI**: bubbletea, bubbles, lipgloss, glamour
- **CLI**: cobra
- **Utilities**: 
  - go-yaml (YAML parsing)
  - go-toml/v2 (TOML parsing)
  - uuid (UUID generation/decoding)
  - go-xmlfmt (XML formatting)
  - csstool (CSS minification)
  - json-repair (LLM JSON repair)
  - toon (JSON to TOON conversion)
  - vektah/gqlparser (GraphQL query formatting)
  - xurls (URL extraction)

## Gotchas & Non-Obvious Patterns

### 1. Dual TUI/CLI Implementation
Some utilities have both CLI command (`cmd/`) and TUI module (`tui/`). The shared logic is often in the TUI package and imported by CLI:

```go
// cmd/jsonfmt.go imports from tui/json
import "github.com/skatkov/devtui/tui/json"
result := json.FormatJSON(string(data))
```

### 2. BasePagerModel Must Be Embedded
TUI models should embed `ui.BasePagerModel`, not use it as a pointer field:
```go
// Correct
type MyModel struct {
    ui.BasePagerModel
}

// Wrong - methods won't be available
type MyModel struct {
    base *ui.BasePagerModel
}
```

### 3. Viewport Content Setting
When setting viewport content, use `Viewport.SetContent()`, not direct assignment:
```go
// Correct
m.Viewport.SetContent(formattedContent)

// Wrong
m.Viewport.Content = formattedContent
```

### 4. Window Size Handling
Always handle `tea.WindowSizeMsg` and update both `Common` and viewport:
```go
case tea.WindowSizeMsg:
    cmd = m.HandleWindowSizeMsg(msg)
    cmds = append(cmds, cmd)
```

### 5. Command Flag Variables
CLI command flags use package-level variables prefixed with command name to avoid conflicts:
```go
var (
    json2toonIndent       int    // Not just "indent"
    json2toonLengthMarker string // Not just "lengthMarker"
)
```

### 6. Menu Registration Order
Menu items are sorted by usage, but initial order in `getMenuOptions()` doesn't matter - frequency overrides it.

### 7. Testing CLI with GetRootCmd()
CLI tests must use `GetRootCmd()` to get a fresh command instance:
```go
cmd := GetRootCmd()  // Not rootCmd directly
```

### 8. Testdata Files
Test fixtures live in `testdata/` at repo root. Reference them as `testdata/example.json`.

### 9. Go Generate
The project has `go generate` hooks (run by goreleaser), though none are currently obvious in the code. Safe to run: `go generate ./...`

### 10. Documentation Is Auto-Generated
Don't manually edit files in `site/cli/` or `site/tui/` - they're generated by `docs/*.go`. Edit the generators instead.

### 11. No Makefile
The project doesn't use Make. Use `go` commands directly or `goreleaser`.

### 12. Version Injection
Version info is injected at build time via ldflags (see `.goreleaser.yaml`):
```go
// cmd/version.go references these injected vars
var version string
var commit string  
var date string
```

### 13. Error Messages in TUI
Use `ShowErrorMessage()` and `ShowStatusMessage()` from BasePagerModel - they auto-timeout after 3 seconds:
```go
return m, m.ShowErrorMessage("Something went wrong")
return m, m.ShowStatusMessage("Copied!")
```

### 14. Race Detection in Tests
CI runs tests with `-race` flag. Avoid shared mutable state without proper locking.

### 15. golangci-lint Excludes Generated Code
The linter config has `generated: lax` to skip generated files. Mark generated files with standard comment:
```go
// Code generated by <tool>. DO NOT EDIT.
```

## Adding New Features

### Adding a New CLI Command

1. Create `cmd/mycommand.go`:
```go
package cmd

var mycommandCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(mycommandCmd)
}
```

2. Add tests in `cmd/mycommand_test.go`

3. Test: `go test -v ./cmd/`

4. Build: `go build`

5. Document: Run `cd docs && go run *.go` to regenerate CLI docs

### Adding a New TUI Module

1. Create package `tui/mymodule/`

2. Create `tui/mymodule/main.go` (or `mymodule.go`):
```go
package mymodule

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/skatkov/devtui/internal/ui"
)

const Title = "My Module"

type MyModel struct {
    ui.BasePagerModel
}

func NewMyModel(common *ui.CommonModel) MyModel {
    return MyModel{
        BasePagerModel: ui.NewBasePagerModel(common, Title),
    }
}

// Implement Init, Update, View, SetContent, helpView
```

3. Register in `tui/root/list.go`:
```go
import "github.com/skatkov/devtui/tui/mymodule"

// In getMenuOptions():
{
    id:    "mymodule",
    title: mymodule.Title,
    model: func() tea.Model { return mymodule.NewMyModel(common) },
},
```

4. Test: `go run main.go` and select your module

5. Document: Run `cd docs && go run *.go` to regenerate TUI docs

### Adding Shared Utility Code

Place in `internal/<packagename>/` if it's DevTUI-specific logic not suitable for external use.

## CI/CD Workflows

### On Push / Pull Request
1. **Build Job** (`.github/workflows/build.yml`):
   - Runs on Ubuntu and macOS
   - `go build -v ./...`
   - `go vet ./...`
   - `go test -v -failfast -race ./...`

2. **Lint Job** (`.github/workflows/lint.yml`):
   - Runs `golangci-lint` via GitHub Action
   - Uses config from `.golangci.yml`

3. **Dependabot Auto-Merge**:
   - If build passes and PR author is dependabot, auto-approve and merge

### On Tag Push
- goreleaser runs automatically (via GitHub Actions, config in `.goreleaser.yaml`)
- Builds binaries for Linux, macOS, Windows
- Creates GitHub release
- Updates Homebrew tap (skatkov/homebrew-tap)

### Jekyll Site Deployment
- `site/` directory is Jekyll-based documentation
- Deployed to GitHub Pages via `.github/workflows/jekyll-gh-pages.yml`
- Documentation is auto-generated from code

## Common Development Tasks

### Run TUI locally
```bash
go run main.go
```

### Test a CLI command
```bash
echo '{"test": "value"}' | go run main.go jsonfmt
go run main.go base64 "hello" 
```

### Run specific tests
```bash
go test -v -run TestJson2toon ./cmd/
```

### Check linting
```bash
golangci-lint run
```

### Generate docs
```bash
cd docs && go run *.go
```

### Test against testdata
```bash
go run main.go jsonfmt < testdata/example.json
go run main.go base64 testdata/sample.txt
```

### View test coverage
```bash
go test -coverprofile=coverage.txt ./...
go tool cover -html=coverage.txt
```

## Style Guide

### Go Code Style
- Follow standard Go conventions
- Use gofumpt and goimports (enforced by golangci-lint)
- Keep functions focused and small
- Export only what's necessary

### Error Messages
- Lowercase, no trailing punctuation: `"failed to parse JSON"`
- Use `fmt.Errorf()` with `%w` for wrapping: `fmt.Errorf("context: %w", err)`
- Provide context: not just "invalid input", but "invalid indent: -1 (must be non-negative)"

### Comments
- Public functions/types must have doc comments
- Doc comments start with the name: `// FormatJSON formats the given JSON string.`
- Use `//` not `/* */` except for package docs

### Variable Names
- Short names for short scopes: `i`, `err`, `ok`
- Descriptive names for longer scopes: `formattedContent`, `usageCount`
- No camelCase for unexported vars with underscores (use `cmdFlag` not `cmd_flag`)

### Test Names
- Use table-driven tests with `name` field
- Test function names: `Test<FunctionName>` or `Test<Type><Method>`
- Subtest names: lowercase with underscores: `"simple_object_with_defaults"`

## Debugging Tips

### TUI Debugging
Bubble Tea takes over stdout, so normal `fmt.Println` won't work:

```go
// Log to a file
f, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
defer f.Close()
fmt.Fprintf(f, "Debug: %v\n", value)

// Or use bubbletea's debug mode
p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
```

### Test Debugging
```bash
# Verbose output
go test -v ./cmd/

# Run single test
go test -v -run TestJson2toonCmd/simple_object ./cmd/

# See what CLI command outputs
go test -v -run TestBase64CLI ./cmd/
```

### Linter Confusion
If linter complains but you don't see the issue:
```bash
# Run specific linter
golangci-lint run --enable-only=goconst

# See which linter triggered
golangci-lint run --print-issued-lines
```

## Resources

- **Bubble Tea Docs**: https://github.com/charmbracelet/bubbletea
- **Cobra Docs**: https://cobra.dev/
- **golangci-lint Linters**: https://golangci-lint.run/usage/linters/
- **Project Site**: https://devtui.com
- **Repository**: https://github.com/skatkov/devtui

## Questions or Improvements?

This guide should cover most scenarios. If something is unclear or you discover new patterns while working on the codebase, update this file to help future agents and developers.
