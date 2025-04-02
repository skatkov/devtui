**DevTUI** - A Swiss-army app for developers

It comes with a number of small utility apps that developers use in their day-to-day work. Such as:
- JSON Formatter
- Cron parser
- UUID generator
- Base Number Converter

It is still being actively developed and improvements and more tools are being added often.

![devtui](devtui.gif)

## Install
Through homebrew:
```
brew install skatkov/tap/devtui
```

Or download executable from [GitHub Releases](https://github.com/skatkov/homebrew-tap/releases?q=devtui&expanded=true)

## Requirements
With **OSX** everything should work out of the box.

On Linux, **Wayland** requires `wl-clipboard` and **X11** requires `xclip` or `xsel` to be installed.

Check your session type to correctly identify clipboard manager to use.

```
echo $XDG_SESSION_TYPE
# wayland or X11
```

Windows should work, but currently is not tested. Please see [Github Releases] (https://github.com/skatkov/homebrew-tap/releases) for a windows binary.
## Usage
Application comes with TUI and CLI interfaces.

### TUI
TUI could be accessed by running a `devtui` in your terminal and you can see the list of available apps. Pick one and go.

### CLI
CLI interface is still experimental and could be a subject to change.

#### Autocomplete
 //TODO

#### CSS Formatter
```bash
devtui cssfmt < testdata/bootstrap.min.css > output.css
```

There are also additional flags that you can pass on
```
  -i, --indent int   spaces for indentation (default 2)
      --semicolon    always end rule with semicolon, even if not needed (default true)
  -t, --tab          use tabs for indentation
```
#### CSS Minimizer
This tools just removes whitespace from css files. Basically similar to CSS Formatter, but with preconfigured options to remove whitespaces.

```bash
devtui cssmin < input.css > output.min.css
```

#### XML Formatter
```bash
devtui xmlfmt < testdata/input.xml > output.xml
```

There are also additional flags that you can pass on

```
  -i, --indent string   Indent string for nested elements (default "  ")
  -n, --nested          Nested tags in comments
  -p, --prefix string   Each element begins on a new line and this prefix
```  

## Contact
I love when people reach out, so please don't hesitate to do that.

- [https://t.me/skatkov](https://t.me/skatkov)
- [https://bsky.app/profile/skatkov.com](https://bsky.app/profile/skatkov.com)
- [https://x.com/5katkov](https://x.com/5katkov)
