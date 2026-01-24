---
title: serve
parent: CLI
---

## devtui serve

Start DevTUI as an SSH server

### Synopsis

Start DevTUI as an SSH server that users can connect to via SSH.

This also starts a simple HTTP server to show a landing page for web browsers.

Examples:
  # Start with default settings (SSH on :2222, HTTP on :8080)
  devtui serve

  # Custom ports
  devtui serve --ssh-port 22 --http-port 80

  # Specify host key path
  devtui serve --host-key /path/to/host_key

```bash
devtui serve [flags]
```

### Options

```
  -h, --help               help for serve
      --host-key string    Path to SSH host key (default ".ssh/devtui_host_key")
      --http-host string   HTTP server host (default "0.0.0.0")
      --http-port string   HTTP server port (default "8080")
      --ssh-host string    SSH server host (default "0.0.0.0")
      --ssh-port string    SSH server port (default "2222")
```
