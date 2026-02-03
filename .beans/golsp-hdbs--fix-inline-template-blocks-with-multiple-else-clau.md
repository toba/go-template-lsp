---
# golsp-hdbs
title: Fix inline template blocks with multiple else clauses breaking indentation
status: completed
type: bug
priority: normal
created_at: 2026-02-03T01:20:37Z
updated_at: 2026-02-03T01:24:33Z
---

When an HTML attribute contains a template action with multiple else clauses (e.g., if/else if/else/end), the inline cancellation logic fails to cancel all the elses, causing indentation reset to 0.