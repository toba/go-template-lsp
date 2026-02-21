---
# golsp-33oi
title: Fix mixed tabs/spaces and void element indentation
status: completed
type: bug
priority: normal
created_at: 2026-02-01T18:48:36Z
updated_at: 2026-02-01T18:58:00Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:25Z"
        task_id: 868hk17kc
---

Two formatting issues reported:

1. **Mixed tabs and spaces**: With insertSpaces enabled, the formatter still outputs a mix of tabs and spaces
2. **Void element indentation**: `<img>` (a void element) is not indenting to the correct level

## Reproduction input

```
{{define "sidebar" -}}
<aside id="sidebar">
   <header>
      <a href="/" aria-label="Pacer home">
   <img src="/web/images/logo.svg" alt=""/>
           <span></span>{{/* "pacer" SVG assigned in CSS */}}
       </a>
   </header>
</aside>
{{end}}
```

## Checklist

- [ ] Investigate root cause of mixed tabs/spaces
- [ ] Investigate void element indentation
- [ ] Write failing test(s) demonstrating both issues
- [ ] Fix the formatter code
- [ ] Verify all tests pass and lint is clean
