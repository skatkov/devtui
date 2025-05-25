# TUI Demo Tapes

This directory contains VHS tape files for generating animated GIF demonstrations of all TUI modules in the devtui project.

## Prerequisites

Make sure you have VHS installed:

```bash
go install github.com/charmbracelet/vhs@latest
```

## Usage

### Generate All Screenshots

To generate screenshots for all TUI modules at once:

```bash
cd docs
./run_demos.sh
```

This will:
- Generate a main menu demo as `site/assets/img/devtui-main.png`
- Generate individual screenshots for each TUI module in `site/assets/img/tui/`

### Generate Individual Screenshots

To generate a screenshot for a specific module:

```bash
cd docs
vhs tapes/demo-<module-name>.tape
```

For example:
```bash
vhs tapes/demo-json.tape
vhs tapes/demo-yaml.tape
vhs tapes/demo-base64-encoder.tape
```

## File Structure

- `demo-main.tape` - Main menu navigation demo
- `demo-<module>.tape` - Individual module demonstrations
- `run_demos.sh` - Script to generate all screenshots

## Output Locations

Generated files are placed in:
- **PNG screenshots**: `site/assets/img/tui/` (for documentation)
- **Main demo**: `site/assets/img/devtui-main.png`

## How It Works

Each tape file:
1. Starts the devtui application with `go run .`
2. Presses `/` to start search mode
3. Types the module name to search for it
4. Presses `Enter` to select the top result
5. Presses `Enter` again to enter the module
6. Takes a screenshot immediately showing the tool's interface
7. Exits back to the main menu

The search approach ensures reliable navigation regardless of menu item order or usage statistics.

## Customization

To modify the demonstrations:
1. Edit the tape files directly, or
2. Modify the generation logic in `../tape_docs.go` and regenerate with `go run .`