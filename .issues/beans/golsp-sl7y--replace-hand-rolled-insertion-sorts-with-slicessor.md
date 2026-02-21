---
# golsp-sl7y
title: Replace hand-rolled insertion sorts with slices.SortFunc
status: completed
type: task
priority: normal
created_at: 2026-02-08T17:00:28Z
updated_at: 2026-02-08T17:02:39Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:26Z"
        task_id: 868hk17ka
---

Four insertion sort implementations should use slices.SortFunc instead:

## Checklist
- [ ] format.go:88 — tagBracket sort by pos
- [ ] format.go:482 — tagEvent sort by pos
- [ ] template.go:840 — sortSemanticTokens by line then char
- [ ] check.go:156 — sort.Slice diagnostics → slices.SortFunc with cmp
