# awesome-tuis listing brief

Status: drafted, not submitted · Created / last verified: 2026-11-14 · Finding: OPP-01

### Brief: Awesome TUIs — https://github.com/rothgar/awesome-tuis#development
**Relationship:** none established; independent community directory · **Contact:** repository maintainers through their documented pull-request process; no personal contact known.
**Closes:** OPP-01 (only when inclusion is confirmed, not when this brief is drafted).

**Currently says:**
> "This is a community maintained list of TUI applications. A TUI application runs in your terminal and has some level of interactivity."
>
> "If you have a cool tool you'd like to share please [open a PR](https://github.com/rothgar/awesome-tuis/pulls)."

There is no DevTUI entry in the current README. GitHub issue/PR title-and-body search for `devtui` returned zero results; this is a limited search, not proof that no differently named submission ever existed.

**Problem:** DevTUI is absent from a relevant terminal-application directory. This is a verified listing gap, not a measured backlink gap or answer-engine citation gap. Citation frequency, authority, traffic and expected clicks are unmeasured; no uplift is promised.

**Requested:**
> - [DevTUI](https://github.com/skatkov/devtui) An interactive developer toolkit with formatters, converters, and CLI subcommands for pipes and scripts.

**Why it is accurate:** this describes the owner-confirmed positioning in `.seo/brand.md` and `.seo/aeo/battlecard.md`. `cmd/root.go:22-28` launches the Bubble Tea app without a subcommand; `tui/root/list.go` registers the native interactive modules, including JSON, TOML, YAML, Markdown and base64 tools. The CLI is implemented in `cmd/`. The entry makes no platform, benchmark, adoption or superiority claims.

**Ask:** add this one entry to the Development section of README.md, between DevProjex and differ, if the maintainers consider DevTUI a useful distinct addition.
**Effort for them:** approximately 2–5 minutes to review one line; testing the app may take longer. Acceptance and timing are their decision.

## Human submission instructions

1. Recheck the README and open PRs immediately before submitting, to avoid a duplicate.
2. Read the [pull-request template](https://github.com/rothgar/awesome-tuis/blob/main/.github/pull_request_template.md).
3. In your own fork of `rothgar/awesome-tuis`, edit only README.md and insert the requested line after DevProjex and before differ in Development.
4. Open a pull request against `rothgar/awesome-tuis:main`, titled **Add DevTUI to Development**. Keep the venue's template and add the short explanation below.
5. Report the PR URL and submission date here in chat. Recheck for inclusion 14 days after submission; do not send unsolicited follow-ups on a schedule.

Suggested PR explanation:

> I maintain DevTUI, an interactive terminal toolkit for everyday developer tasks, with a CLI mode for pipes and scripts.
>
> The repository was created on February 8, 2025, so it meets the six-month age requirement. Running `devtui` without a subcommand launches its own interactive Bubble Tea interface; it is not a wrapper around another interactive command or just a one-shot output formatter.
>
> Its proposed distinction is a single interactive toolkit covering multiple everyday developer utilities, alongside CLI subcommands. Please judge whether that is sufficiently distinct from the tools already listed.
>
> Documentation: https://devtui.com/
> Screenshot: https://github.com/skatkov/devtui/blob/main/site/assets/img/devtui.png

## Eligibility and fact checks

| Check | Result | Source / last verified |
|---|---|---|
| Repository at least six months old | PASS: created 2025-02-08T12:35:53Z | GitHub repository API, 2026-11-14 |
| Repository not archived or disabled | PASS; latest push reported 2026-11-12T21:51:47Z | GitHub repository API, 2026-11-14; recency is not a promise of future maintenance |
| Interactive native interface | PASS by source inspection: root launches Bubble Tea, menu registers modules | cmd/root.go and tui/root/list.go, 2026-11-14 |
| Not merely a formatter / interactive-command wrapper | PASS by source inspection of the root and module menu | Same source paths, 2026-11-14; no fresh runtime test |
| Unique application/interface | EDITORIAL REVIEW REQUIRED | Toolkit distinction stated, not certified; upstream maintainers decide |
| Existing entry / prior submission | No entry; zero title/body issue-search results for devtui | Upstream main README and GitHub search API, 2026-11-14 |
| Requested wording / brand | PASS: matches interactive toolkit + pipeable CLI positioning | .seo/brand.md and .seo/aeo/battlecard.md |
| Citation or traffic impact | UNCHECKED: no answer-engine sample or directory analytics | Do not imply this is a proven citation source |
| Submission, contact or publication | NOT DONE, intentionally | Human submits; approval of this brief is not authorization for third-party edits |

Evidence: `.seo/evidence/2026-11-14-awesome-tuis.json`.

Off-page inclusion is slower and less certain than an on-page correction. Do not hold future on-page work for this placement, pay for inclusion, solicit reviews, or claim an outcome before it happens.
