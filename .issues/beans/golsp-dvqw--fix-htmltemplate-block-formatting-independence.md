---
# golsp-dvqw
title: Fix HTML/template block formatting independence
status: completed
type: bug
priority: normal
created_at: 2026-02-02T23:29:35Z
updated_at: 2026-02-02T23:45:17Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:25Z"
        task_id: 868hk17k2
---

The formatter's single level variable drifts when HTML tags appear in mutually exclusive template branches (if/else). Implement save/restore stack for template block boundaries.
