---
# golsp-whzj
title: Fix inline template blocks corrupting blockStack
status: completed
type: bug
priority: normal
created_at: 2026-02-03T00:05:00Z
updated_at: 2026-02-03T00:06:04Z
---

When {{if ...}}...{{end}} appears on a single line mixed with HTML, the formatter mis-indents subsequent lines. The fix cancels matching open/end pairs on the same line before stack processing.