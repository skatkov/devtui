# DevTUI: brand and voice contract

Written 2026-11-10 by `/seo` foundation. Positioning confirmed by the owner on 2026-11-10. Voice and persona confirmed by the owner on 2026-11-10.

## Positioning (confirmed 2026-11-10)

> We are a **single-binary terminal toolkit** for **developers who live in the terminal** who need **everyday formatters and converters (JSON, YAML, TOML, CSV, base64, UUID and more) without pasting data into websites**.
> We win because **it is one local, offline app with both an interactive TUI and a pipeable CLI**.
> We are **not** for **people who want a GUI or web tool, or an IDE plugin**.

## Product one-liner

All-in-one terminal toolkit that brings everyday developer utilities together in one TUI and CLI. (From the README and the homepage.)

## Persona

- Primary: individual developers on macOS and Linux who work in a terminal all day (Go, Ruby, JS and backend devs, SREs). (confirmed 2026-11-10)
- Not B2B. Free, MIT-licensed open source with GitHub Sponsors.

## Competitors (confirmed 2026-11-10)

DevToys (desktop GUI), DevUtils (macOS app), IT-Tools (web), CyberChef (web), single-purpose CLIs (jq, yq).

## Differentiators

1. Both modes in one binary: an interactive TUI with a menu, and CLI subcommands for pipes and scripts.
2. Runs locally and offline, so no data leaves the machine (homepage claim).
3. Terminal-native: keyboard-driven, works over SSH, clipboard and `$EDITOR` integration.

## Anti-positioning

Not a GUI. Not a web app. Not an IDE extension. Does not try to replace jq or yq for complex queries; it covers the everyday format-and-convert jobs.

## Voice (confirmed 2026-11-10)

Plain, short, developer-to-developer. Bullet points over paragraphs. Concrete commands over adjectives. Light humour is fine ("No internet? No problem").
Forbidden: hype words ("revolutionary", "blazing", "game-changing"), invented benchmarks, claims of platform support that `.goreleaser.yaml` does not build.

## Forbidden claims

- Windows support (goreleaser builds darwin and linux only).
- Any number of users, stars or downloads that isn't read from a live source the same day.

## Author

Stanislav (Stas) Katkov, https://skatkov.com, the maintainer.

## Proprietary data

None reachable. The local usage stats in `~/.config/devtui/usage_stats.json` never leave the user's machine and must never be used.
