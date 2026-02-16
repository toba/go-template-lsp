---
# golsp-vzjn
title: Fix false positive 'character(s) not recognized' for negative numbers
status: completed
type: bug
priority: normal
created_at: 2026-02-02T23:50:54Z
updated_at: 2026-02-02T23:51:35Z
---

The lexer reports errors on negative number literals like -1 in template expressions. Fix by extending numeric regex patterns to optionally accept a leading minus sign.