---
# 7j3-htl
title: go-template-lsp check crashes with nil pointer dereference
status: completed
type: bug
priority: normal
created_at: 2026-03-16T19:40:03Z
updated_at: 2026-03-16T20:02:39Z
sync:
    clickup:
        synced_at: "2026-03-16T20:03:20Z"
        task_id: 868hwye6g
---

Running `go-template-lsp check web/` from the pacer/core repo produces:\n\n```\ngo-template-lsp check: internal error: runtime error: invalid memory address or nil pointer dereference\n```\n\nThis blocks the pacer/core pre-push lint hook (lint.sh runs go-template-lsp check).

## Summary of Changes

The nil pointer dereference was caused by two conditions in `analyzer_statements.go` (around line 196 and 228) where `rhs` (`inferences.uniqueVariableInExpression`) was dereferenced **before** its nil check. When a `{{with}}` expression used a function call returning `any` (e.g. `{{with index . "key"}}`), `rhs` was nil, causing the crash.

**Fix**: moved `rhs != nil` to the front of both conditional chains so it short-circuits before any dereference.

**Files changed**:
- `internal/template/analyzer/analyzer_statements.go` — reordered nil checks (two locations)
- `cmd/go-template-lsp/check_test.go` — added regression test
- `cmd/go-template-lsp/testdata/check/with-nil-rhs/page.gohtml` — test fixture
