---
# golsp-34xz
title: Support custom FuncMap declarations for check/diagnostics
status: completed
type: feature
priority: normal
created_at: 2026-02-08T17:19:41Z
updated_at: 2026-02-08T17:24:24Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:25Z"
        task_id: 868hk17jt
---

The `check` subcommand (and LSP diagnostics) report false-positive **"function undefined"** and **"only functions and methods accept arguments"** errors for custom Go template FuncMap functions (e.g. `asset`, `timehtml`, custom helpers).

## Problem

Projects using `html/template` with custom FuncMaps get hundreds of spurious diagnostics. In one real codebase (Pacer), 487 of ~550 lint errors are these false positives, making the tool unusable for CI or editor integration.

Example template:
```html
<link rel="stylesheet" href="{{asset "css/base.css"}}">
```
Produces:
```
app.gohtml:9:36: function undefined
app.gohtml:9:36: only functions and methods accept arguments
```

## Proposed Solution

Allow declaring custom template functions via a config file (e.g. `.gotmpls.yaml`, `.gotmpls.json`, or `gotmpls.config` in project root), so the linter knows they exist and can skip "function undefined" errors for them.

Minimal example:
```yaml
funcs:
  - asset
  - timehtml
  - route
  - dict
```

Ideally, the linter could also auto-discover FuncMap declarations by scanning Go source files for `template.FuncMap{...}` patterns, but a static config file would solve the immediate problem.

## Impact

- **487 false positives eliminated** in one real-world project
- Makes `check -json` viable for CI pipelines
- Makes LSP diagnostics useful instead of noisy

## Related

There's a secondary class of false positives from type inference on range variables (linter infers `any` and can't resolve field access), which `/*gotype:*/` annotations partially address. But the FuncMap issue is the bigger blocker — it affects every project with custom template functions.

## Breakdown

- [x] Auto-discover FuncMap from Go source files (scans from module root)
- [x] Add `-funcs` flag to `check` subcommand for additional functions
- [x] LSP hot-reload: re-scans when Go files change
- [x] Go-to-definition support for custom functions
- [x] Update README with documentation
