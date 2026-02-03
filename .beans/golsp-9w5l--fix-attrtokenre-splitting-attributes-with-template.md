---
# golsp-9w5l
title: Fix attrTokenRe splitting attributes with template actions containing quotes
status: completed
type: bug
priority: normal
created_at: 2026-02-03T00:56:48Z
updated_at: 2026-02-03T01:01:19Z
---

The attrTokenRe regex uses [^"]* for double-quoted attribute values, which breaks on attributes like aria-sort="{{if eq .SortBy "check_in"}}ascending{{else}}none{{end}}". The regex matches up to the first inner quote, splitting the attribute into garbage tokens and corrupting the formatted output.