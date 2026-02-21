---
# golsp-mk5u
title: 'Format-on-save race: didChange not processed before formatting request'
status: completed
type: bug
priority: high
created_at: 2026-02-01T20:42:55Z
updated_at: 2026-02-01T20:49:56Z
sync:
    clickup:
        synced_at: "2026-02-21T04:27:25Z"
        task_id: 868hk17jy
---

## Problem

When using Zed with autosave (1s delay) and manual Cmd+S, a race condition can cause the LSP to format stale content, wiping recent edits. The sequence:

1. User types something
2. Zed sends `didChange` to the LSP
3. User presses Cmd+S, triggering format-on-save
4. Zed sends `textDocument/formatting` request
5. If the formatting request is processed before the `didChange`, the LSP formats old content from `textFromClient`
6. The formatted (stale) result replaces the buffer, losing recent edits

This is distinct from golsp-bue2 (which addressed the case where `textFromClient` had no entry at all). Here, an entry exists but is outdated because `didChange` hasn't been processed yet.

## Possible Causes

- The LSP processes messages on separate goroutines and `didChange` updates to `textFromClient` race with the formatting handler reading it
- No synchronization ensures `didChange` is fully applied before a formatting request reads the document state

## Likely Fix

Ensure `textFromClient` reads/writes are properly synchronized, or have the formatting handler use the document version to confirm it has the latest content before proceeding.

## Related

- golsp-bue2: Format returns null when file not in textFromClient (completed — addressed the nil case, not the stale case)
