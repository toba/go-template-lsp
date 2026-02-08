---
# golsp-909x
title: Implement textDocument/documentLink for hyphenated template names
status: completed
type: feature
priority: normal
created_at: 2026-02-08T16:30:43Z
updated_at: 2026-02-08T16:32:10Z
---

Implement the textDocument/documentLink LSP method so that hyphenated template names like "logo-particles" are treated as a single clickable link, overriding the editor's word-boundary detection that splits at hyphens.

## Checklist
- [ ] Add DocumentLinkInfo struct and DocumentLinks function to internal/template/template.go
- [ ] Add MethodDocumentLink constant to lsp/protocol.go
- [ ] Add DocumentLinkParams, DocumentLinkResult types and ProcessDocumentLinkRequest handler to lsp/methods.go
- [ ] Add DocumentLinkProvider capability to ServerCapabilities
- [ ] Wire up handler in cmd/go-template-lsp/main.go
- [ ] Add TestDocumentLinks test to template_test.go
- [ ] Run tests and linter