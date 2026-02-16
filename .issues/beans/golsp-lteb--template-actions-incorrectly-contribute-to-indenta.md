---
# golsp-lteb
title: Template actions incorrectly contribute to indentation level
status: completed
type: bug
priority: normal
created_at: 2026-02-01T18:40:44Z
updated_at: 2026-02-01T18:43:13Z
---

The formatter compounds indentation from both HTML tags and template actions (define/range/if/with/block/end). Template actions should NOT contribute to the indent level — only HTML nesting should drive indentation. Template actions should be indented at the current HTML-determined level but not push it further.

## Example

Input:
```
{{define "client-list" -}}
<ul class="entity">
    {{range . -}}
<li><a {{itemLinkAttrs . 0}}><h3>{{.Name}}</h3></a></li>
    {{- end -}}
</ul>
{{- end -}}
```

Expected:
```
{{define "client-list" -}}
<ul class="entity">
   {{range . -}}
   <li><a {{itemLinkAttrs . 0}}><h3>{{.Name}}</h3></a></li>
   {{- end -}}
</ul>
{{- end -}}
```

## Checklist

- [ ] Add failing test case with the above input/expected
- [ ] Update computeLineDeltas to ignore template action events
- [ ] Update all existing tests to reflect new indentation behavior
- [ ] Run tests and golangci-lint