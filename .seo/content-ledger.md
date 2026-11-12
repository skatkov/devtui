# DevTUI — Content Ledger

> The memory of the engine. Read first on every `/seo` run: the **Shipped** table is the dedup record (never re-write a covered topic); the **Performance** table is the scoreboard (did the last pieces actually work?); the **Candidate backlog** is the scored shortlist so each run starts warm. Updated in the same edit batch as every piece shipped.

---

## Shipped

| Date | Title | Type | Slug / URL | Target keyword | Vol | Bucket | Original data (source · n · as-of) | Refresh due | Primary internal links | Commit / PR |
|---|---|---|---|---|---|---|---|---|---|---|


<!-- Append one row per piece at Step 5. Type ∈ guide | how-to | listicle | definition | comparison | data-study | resource | opinion | case-study | tool -->
<!-- Tools: add a maintenance note to the row for anything with decaying data (fee schedules, rates, limits) — what expires and roughly when. -->
<!-- Original data: name the first-party element the piece carries (or "none — <ladder rung used>"). Refresh due: quarterly for volatile metrics, annual for structural ones, n/a if the piece carries no original data. The saved query lives in .seo/briefs/<slug>.md — re-running it is the cheapest high-value content run this engine has. -->

---

## Performance

> The closed loop. Filled by **Step 0.5** of a later run from GSC (`references/gsc.md` §2), seeded as `unmeasured` when a piece ships. Runs are stateless, so an unwritten measurement never happened. State ∈ winning | close | wrong-query | invisible | unmeasured | too-new | no-data (pre-window).
>
> **Backfilled:** 2026-11-10 (no editorial pieces exist; 61 generated docs pages, all `watch` in census) — set on the first GSC-aware run over an existing library (`gsc.md` §2b). Until it's set, assume no piece has ever been measured and backfill all of them in one pass before selecting anything.
>
> Adding this section to an older ledger? **Append it and match that file's own column style.** Don't renumber, reorder, or normalize the sections already there.

| Slug / URL | Published | Indexed? (state · checked) | Read @28d (clicks · impr · pos) | Read @56d (clicks · impr · pos) | Site-wide same window (clicks · impr) | Best lever (recover/CTR/rank · est. clicks) | State | Note / next action |
|---|---|---|---|---|---|---|---|---|
| https://devtui.com/sitemap.xml (site-wide repair) | 2026-11-10 | deployed; sitemap 59 URLs (2026-11-12) | | | 5 · 133 (2026-10-10..11-06) | discovery · 41 URLs unknown to Google | unmeasured | Read GSC known-URL count at 21d: was 20 of 61 |
| https://devtui.com/ (title + description boost) | 2026-11-12 | indexed (homepage) | | | 7 · 217 (2026-10-13..11-09) | CTR · pre 5.7% at pos 5.8 | unmeasured | Live 2026-11-12. Read CTR at 28d (~2026-12-10) |

<!-- Two clocks. Index status: check ~3 days after deploy, on anything published in the last 30 days not yet confirmed indexed. Performance: measure any piece live 21+ days and not measured in the last 30. -->
<!-- Indexed? holds the coverage state verbatim (e.g. "Submitted and indexed", "Crawled - currently not indexed", "URL is unknown to Google", "noindex") plus the date checked. Each state has a different fix — see references/gsc.md §3a. -->
<!-- Read at ~28 and ~56 days after publish (or after a boost lands), not "whenever the skill next runs" — fixed checkpoints make rows comparable across pieces. -->
<!-- Site-wide same window is the CONTROL. Without it a Google-wide lift reads as your edit working. One get_performance_overview call. -->
<!-- Best lever: the largest of recover / CTR / rank in estimated clicks (references/gsc.md §2c). It decides both whether to boost and what kind of fix to write. -->
<!-- LOOP GATE: three published pieces in a row not indexed = the next run fixes that, not writes a fourth. -->
<!-- close (pos 5-15) → offer the boost at the next checkpoint, it often beats a new piece. -->
<!-- wrong-query → the queries it DID attract are a free candidate list for selection. -->
<!-- invisible after 21+ days → run inspect_url_enhanced BEFORE writing anything adjacent; not-indexed, noindex, wrong canonical, and thin each need a different fix. -->
<!-- Three invisible pieces in a row on one vein = the vein is dead. That verdict outranks any keyword-tool score. -->

---

## Candidate backlog

> Scored shortlist from the last selection run. Re-score when this is >30 days old or the user says "re-research." The next run reads this before regenerating the pool — it starts from here, validates the top pick is still open and winnable, and only does fresh research if needed.

