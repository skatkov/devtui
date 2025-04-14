# **DevTUI** - A Swiss-army app for developers

**DevTUI** is an all-in-one terminal toolkit designed for developers. It consolidates your everyday utilities into a single, unified TUI (Text User Interface) and CLI (Command Line Interface). No more juggling multiple tools — everything you need is now just one command away.

🚀 Actively developed and continuously improving.

![devtui](devtui.gif)

## ✨ Features

DevTUI includes a growing collection of handy tools for developers:

- ✅ JSON Formatter
- ⏰ Cron Parser
- 🔑 UUID Generator
- 🔢 Base Number Converter
- 🎯 And more...

---

## 📦 Install
### 🧃 Homebrew
```
brew install skatkov/tap/devtui
```
### 📥 Download Executable

Or download executable from [GitHub Releases](https://github.com/skatkov/homebrew-tap/releases?q=devtui&expanded=true)

---

## 💡 Why DevTUI?

 - 🧰 Unified experience – Replace scattered tools with a single app
 - 🔒 Privacy-focused – Everything runs locally, no data ever leaves your computer
 - 🌐 Offline support – No internet? No problem
 - ⌨️ Built for the terminal – No need to reach for your mouse or browser
 - 🛠️ Actively maintained – Not just another abandoned open-source project

---

## 📋 Requirements

### macOS
- ✅ Works out of the box

### Linux

- 🖱 Wayland requires: `wl-clipboard`
- 🧮 X11 requires: `xclip` or `xsel`

To check your session type:
```bash
echo $XDG_SESSION_TYPE
# Output: wayland or x11
```
### Windows

- ⚠️ Should work, but not officially tested.
- Grab the executable from the GitHub Releases page.

---

## 🚀 Usage
DevTUI includes both a TUI and CLI interface.

### 🖥 TUI
Run `devtui` you’ll see a list of available tools — just pick one and go!

### Autocompletion
Run a one of these commands depending on shell

```
devtui completion bash > ~/.bashrc
devtui completion zsh  > ~/.zshrc
devtui completion fish > ~/.fishrc
```

### CLI (Experimental)
The CLI interface is still in development and may change in future versions.

#### 🎨 CSS Formatter
```bash
devtui cssfmt < testdata/bootstrap.min.css > output.css
```

Optional flags:
```
  -i, --indent int   spaces for indentation (default 2)
      --semicolon    always end rule with semicolon, even if not needed (default true)
  -t, --tab          use tabs for indentation
```
#### 🧼 CSS Minimizer
Strip unnecessary whitespace from CSS files:

```bash
devtui cssmin < input.css > output.min.css
```

#### 🗂 XML Formatter
```bash
devtui xmlfmt < testdata/input.xml > output.xml
```

Optional flags:

```
  -i, --indent string   Indent string for nested elements (default "  ")
  -n, --nested          Nested tags in comments
  -p, --prefix string   Each element begins on a new line and this prefix
```

#### 📝 GraphQL Query Formatter
Format GraphQL queries:

```bash
devtui gqlfmt < testdata/query.graphql

devtui gqlfmt < testdata/query.graphql > formatted.graphql
devtui gqlfmt --indent "    " --with-comments --with-descriptions < testdata/query.graphql
```

Optional flags:
```
  -i, --indent string       Indent string for nested elements (default is 2 spaces) (default "  ")
  -c, --with-comments       Include comments in the formatted output
  -d, --with-descriptions   Include descriptions in the formatted output (omitted by default)
```

#### 🗒️ TSV to Markdown Table
Convert TSV to Markdown Table:

```bash
devtui tsv2md -t < example.tsv          - convert tsv from stdin and view result in stdout

devtui tsv2md < example.tsv > output.md - convert tsv from stdin and write result in new file

cat example.tsv | devtui tsv2md         - convert tsv from stdin and view result in stdout
```

Optional flags:
```
  -a, --align           Align columns width
  -t --header string    Add main header (h1) to result
```

#### 🗒️ CSV to Markdown Table
Convert CSV to Markdown Table:

```bash
devtui csv2md -t < example.csv          - convert csv from stdin and view result in stdout

devtui csv2md < example.csv > output.md - convert csv from stdin and write result in new file

cat example.csv | devtui csv2md         - convert csv from stdin and view result in stdout
```

Optional flags:
```
  -a, --align           Align columns width
  -t --header string    Add main header (h1) to result
```

#### 🗒️ TSV to CSV
Convert TSV to CSV:

```bash
devtui tsv2csv < testdata/input.tsv > output.csv
```

Optional flags:
```
  -i, --indent string   Indent string for nested elements (default "  ")
  -n, --nested          Nested tags in comments
  -p, --prefix string   Each element begins on a new line and this prefix
```

---

## 🧑‍💻 Contact
I love when people reach out, so please don't hesitate to do that.

- contact@devtui.com
- [https://t.me/skatkov](https://t.me/skatkov)
- [https://bsky.app/profile/skatkov.com](https://bsky.app/profile/skatkov.com)
- [https://x.com/5katkov](https://x.com/5katkov)
