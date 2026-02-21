---
# golsp-k9kf
title: Add check subcommand for CLI linting
status: completed
type: feature
priority: normal
created_at: 2026-02-08T16:28:03Z
updated_at: 2026-02-08T16:32:41Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:25Z"
        task_id: 868hk17ju
---

Add a `check` subcommand that runs the same diagnostic pipeline as the LSP server and outputs errors to stdout for CI/CD and agents.

## Checklist
- [x] Add subcommand dispatch in main.go
- [x] Create check.go with runCheck, collectDiagnostics, output formatting
- [x] Create testdata fixtures (valid and errors directories)
- [x] Create check_test.go with tests
- [x] Run tests and linter
