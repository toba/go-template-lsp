---
# golsp-opnu
title: Don't wrap attributes that still exceed printWidth on their own line
status: completed
type: bug
priority: normal
created_at: 2026-02-01T19:37:00Z
updated_at: 2026-02-02T16:55:22Z
---

When an element has an attribute that exceeds printWidth and will STILL exceed printWidth even on its own line, don't bother reformatting. Example with printWidth 100:

```html
<path d="M12 6V2H8M2 12h2m16 0h2m-2 4a2 2 0 0 1-2 2H8.828a2 2 0 0 0-1.414.586l-2.202 2.202A.71.71 0 0 1 4 20.286V8a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2z" />
```

This should remain on one line since wrapping the `d` attribute to its own line won't help it fit within printWidth.