| Rank | Candidate | Proposed type | Target keyword | Vol | Bucket (E/M/H · src) | Intent | Data angle | Score | Notes / angle |
|---|---|---|---|---|---|---|---|---|---|
| 1 | /tui/markdown.html close on "tui markdown renderer" | refresh | tui markdown renderer | | | info | none | 19 | 36 impr pos 9.6 (2026-11-12); query pos 5.4, 0 clicks |
| 2 | Offpage brief: awesome-tuis / Terminal Trove listings | offpage-brief | — | | | nav | none | 17 | backlinks never measured; link-health vital has no authority source |
| 3 | Docs page titles are short ("base64 \| DevTUI", ~15 chars) and say nothing about the job | refresh | e.g. "base64 decode cli" | | | info | none | | title is also the just-the-docs nav label; needs a separate SEO title field or a layout override |

<!-- Score = sum of winnability + traffic-potential + conversion-intent + strategic-value + data-angle + (6 - effort), each 1-5. See references/opportunity-research.md Step C. -->
<!-- Data angle: the specific first-party number or public-data combine this piece would carry, per references/proprietary-data.md. "None found" is a legitimate entry — an angle nobody can name isn't one. -->

---

## Coverage map (optional)

> A running view of which clusters/themes have content and which are thin or empty. Helps spot topical-depth gaps the keyword tools miss. Fill in as the library grows.

| Cluster / theme | Pieces shipped | Gaps still open |
|---|---|---|
| CLI command reference | 31 generated pages | thin (median 94 words), no examples narrative |
| TUI module reference | 29 generated pages (3 stale dupes) | thin |

---

## Notes

- **Difficulty buckets, not a KD cap:** Easy / Medium / Hard (`references/research-recipes.md`). Your playable bucket is the hardest one where GSC shows 2+ page-1 positions. Target it; allow one stretch pick in four one bucket up. Record the bucket AND which vendor's KD produced it — KD scales are vendor-specific and must never be compared across vendors.
- **No duplication:** before adding a candidate, check it isn't already a programmatic page in `.seo/roadmap.md` or a shipped row above.
- **One piece per run.** This ledger grows by one `Shipped` row per invocation.
- **Measure before you write.** Step 0.5 reads GSC for every piece live 21+ days and updates `Performance`. A `close` piece (position 5-15) is a boost candidate that regularly beats writing anything new, and an `invisible` piece is a question to answer before shipping its siblings.
- **Refresh beats rewrite.** Scan the `Refresh due` column at the start of every run. A piece whose original data is stale re-cuts in a fraction of a new run's cost — same URL, new numbers, renewed freshness, no new research — and answer engines reward the update. Do the refresh before starting a new piece when one is due.
- **Tool cadence:** if no `tool` row appears in the last 6 shipped pieces and a viable tool candidate exists, it takes the explore slot at the next checkpoint; none in the last 10 and it becomes the recommended pick. A library with no tool in it has no link magnet.

---

## Actions log

| date | action | target | type | lane | what changed | commit |
|---|---|---|---|---|---|---|
| 2026-11-10 | repair | https://devtui.com/sitemap.xml, https://devtui.com/robots.txt | technical | technical | Added jekyll-sitemap (62 URLs, no lastmod) and robots.txt with Sitemap line; both were 404. Pre-state: GSC knew 20 of 61 URLs. | uncommitted |
| 2026-11-10 | repair+consolidate | https://devtui.com/cli/*, https://devtui.com/tui/*, /start.html | technical | technical | Generators (docs/cli-docs.go, docs/tui-docs.go) now emit a unique description per page and an H1 on CLI pages; TUI title regex handles const blocks (url-extractor titled, de-orphaned); stale TUI pages auto-removed (deleted /tui/csvjson.html, /tui/jsontoml.html, /tui/tomljson.html, no 301: GitHub Pages; 0 GSC impressions); H1 + description on /cli/ and /start.html, description on /tui/. Pre-state: 38 violations live (62 shared descriptions, 33 h1-count, 3 dup titles, 1 orphan). | 21cae3f (deployed 2026-11-12) |
| 2026-11-12 | refresh | https://devtui.com/ | boost (CTR) | editorial | Title "Home \| DevTUI" → "DevTUI \| Developer Tools for the Terminal (TUI + CLI)" (53 ch) via page title = site title + new `site.tagline`; own meta description (153 ch, was site.description at 259 ch, cut in SERP); last_modified_at 2026-11-12 (sitemap lastmod + JSON-LD dateModified). Pre-state: 6 clicks / 106 impr / CTR 5.7% / pos 5.8 (2026-10-13..11-09); "tui dev" 2/13 pos 1.5. | d442789 (deployed 2026-11-12) |
