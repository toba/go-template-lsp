---
# golsp-94xu
title: Fix false positive type errors for map[string]any template data
status: completed
type: bug
priority: normal
created_at: 2026-02-08T18:13:47Z
updated_at: 2026-02-08T18:17:42Z
---

The check subcommand reports false positives on projects using map[string]any as template data. Two bugs: (1) map field access is blocked in getRealTypeAssociatedToVariable, (2) strict implicit type compatibility requires ALL fields to exist.