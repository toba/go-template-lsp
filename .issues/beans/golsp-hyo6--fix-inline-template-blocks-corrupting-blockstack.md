---
# golsp-hyo6
title: Fix inline template blocks corrupting blockStack
status: completed
type: bug
priority: normal
created_at: 2026-02-03T00:03:09Z
updated_at: 2026-02-08T17:24:04Z
---

When {{if ...}}<p>...</p>{{end}} appears on a single line, ends pop from blockStack before opens push, consuming an outer block's entry and leaving an orphaned entry. Fix by cancelling matched open/end pairs on the same line before stack processing.