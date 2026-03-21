---
# b37-bk1
title: Migrate to toba/lsp server harness
status: completed
type: task
priority: high
created_at: 2026-03-21T17:45:43Z
updated_at: 2026-03-21T18:08:06Z
sync:
    clickup:
        synced_at: "2026-03-21T17:46:00Z"
        task_id: 868hz71tu
---

Replace the hand-rolled LSP main loop with the new `github.com/toba/lsp/server` package (v0.2.0+).

## What changes

The `server` package handles all lifecycle boilerplate:
- JSON-RPC transport via `go.lsp.dev/jsonrpc2`
- `initialize` / `initialized` / `shutdown` / `exit` lifecycle
- Document state management (open/change/close)
- Diagnostic publishing with debouncing
- Optional handler delegation (Hover, Completion, Definition, Formatting, CodeAction, References, Rename, DocumentSymbol)

## Steps

- [ ] Add `github.com/toba/lsp v0.2.0` dependency
- [ ] Implement `server.Handler` interface (Initialize, Diagnostics, Shutdown)
- [ ] Implement any optional handler interfaces (e.g. `server.HoverHandler`, `server.CompletionHandler`)
- [ ] Replace main loop with `server.Server{Name: "go-template-lsp", Version: version, Handler: h}.Run(ctx)`
- [ ] Remove hand-rolled JSON-RPC dispatch, document store, and diagnostic goroutine
- [ ] Remove direct dependency on `toba/lsp/transport` if no longer needed
- [ ] Run tests and linter
- [ ] Verify in editor (VS Code or Zed)

## Summary of Changes

- Created `handler.go` implementing `server.Handler` (Initialize, Diagnostics, Shutdown) plus optional interfaces: `server.HoverHandler`, `server.DefinitionHandler`, `server.FormattingHandler`
- Replaced hand-rolled main loop with `server.Server{Name, Version, Handler}.Run(ctx)`
- Workspace state (parsed files, analyzed files, errors) now managed inside the handler struct with mutex synchronization
- Cross-file semantic analysis runs within `Diagnostics()` callback on each document change
- Server harness handles: LSP lifecycle (initialize/shutdown/exit), document sync (didOpen/didChange/didClose), diagnostic publishing with debouncing, and JSON-RPC transport
- `check` subcommand preserved unchanged
- All existing tests pass
