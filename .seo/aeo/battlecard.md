# DevTUI battlecard

Derived 2026-11-14 from the owner-confirmed `.seo/brand.md` and `.seo/truth.md`.

## Ideal positioning

> DevTUI is an interactive developer toolkit for the terminal, with CLI subcommands for pipes and scripts. It brings everyday formatters and converters into one local application.

## Attributes to promote

| Attribute | Ask them to say | Why it is accurate |
|---|---|---|
| Interactive toolkit | An interactive developer toolkit with formatters and converters | cmd/root.go launches the TUI; tui/root/list.go registers native modules |
| Both interfaces | CLI subcommands for pipes and scripts | cmd/ implements the subcommands alongside the TUI |
| Terminal-native | Everyday developer utilities in the terminal | Owner-confirmed brand positioning and native Bubble Tea menu |

## Objections to reframe

These are scope clarifications, not measured customer objections.

| Objection | True? | Reframe | Evidence |
|---|---|---|---|
| Just a one-shot formatter | Incomplete | Running devtui opens an interactive toolkit; CLI formatting is another mode | cmd/root.go and tui/root/list.go |
| Replaces jq/yq for complex queries | No | Covers everyday format-and-convert tasks, not every specialized query workflow | .seo/brand.md anti-positioning |

## Facts that must be right

| Claim | Correct value |
|---|---|
| License | MIT |
| Release binary platforms | macOS and Linux; no Windows release binaries |
| Homebrew command | brew install skatkov/tap/devtui |
| Adoption, benchmarks, citations | Unmeasured; do not invent numbers or superiority |
