Setup# DevTUI Documentation

Welcome to the **DevTUI** project! This is a comprehensive developer tool built with Go and the Bubble Tea framework.

## Quick Links

- [GitHub Repository](https://github.com/skatkov/devtui)
- [Official Documentation](https://devtui.dev/docs)
- [API Reference](https://api.devtui.dev/v1)

## Installation

You can install DevTUI from multiple sources:

1. **Direct download**: https://releases.devtui.dev/latest
2. **Homebrew**: Visit brew.sh for installation instructions
3. **Docker**: Pull from docker.io/devtui/devtui:latest
4. **Source**: Clone from git.sr.ht/~skatkov/devtui

## Configuration

The configuration file should be placed at:
- Linux: `~/.config/devtui/config.yaml`
- macOS: `~/Library/Application Support/devtui/config.yaml`
- Windows: `%APPDATA%\devtui\config.yaml`

Example configuration available at: https://raw.githubusercontent.com/skatkov/devtui/main/config.example.yaml

## Features

### Base64 Tools
- Encode and decode base64 strings
- Support for binary files
- More info: https://en.wikipedia.org/wiki/Base64

### JSON Tools
- Format and validate JSON
- Convert between JSON and other formats
- JSON Schema validation: jsonschema.org

### YAML Processing
- Format YAML files
- Convert YAML to JSON and vice versa
- Specification: yaml.org/spec/1.2/spec.html

## Community

Join our community channels:

- Discord: discord.gg/devtui-community
- Reddit: reddit.com/r/devtui
- Stack Overflow: stackoverflow.com/questions/tagged/devtui
- Discussions: github.com/skatkov/devtui/discussions

## API Endpoints

Our REST API is available at multiple endpoints:

| Environment | Base URL | Documentation |
|------------|----------|---------------|
| Production | https://api.devtui.dev | https://docs.devtui.dev/api |
| Staging | https://staging-api.devtui.dev | https://staging-docs.devtui.dev/api |
| Development | http://localhost:8080 | http://localhost:8080/docs |

## External Resources

- Go documentation: go.dev/doc
- Bubble Tea framework: github.com/charmbracelet/bubbletea
- Cobra CLI library: github.com/spf13/cobra
- Lipgloss styling: github.com/charmbracelet/lipgloss

## Support

Need help? Check these resources:

1. [FAQ](https://devtui.dev/faq)
2. [Troubleshooting Guide](https://devtui.dev/troubleshooting)
3. [Email Support](mailto:support@devtui.dev)
4. [Bug Reports](https://github.com/skatkov/devtui/issues)

## Contributing

We welcome contributions! Please read our [contributing guide](https://github.com/skatkov/devtui/blob/main/CONTRIBUTING.md) and check the [good first issues](https://github.com/skatkov/devtui/labels/good%20first%20issue).

### Development Setup

```bash
git clone https://github.com/skatkov/devtui.git
cd devtui
go mod download
go run main.go
```

For more detailed setup instructions, see: https://devtui.dev/development-setup

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/skatkov/devtui/blob/main/LICENSE) file for details.

## Acknowledgments

Special thanks to:
- The Charm team for Bubble Tea: charm.sh
- The Go community: golang.org/community
- All our contributors: github.com/skatkov/devtui/graphs/contributors

---

For more information, visit our website at https://devtui.dev or check out the source code at https://github.com/skatkov/devtui.
