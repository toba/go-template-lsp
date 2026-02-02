---
# golsp-ijrh
title: Fix concurrent stdout writes and add protection tests
status: completed
type: bug
priority: high
created_at: 2026-02-02T16:51:37Z
updated_at: 2026-02-02T16:57:40Z
---

The main goroutine and diagnostic goroutine both call `lsp.SendToLspClient(os.Stdout, ...)` without synchronization, which can corrupt the LSP Content-Length framing and cause the editor to report broken pipe errors.

## Checklist

- [ ] Add `muStdout := new(sync.Mutex)` in `main()` alongside `muTextFromClient`
- [ ] Pass `muStdout` to `processDiagnosticNotification`
- [ ] Wrap the 3 `lsp.SendToLspClient(os.Stdout, ...)` call sites with `muStdout.Lock()`/`Unlock()`:
  - `main.go:127` (shutdown/illegal-request path)
  - `main.go:243` (normal response path)
  - `main.go:495` (diagnostic notification goroutine)
- [ ] Add race detector test: spawn concurrent goroutines writing to a shared writer with mutex, verify no race and all messages decode
- [ ] Add frame integrity test: concurrent writes through io.Pipe, verify every decoded message has valid Content-Length framing and well-formed JSON body
- [ ] Run `go test -race ./...` and `golangci-lint run --fix`

## Reference

See go-css-lsp commit 164e103 and 0152761 for the mutex fix and test implementations.