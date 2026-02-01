---
# golsp-bue2
title: Format returns null when file not in textFromClient
status: completed
type: bug
priority: high
created_at: 2026-02-01T19:25:45Z
updated_at: 2026-02-01T19:27:36Z
---

## Problem

When Zed sends a `textDocument/formatting` request for a file that hasn't received a `didOpen` notification, the LSP returns `{"result": null}` instead of formatted content. This means format-on-save silently does nothing.

## Root Cause

`ProcessFormattingRequest` in `cmd/go-template-lsp/lsp/methods.go:1001` checks `textFromClient[fileUri]` and if nil, returns an empty response. The `textFromClient` map is only populated via `didOpen`/`didChange`, but Zed only sends `didOpen` for files actively focused in the editor — not all files it considers open.

## Fix

When `fileContent == nil`, fall back to reading the file from disk using the URI (strip `file://` prefix and `os.ReadFile`). The `fileUri` is already a `file://` URI like `file:///Users/jason/Developer/pacer/core/web/components/sidebar/sidebar.gohtml`.

Around line 1001:

```go
if fileContent == nil {
    // Fall back to reading from disk
    path := strings.TrimPrefix(fileUri, "file://")
    fileContent, err = os.ReadFile(path)
    if err != nil {
        slog.Warn("Cannot read file for formatting: " + err.Error())
        responseData, err := json.Marshal(res)
        if err != nil {
            slog.Warn("Error marshalling formatting response: " + err.Error())
            return nil, fileName
        }
        return responseData, fileName
    }
}
```

Requires adding `"os"` and `"strings"` to imports (strings may already be imported).