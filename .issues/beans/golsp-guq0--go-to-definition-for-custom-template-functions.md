---
# golsp-guq0
title: Go-to-definition for custom template functions
status: completed
type: feature
priority: normal
created_at: 2026-02-08T16:46:38Z
updated_at: 2026-02-08T16:50:16Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:25Z"
        task_id: 868hk17k7
---

Custom template functions should be cmd-clickable, navigating to their definition in Go source. Two bugs: (1) No source position captured in FuncMap scanner (2) Wrong range returned in FindSourceDefinitionFromPosition.

## Checklist
- [ ] Capture source positions in FuncMap scanner (funcmap_scanner.go)
- [ ] Return function definition directly for go-to-definition (analyzer_lsp.go)
- [ ] Add tests for position capture (funcmap_scanner_test.go)
- [ ] Add integration test for go-to-definition with custom functions (template_test.go)
- [ ] Run tests and lint
