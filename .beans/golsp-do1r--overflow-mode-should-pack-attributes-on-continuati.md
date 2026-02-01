---
# golsp-do1r
title: Overflow mode should pack attributes on continuation lines
status: completed
type: bug
priority: normal
created_at: 2026-02-01T19:42:49Z
updated_at: 2026-02-01T19:45:17Z
---

In overflow attrWrapMode, once an attribute overflows the printWidth, all remaining attributes are placed one-per-line. Instead, remaining attributes should be packed onto continuation lines, fitting as many as possible per line before wrapping again.

## Expected
```
<section class="portfolio" hx-boost:inherited="true" hx-target:inherited="main"
      hx-swap:inherited="innerHTML" hx-push-url:inherited="true">
```

## Actual
```
<section class="portfolio" hx-boost:inherited="true" hx-target:inherited="main"
      hx-swap:inherited="innerHTML"
      hx-push-url:inherited="true">
```

## Checklist
- [ ] Add failing test demonstrating the packing behavior
- [ ] Fix overflow mode to pack attributes on continuation lines
- [ ] Run tests and lint