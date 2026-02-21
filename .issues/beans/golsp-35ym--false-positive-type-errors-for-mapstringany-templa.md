---
# golsp-35ym
title: False positive type errors for map[string]any template data
status: completed
type: bug
priority: normal
created_at: 2026-02-08T17:58:40Z
updated_at: 2026-02-08T21:35:08Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:25Z"
        task_id: 868hk17k6
---

## Problem

The \`check\` subcommand reports false positive type errors when template data is passed as \`map[string]any\`. Since the linter infers types purely from template usage (no Go source analysis or type annotations), it cannot verify fields accessed through map indirection. All reported errors are for fields that **do** exist and work correctly at runtime.

This is the most common pattern in real-world Go web apps — many projects (including ours) use \`map[string]any\` as template data rather than typed structs, because handlers build page data dynamically with varying fields per page.

## Reproduction

Given this Go handler:

\`\`\`go
data := map[string]any{
    "Title":       "My Page",
    "Page":        "login",
    "IsRunning":   true,
    "NumRuns":     42,
    "PortfolioID": 123,
    "View":        "table",
    "Units": unitsPageData{
        Sort:  "name",
        View:  "table",
        Query: "",
    },
}
tmpl.ExecuteTemplate(w, "my-page", data)
\`\`\`

And this template:

\`\`\`html
{{define "my-page"}}
  {{if eq .Page "login"}}Login{{end}}
  {{if .IsRunning}}Running ({{.NumRuns}} runs){{end}}
  {{template "sub-page" .}}
{{end}}

{{define "sub-page"}}
  {{- $u := .Units -}}
  {{if eq $u.View "table"}}Table View{{end}}
  Portfolio: {{.PortfolioID}}
{{end}}
\`\`\`

The linter reports errors like:

\`\`\`
my-page.gohtml:2:14: type mismatch, field not found: '$.Page' of type 'string'
my-page.gohtml:3:8: type mismatch, field not found: '$.IsRunning' of type 'any'
my-page.gohtml:3:30: type mismatch, field not found: '$.NumRuns' of type 'any'
sub-page.gohtml:3:14: type mismatch, expected 'string' but got 'invalid type'
sub-page.gohtml:4:16: type mismatch, field not found: '$.PortfolioID' of type 'any'
\`\`\`

## Real-world examples from Pacer

### ~~Category 1: "field not found" on map[string]any top-level fields~~ ✅ FIXED

Fixed in latest version. These 10 errors no longer appear:

| ~~\`auth.gohtml:15\`~~ | ~~\`$.Page\`~~ | ~~field not found on type 'string'~~ |
| ~~\`job/job.gohtml:7-9\`~~ | ~~\`$.IsRunning\`, \`$.AverageDuration\`, \`$.NumRuns\`~~ | ~~field not found~~ |
| ~~\`rm.gohtml:12\`~~ | ~~\`$.PortfolioID\`~~ | ~~field not found~~ |
| ~~\`res-search.gohtml:3-4\`~~ | ~~\`$.OTAOptions\`, \`$.Results\`~~ | ~~field not found~~ |
| ~~\`res-trends.gohtml:20\`~~ | ~~\`$.SortBy\`, \`$.Trends\`~~ | ~~field not found~~ |
| ~~\`portfolio.gohtml:60\`~~ | ~~\`$.Billing\`~~ | ~~field not found~~ |
| ~~\`calendar.gohtml:52\`~~ | ~~\`$.Grid.Dates\`~~ | ~~field not found~~ |

### ~~Category 1b: "field not found" — sub-template fields~~ ✅ FIXED

Fixed in v0.12.3. These 2 errors no longer appear:

| ~~\`auth-router.gohtml:2\`~~ | ~~\`$.Message\`~~ | ~~field not found on type 'any'~~ |
| ~~\`auth-router.gohtml:3\`~~ | ~~\`$.AllowedDomains\`~~ | ~~field not found on type 'map[any]any'~~ |

### ~~Category 2: "field or method not found" on nested struct fields~~ ✅ FIXED

Fixed in v0.12.3. These 8 errors no longer appear:

| ~~\`calendar-grid.gohtml:17\`~~ | ~~\`$unit.Name\`, \`$unit.Bedrooms\`, \`$unit.City\`~~ | ~~field or method not found~~ |
| ~~\`calendar-grid.gohtml:19-20\`~~ | ~~\`$unit.Name\`, \`$unit.Bedrooms\`, \`$unit.City\`~~ | ~~same~~ |
| ~~\`calendar-grid.gohtml:25\`~~ | ~~\`$unit.ID\`~~ | ~~field or method not found~~ |

### ~~Category 3: "expected 'string' but got 'invalid type'" for struct fields through map indirection~~ ✅ FIXED

Fixed in v0.12.2. These 7 errors no longer appear:

| ~~\`units.gohtml:29\`~~ | ~~\`eq $u.View "table"\`~~ | ~~expected 'string' but got 'invalid type'~~ |
| ~~\`units.gohtml:35\`~~ | ~~\`eq $u.View "map"\`~~ | ~~same~~ |
| ~~\`units.gohtml:71\`~~ | ~~\`eq (printf "%d" .) $u.FilterBedrooms\`~~ | ~~same~~ |
| ~~\`units.gohtml:85-87\`~~ | ~~\`eq $u.FilterConfidence "low/medium/high"\`~~ | ~~same~~ |
| ~~\`units.gohtml:102\`~~ | ~~\`eq $u.View "map"\`~~ | ~~same~~ |

## Possible solutions

- [x] Support \`map[string]any\` as template data — treat map keys as valid field names
- [ ] Add type annotation comments (\`{{/* @type MyStruct */}}\`) so templates can declare their expected data type
- [ ] Add a \`//gotmpls:ignore\` or \`{{/* gotmpls:ignore */}}\` directive to suppress specific errors
- [ ] Read Go source to discover the actual types passed to \`ExecuteTemplate\` calls (ambitious but ideal)

## Progress

- v0.12.1: Fixed Category 1 (map field access + cross-template partial struct compat). Error count: **27 → 17**.
- v0.12.2: Fixed Category 3 (struct fields through map indirection). Error count: **17 → 10**.
- v0.12.3: Fixed Categories 1b and 2 (sub-template fields + nested struct fields in range loops). Error count: **10 → 0**.

## Summary of Changes

All 27 false positive type errors from `map[string]any` template data have been resolved across three releases (v0.12.1–v0.12.3). The `check` subcommand now passes cleanly and is ready for CI adoption.
