---
# golsp-f3xq
title: Hot-reload custom template functions when Go files change
status: completed
type: feature
priority: normal
created_at: 2026-02-06T18:31:47Z
updated_at: 2026-02-06T18:37:32Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:25Z"
        task_id: 868hk17jf
---

Use LSP file watching (client/registerCapability + workspace/didChangeWatchedFiles) to re-scan Go files for FuncMap definitions when they change, then re-analyze templates automatically.

## Checklist
- [x] Add constants to protocol.go (MethodRegisterCapability, MethodDidChangeWatchedFiles, FileChangeType, WatchKind)
- [x] Add types to methods.go (Registration, RegistrationParams, FileSystemWatcher, etc.)
- [x] Add BuildRegisterFileWatcherRequest function to methods.go
- [x] Add ProcessDidChangeWatchedFilesNotification function to methods.go
- [x] Add unit tests for BuildRegisterFileWatcherRequest and ProcessDidChangeWatchedFilesNotification
- [x] Add goFilesChangedNotification channel and wire it through main.go
- [x] Send client/registerCapability in MethodInitialized handler
- [x] Handle MethodDidChangeWatchedFiles in main loop
- [x] Add customFunctionsEqual helper
- [x] Modify processDiagnosticNotification to handle goFilesChanged via select
- [x] Add unit test for customFunctionsEqual
- [x] Add integration test for hot-reload behavior
- [x] Run go test ./... and golangci-lint run --fix
