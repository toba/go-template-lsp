---
# kut-cg8
title: Migrate to go.lsp.dev/protocol + go.lsp.dev/jsonrpc2
status: draft
type: feature
priority: normal
created_at: 2026-03-21T17:10:55Z
updated_at: 2026-03-21T17:10:55Z
sync:
    clickup:
        synced_at: "2026-03-21T17:46:00Z"
        task_id: 868hz71tr
---

Replace the hand-rolled LSP protocol types and JSON-RPC transport in cmd/go-template-lsp/lsp/ with the standard go.lsp.dev/protocol and go.lsp.dev/jsonrpc2 packages. This eliminates maintaining custom LSP struct definitions and transport code, and gives full spec-compliant types for all LSP methods.

## Steps
- Add go.lsp.dev/protocol and go.lsp.dev/jsonrpc2 as dependencies
- Replace all custom types in lsp/protocol.go with imports from go.lsp.dev/protocol
- Replace the custom scanner/transport in lsp/methods.go with go.lsp.dev/jsonrpc2 stream handling
- Update all handler functions to use the standard protocol types
- Remove custom type definitions that are now redundant
- Run tests and linter
