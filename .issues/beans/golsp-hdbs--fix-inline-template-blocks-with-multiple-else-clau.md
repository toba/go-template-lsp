---
# golsp-hdbs
title: Fix inline template blocks with multiple else clauses breaking indentation
status: completed
type: bug
priority: normal
created_at: 2026-02-03T01:20:37Z
updated_at: 2026-02-03T01:24:33Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:25Z"
        task_id: 868hk17jm
---

When an HTML attribute contains a template action with multiple else clauses (e.g., if/else if/else/end), the inline cancellation logic fails to cancel all the elses, causing indentation reset to 0.
