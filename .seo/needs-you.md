# Needs you

Decisions, logins, reviews, and keys only a human can provide. The skill writes the `question` column. You write the `answer` column. The next run acts on any row with an answer and no `closed` date, then stamps it.

Keep answers to one line where you can. If the answer is "never" or "not now," say that; the row closes and the blocked candidate is dropped or deferred.

| id | opened | blocks | question | answer | closed |
|---|---|---|---|---|---|
| NY-1 | 2026-11-10 | sitemap repair (discovery) | After the site deploys, submit https://devtui.com/sitemap.xml in Search Console > Sitemaps for sc-domain:devtui.com. Done? (yes/no) | yes, submitted 2026-11-10 | 2026-11-10 |
| NY-2 | 2026-11-10 | Bing panel | Connect Bing Webmaster Tools (`seo providers bing connect`) and submit the sitemap there too? (a) yes (b) not now (c) never | a (2026-11-10). Site added to Bing Webmaster and sitemap submitted 2026-11-10. API connected via `seo providers bing connect` (keychain) 2026-11-10; no Bing data yet (new site) | 2026-11-10 |
| NY-3 | 2026-11-10 | brand voice | Confirm the inferred voice and persona in .seo/brand.md: "plain, short, developer-to-developer; persona = individual devs on macOS/Linux who live in the terminal". (a) correct (b) change: ... | a, correct | 2026-11-10 |
| NY-4 | 2026-11-12 | production verification of the 2026-11-10 generator fixes and today's homepage title | Local main is 5 commits ahead of origin (plus today's uncommitted homepage change); pushing main deploys devtui.com. Push it? (yes/no) | yes, pushed as d442789; deploy run 43311761439 succeeded, production verified 2026-11-12 | 2026-11-12 |
| NY-5 | 2026-11-12 | AEO value estimates | AI referrals are not measured (no GA4 or log access configured; the site uses Umami at umami.skatkov.com). Expose Umami referrer data to the skill? (a) yes, paste a 28-day referrer export each month (b) not now (c) never | a (2026-11-12): owner pastes a 28-day Umami referrer export monthly | 2026-11-12 |
