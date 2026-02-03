---
# golsp-hyo6
title: Fix inline template blocks corrupting blockStack
status: in-progress
type: bug
created_at: 2026-02-03T00:03:09Z
updated_at: 2026-02-03T00:03:09Z
---

When {{if ...}}<p>...</p>{{end}} appears on a single line, ends pop from blockStack before opens push, consuming an outer block's entry and leaving an orphaned entry. Fix by cancelling matched open/end pairs on the same line before stack processing.