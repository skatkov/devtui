---
title: Getting started
nav_order: 2
---



## Install
{: #install}
### Homebrew
```
brew install skatkov/tap/devtui
```
### Download Executable

Or download executable from [GitHub Releases](https://github.com/skatkov/devtui/releases)

---

## Requirements

### macOS
- Works out of the box

### Linux

- Wayland requires: `wl-clipboard`
- X11 requires: `xclip` or `xsel`

To check your session type:
```bash
echo $XDG_SESSION_TYPE # Output: wayland or x11
```
