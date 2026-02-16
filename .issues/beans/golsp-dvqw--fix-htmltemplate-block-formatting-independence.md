---
# golsp-dvqw
title: Fix HTML/template block formatting independence
status: completed
type: bug
priority: normal
created_at: 2026-02-02T23:29:35Z
updated_at: 2026-02-02T23:45:17Z
---

The formatter's single level variable drifts when HTML tags appear in mutually exclusive template branches (if/else). Implement save/restore stack for template block boundaries.