<!-- merged from seo-content/assets/link-inventory-template.md + seo-sprint/assets/link-inventory-template.md -->

# {{PRODUCT_NAME}} — Internal Link Inventory

> Every URL `/seo` can link to, and every URL it has shipped. Each new page picks ≥3 in-body links from here, takes ≥2 inbound links from existing pages (≥1 of them a frequently-crawled hub), and registers itself here on ship so the next page can link to it. Reuse and append; never overwrite.

## Existing pages (link targets)

### Homepage + core marketing

| Slug | URL | Title (anchor-text candidate) | Linked by |
|---|---|---|---|
| `/` | {{HOMEPAGE_URL}} | {{HOMEPAGE_TITLE}} | All |
| `/pricing` | {{PRICING_URL}} | Pricing | Comparisons, playbooks |
| `/about` | {{ABOUT_URL}} | About | (occasional) |

### Features

| Slug | URL | Title | Linked by |
|---|---|---|---|
{{FEATURES_INVENTORY}}

### Tools / free utilities

| Slug | URL | Title | Linked by |
|---|---|---|---|
{{TOOLS_INVENTORY}}

### Content (blog / guides / published pieces)

| Slug | URL | Title | Topic | Linked by |
|---|---|---|---|---|
{{CONTENT_INVENTORY}}

---

## Crawl hubs

> The pages Google crawls most often, from `batch_url_inspection` → `last_crawled` (`references/gsc.md` §7a), mirrored in `.seo/config.json` under `crawl_hubs`. **At least one inbound link on every new page must come from this list** — two links from pages Google visits quarterly is close to no discovery signal at all. Re-derive every few months.

| URL | Last crawled | Read on |
|---|---|---|
{{CRAWL_HUBS_TABLE}}

---

## Programmatic pages

> Filled by the programmatic lane as pattern batches ship. Each row records where the page links out to and what links in, so the link spine is auditable without re-crawling.

### `/alternatives/[slug]`

| Slug | Ships in phase | URL | Inbound links from | Outbound links to |
|---|---|---|---|---|
{{ALTERNATIVES_TABLE}}

### `/for/[slug]`

| Slug | Ships in phase | URL | Inbound links from | Outbound links to |
|---|---|---|---|---|
{{FOR_TABLE}}

### `/compare/[slug]`

| Slug | Ships in phase | URL | Inbound links from | Outbound links to |
|---|---|---|---|---|
{{COMPARE_TABLE}}

### `/playbooks/[slug]`

| Slug | Ships in phase | URL | Inbound links from | Outbound links to |
|---|---|---|---|---|
{{PLAYBOOKS_TABLE}}

---

## Editorial pieces shipped

> Appended on every run that ships a piece. Each one becomes a future link target.

| Slug | URL | Title | Type | Inbound links from | Anchor-text variations |
|---|---|---|---|---|---|
{{SHIPPED_PIECES_TABLE}}

---

## Anchor-text variations (avoid repetition)

When linking to the same destination from multiple pages, vary the anchor text. Repeating "click here" or the destination's exact title weakens diversity signals.

Example variations for a hypothetical `/features/twitter-monitoring`:

- "monitor Twitter mentions"
- "real-time Twitter monitoring"
- "Twitter mention tracking"
- "our Twitter monitoring tool"
- "{{PRODUCT_NAME}}'s Twitter watch feature"

When you ship a page, append 4-5 anchor variations for it here so future runs can rotate them.
