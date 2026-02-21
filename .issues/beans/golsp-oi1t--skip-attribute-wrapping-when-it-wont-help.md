---
# golsp-oi1t
title: Skip attribute wrapping when it won't help
status: completed
type: bug
priority: normal
created_at: 2026-02-01T19:40:46Z
updated_at: 2026-02-01T19:41:29Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:25Z"
        task_id: 868hk17jh
---

When an element has attributes that exceed printWidth, the formatter wraps them onto new lines. But if the wrapped lines STILL exceed printWidth (e.g., a single long SVG d attribute), wrapping is pointless and makes the output worse.
