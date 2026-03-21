---
# ei0-agl
title: Adopt shared github.com/toba/lsp module
status: completed
type: feature
priority: normal
created_at: 2026-03-21T17:18:32Z
updated_at: 2026-03-21T18:07:55Z
sync:
    clickup:
        synced_at: "2026-03-21T17:46:00Z"
        task_id: 868hz71tt
---

Replace duplicated LSP infrastructure code with the shared github.com/toba/lsp module. This eliminates ~500-700 lines of copy-pasted code that is identical across all four toba LSP projects.

## Packages to adopt

- `github.com/toba/lsp/transport` — replaces hand-rolled `parsing.go` (ReceiveInput, decode, Encode, SendToLspClient, SendOutput)
- `github.com/toba/lsp/logging` — replaces createLogFile(), configureLogging(), MaxLogFileSize/DirPermissions/FilePermissions constants
- `github.com/toba/lsp/pathutil` — replaces uriToFilePath(), filePathToUri(), convertKeysFromFilePathToUri()
- `github.com/toba/lsp/position` — replaces OffsetToLineChar/offsetToLineCol, LineCharToOffset, intToUint, uintToInt

## Steps

1. Add `github.com/toba/lsp` dependency
2. Replace transport layer: `lsp.ReceiveInput` → `transport.NewScanner`, `lsp.SendToLspClient` → `transport.Send`, `lsp.Encode` → `transport.Encode`
3. Replace logging: `createLogFile()`/`configureLogging()` → `logging.Configure(appName)`
4. Replace path utilities: `uriToFilePath` → `pathutil.URIToFilePath`, `filePathToUri` → `pathutil.FilePathToURI`
5. Replace position utilities: `OffsetToLineChar` → `position.OffsetToLineCol`, etc.
6. Delete the replaced local code (parsing.go, logging functions, URI functions, position functions)
7. Remove now-unused constants (ContentLengthHeader, HeaderDelimiter, etc.)
8. Run tests and linter

## Summary of Changes

- Added `github.com/toba/lsp` v0.2.1 as a dependency
- Replaced local `uriToFilePath`/`filePathToUri`/`convertKeysFromFilePathToUri` with `pathutil.URIToFilePath`/`pathutil.FilePathToURI`/`pathutil.ConvertMapKeysToURI`
- Replaced local `createLogFile`/`configureLogging` with `logging.Configure` (via server harness)
- Replaced local LSP transport (`parsing.go`: `ReceiveInput`, `SendToLspClient`, `Encode`, `decode`) with server harness transport
- Deleted entire `cmd/go-template-lsp/lsp/` package (protocol.go, methods.go, parsing.go and tests)
- Removed local constants `DirPermissions`, `FilePermissions`, `MaxLogFileSize`, `ContentLengthHeader`, etc.
