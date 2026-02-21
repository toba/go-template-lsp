---
# 5q5-or7
title: Tokenizer breaks on JSON with embedded Go template expressions in HTML attributes
status: scrapped
type: bug
priority: high
created_at: 2026-02-18T22:39:35Z
updated_at: 2026-02-18T23:02:52Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:25Z"
        task_id: 868hk17k8
---

## Description

The tokenizer produces incorrect tokens when an HTML attribute value contains JSON that itself embeds Go template expressions. This pattern is common with htmx's `hx-vals` attribute.

## Reproducer

From `pacer/core/web/portfolio/portfolio/gohtml`, around line 65:

```html
hx-vals='{"entity_kind":"{{.EntityKind}}","entity_id":"{{.EntityID}}","issue_type":"{{.IssueType}}"}'
```

The JSON curly braces (`{`, `}`) interleave with Go template delimiters (`{{`, `}}`), causing the tokenizer to misparse the region.

## Expected

The tokenizer should correctly identify Go template actions (`{{.EntityKind}}` etc.) within the attribute value and treat the surrounding JSON literal braces as plain text.

## Acceptance Criteria

- [ ] Add a failing test demonstrating incorrect tokenization of JSON-in-attribute with embedded Go template expressions
- [ ] Fix the tokenizer to handle this case correctly
- [ ] Existing tests continue to pass


## Reasons for Scrapping

Not a go-template-lsp bug. The LSP lexer, parser, and semantic tokens all handle this pattern correctly (verified with comprehensive tests). The visual issue is caused by the tree-sitter HTML injection in the gozer Zed extension — the HTML parser can't handle single-quoted attribute values split across disjoint byte ranges by template actions. Moved to gozer as afc-vlg.
