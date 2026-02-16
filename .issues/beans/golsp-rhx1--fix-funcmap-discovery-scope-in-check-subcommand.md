---
# golsp-rhx1
title: Fix FuncMap discovery scope in check subcommand
status: completed
type: bug
priority: normal
created_at: 2026-02-08T17:42:40Z
updated_at: 2026-02-08T17:45:51Z
---

The check subcommand scans for FuncMap definitions in the same directory it scans for templates. When templates live in a subdirectory but FuncMap definitions are elsewhere in the module, the scanner misses them — producing false 'function undefined' errors. Fix: scan for FuncMaps at the Go module root (walk up to find go.mod), not at the template directory. Also add -funcs flag as escape hatch.

## Checklist
- [ ] Add FindModuleRoot to funcmap_scanner.go
- [ ] Add FunctionsFromNames to funcmap_scanner.go
- [ ] Add tests for FindModuleRoot and FunctionsFromNames
- [ ] Modify check.go to use FindModuleRoot for FuncMap scanning
- [ ] Add -funcs flag to check subcommand in main.go
- [ ] Wire -funcs flag into runCheck
- [ ] Update README
- [ ] Run tests and lint