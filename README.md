# **DevTUI** - A Swiss-army app for developers

All-in-one terminal toolkit that consolidates everyday developer utilities into a unified TUI and CLI interfaces.

![devtui](/site/assets/img/devtui.png)

## 💡 Why DevTUI?

 - 🧰 Unified experience – Replace scattered tools with a single app
 - 🔒 Privacy-focused – Everything runs locally, no data ever leaves your computer
 - 🌐 Offline support – No internet? No problem
 - ⌨️ Built for the terminal – No need to reach for your mouse or browser

---

## 📦 Install
### 🧃 Homebrew
```
brew install skatkov/tap/devtui
```
### 📥 Download Executable

Or download executable from [GitHub Releases](https://github.com/skatkov/devtui/releases)

## 🚀 Docs
-> [devtui.com/start](https://devtui.com/start)

---

## 📚 Documentation Generator

DevTUI includes automated documentation generators for both CLI and TUI interfaces (not complete yet, though).

### Generate All Documentation
To regenerate both CLI and TUI documentation:

```bash
cd docs && go run *.go
```

This will:
- Generate CLI documentation in `site/cli/` with proper Jekyll front matter
- Generate TUI documentation in `site/tui/` with key bindings and usage instructions
- Clean up auto-generated content (remove footers, SEE ALSO sections, etc.)
- Apply proper formatting and language hints for code examples

### Individual Generators
You can also run generators separately:

```bash
# CLI documentation only
cd docs && go run cli-docs.go docs.go

# TUI documentation only
cd docs && go run tui-docs.go docs.go
```

---

## Logo
Logo was done by [Andrei Kedrin](https://linktr.ee/andreikedrin).

Figma original:
https://www.figma.com/design/JTS0mzphMDiRuuC3xNprLu/Untitled?node-id=0-1&p=f&t=0LeB0uhXSUmZpE3Q-0
