---
# golsp-35ym
title: False positive type errors for map[string]any template data
status: in-progress
type: bug
priority: normal
created_at: 2026-02-08T17:58:40Z
updated_at: 2026-02-08T18:04:04Z
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

## Real-world examples from Pacer (27 errors, all false positives)

### Category 1: "field not found" on map[string]any top-level fields

These fields are all set via \`map[string]any\` data and work at runtime:

| File | Field | Linter says | Reality |
|------|-------|-------------|---------|
| \`auth.gohtml:15\` | \`$.Page\` | field not found on type 'string' | Set as \`"Page": "login"\` |
| \`auth-router.gohtml:2\` | \`$.Message\` | field not found on type 'any' | Accessed conditionally, nil-safe |
| \`auth-router.gohtml:3\` | \`$.AllowedDomains\` | field not found on type 'map[any]any' | Set as \`"AllowedDomains": slice\` |
| \`job/job.gohtml:7\` | \`$.IsRunning\` | field not found on type 'any' | Set as \`"IsRunning": bool\` |
| \`job/job.gohtml:8\` | \`$.AverageDuration\` | field not found on type 'struct{...}' | Set as \`"AverageDuration": *time.Duration\` |
| \`job/job.gohtml:9\` | \`$.NumRuns\` | field not found on type 'any' | Set as \`"NumRuns": int\` |
| \`rm.gohtml:12\` | \`$.PortfolioID\` | field not found on type 'any' | Set by some handlers |
| \`res-search.gohtml:3\` | \`$.OTAOptions\` | field not found on type 'any' | Set as \`"OTAOptions": []MultiSelectOption\` |
| \`res-search.gohtml:4\` | \`$.Results\` | field not found on type 'struct{...}' | Set as \`"Results": *SearchResults\` |
| \`res-trends.gohtml:20\` | \`$.SortBy\` | field not found on type 'string' | Set as \`"SortBy": string\` |
| \`res-trends.gohtml:20\` | \`$.Trends\` | field not found on type 'struct{...}' | Set as \`"Trends": *Trends\` |
| \`portfolio.gohtml:60\` | \`$.Billing\` | field not found on type 'struct{...}' | Set as \`"Billing": portfolio.Billing\` |

### Category 2: "field or method not found" on nested struct fields

The linter partially infers struct types from template usage but misses fields:

| File | Field | Reality |
|------|-------|---------|
| \`calendar-grid.gohtml:17\` | \`$unit.Name\` | \`GridUnit.Name\` exists (string) |
| \`calendar-grid.gohtml:17\` | \`$unit.Bedrooms\` | \`GridUnit.Bedrooms\` exists (int16) |
| \`calendar-grid.gohtml:17\` | \`$unit.City\` | \`GridUnit.City\` exists (string) |
| \`calendar-grid.gohtml:19\` | \`$unit.Name\` | same |
| \`calendar-grid.gohtml:20\` | \`$unit.Bedrooms\`, \`$unit.City\` | same |
| \`calendar-grid.gohtml:25\` | \`$unit.ID\` | \`GridUnit.ID\` exists (int64) |
| \`calendar.gohtml:52\` | \`$.Grid.Dates\` | \`GridData.Dates\` exists (\`[]DateHeader\`) |

### Category 3: "expected 'string' but got 'invalid type'" for struct fields through map indirection

The linter loses type info when a struct is stored in \`map[string]any\` and then accessed:

| File | Expression | Reality |
|------|-----------|---------|
| \`units.gohtml:29\` | \`eq $u.View "table"\` | \`$u\` is \`unitsPageData\`, \`.View\` is \`string\` |
| \`units.gohtml:35\` | \`eq $u.View "map"\` | same |
| \`units.gohtml:71\` | \`eq (printf "%d" .) $u.FilterBedrooms\` | \`.FilterBedrooms\` is \`string\` |
| \`units.gohtml:85\` | \`eq $u.FilterConfidence "low"\` | \`.FilterConfidence\` is \`string\` |
| \`units.gohtml:86\` | \`eq $u.FilterConfidence "medium"\` | same |
| \`units.gohtml:87\` | \`eq $u.FilterConfidence "high"\` | same |
| \`units.gohtml:102\` | \`eq $u.View "map"\` | same as line 29 |
| \`user-link.gohtml:6\` | \`$.PictureURL\` | User struct has PictureURL |

## Possible solutions

- [ ] Support \`map[string]any\` as template data — treat map keys as valid field names
- [ ] Add type annotation comments (\`{{/* @type MyStruct */}}\`) so templates can declare their expected data type
- [ ] Add a \`//gotmpls:ignore\` or \`{{/* gotmpls:ignore */}}\` directive to suppress specific errors
- [ ] Read Go source to discover the actual types passed to \`ExecuteTemplate\` calls (ambitious but ideal)

## Impact

In our codebase, **100% of the 27 check errors are false positives**, making the check output unusable for CI. We cannot adopt \`check\` in CI until false positive rate drops to near zero.