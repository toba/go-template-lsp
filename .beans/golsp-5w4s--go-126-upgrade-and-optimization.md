---
# golsp-5w4s
title: Go 1.26 upgrade and optimization
status: completed
type: task
priority: normal
created_at: 2026-02-12T17:30:22Z
updated_at: 2026-02-12T17:50:57Z
---

Upgrade project to Go 1.26 and apply optimization findings from goptimize analysis.

## Summary of Changes

- Updated Go version from 1.25.6 to 1.26.0 (go.mod + mise.toml)
- Updated golangci-lint from v2.8.0 to v2.9.0 (required for Go 1.26)
- Removed redundant `intrange` and `modernize` linters from .golangci.yaml (now covered by `go fix`)
- Converted `wg.Add(1); go func() { defer wg.Done()` → `wg.Go()` in 3 locations (Go 1.26)
- Converted `errors.As` → `errors.AsType[T]()` in 2 locations (Go 1.26)
- Extracted duplicated file retrieval pattern into `resolveFileNode()` helper in methods.go (4 call sites)