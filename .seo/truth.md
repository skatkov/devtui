# Truth ledger: DevTUI

Derived from code on 2026-11-10. Re-derive whenever `.goreleaser.yaml`, `cmd/`, `tui/` or the install docs change.

| Claim | Value | Source | Read |
|---|---|---|---|
| Price | Free, MIT licence | LICENSE | 2026-11-10 |
| Platforms with release binaries | macOS (darwin) and Linux | .goreleaser.yaml `goos` | 2026-11-10 |
| Windows binaries | **no** | .goreleaser.yaml `goos` (no windows) | 2026-11-10 |
| Install: Homebrew | `brew install skatkov/tap/devtui` | README.md, site/start.md | 2026-11-10 |
| Install: FreeBSD ports | `pkg install devtui` (community-maintained) | README.md (missing from site/start.md) | 2026-11-10 |
| Install: GitHub Releases | yes | README.md, site/start.md | 2026-11-10 |
| CLI subcommands | 31 (base64 … yamlstruct, plus version and completion) | cmd/*.go | 2026-11-10 |
| TUI modules | 26 (tui/* except root) | tui/*/ | 2026-11-10 |
| Runs offline / locally | yes, with no network calls in the tools | homepage claim; tools are local transforms | 2026-11-10 |
| Linux clipboard requirement | wl-clipboard (Wayland), xclip or xsel (X11) | site/start.md | 2026-11-10 |

## Code vs site disagreements

- `site/start.md` doesn't list the FreeBSD install that the README lists. Low risk; it's an omission, not a contradiction.
- `site/tui/` has docs for 3 TUI modules that no longer exist under those names (csvjson, jsontoml, tomljson). They duplicate csv2json, json2toml and toml2json.
- The `tui/password` module has no docs page. `site/tui/url-extractor.md` has an empty title.

## HIGH-RISK

None. This is a free dev tool with no compliance, medical, legal or financial claims. The riskiest claim is platform support (Windows).
