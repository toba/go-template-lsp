---
name: upstream
description: |
  Check upstream repos for new changes that may be worth incorporating. Use when:
  (1) User says /upstream
  (2) User asks to "check upstream" or "what changed upstream"
  (3) User wants to know if upstream repos have new commits
  (4) User asks about syncing with or pulling from upstream sources
---

# Upstream Change Tracker

Check the two upstream repos that go-template-lsp derives code from for new commits, classify changes by relevance, and present a summary.

## Upstream Repos

| Repo | Default Branch | Derived Into |
|------|---------------|--------------|
| `yayolande/gota` | `master` | `internal/template/` — parser, lexer, analyzer, template logic |
| `yayolande/go-template-lsp` | `master` | `cmd/go-template-lsp/` — LSP server architecture, methods, parsing |

## Workflow

### Step 1: Read Marker File

Read `.claude/skills/upstream/references/last-checked.json`.

- **If the file does not exist** → this is a first run. Set `FIRST_RUN=true`.
- **If the file exists** → parse the JSON to get `last_checked_sha` and `last_checked_date` per repo.

### Step 2: Fetch Changes (Both Repos in Parallel)

Run both `gh api` calls in parallel using the Bash tool.

#### First Run (no marker file)

Fetch the last 30 commits per repo:

```bash
gh api "repos/yayolande/gota/commits?per_page=30&sha=master" --jq '[.[] | {sha: .sha, date: .commit.committer.date, message: (.commit.message | split("\n") | .[0]), author: .commit.author.name}]'
```

```bash
gh api "repos/yayolande/go-template-lsp/commits?per_page=30&sha=master" --jq '[.[] | {sha: .sha, date: .commit.committer.date, message: (.commit.message | split("\n") | .[0]), author: .commit.author.name}]'
```

Also fetch the changed files for each repo's recent commits to classify relevance:

```bash
gh api "repos/yayolande/gota/commits?per_page=30&sha=master" --jq '[.[].sha]' | jq -r '.[]' | head -30 | while read sha; do gh api "repos/yayolande/gota/commits/$sha" --jq '{sha: .sha, files: [.files[].filename]}'; done
```

(Repeat for go-template-lsp.)

#### Subsequent Runs (marker file exists)

Use the compare API:

```bash
gh api "repos/yayolande/gota/compare/{LAST_SHA}...master" --jq '{total_commits: .total_commits, commits: [.commits[] | {sha: .sha, date: .commit.committer.date, message: (.commit.message | split("\n") | .[0]), author: .commit.author.name}], files: [.files[].filename]}'
```

```bash
gh api "repos/yayolande/go-template-lsp/compare/{LAST_SHA}...master" --jq '{total_commits: .total_commits, commits: [.commits[] | {sha: .sha, date: .commit.committer.date, message: (.commit.message | split("\n") | .[0]), author: .commit.author.name}], files: [.files[].filename]}'
```

**Fallback:** If the compare API returns 404 (e.g. force-push rewrote history), fall back to date-based query:

```bash
gh api "repos/{owner}/{repo}/commits?since={LAST_DATE}&sha=master&per_page=100" --jq '[.[] | {sha: .sha, date: .commit.committer.date, message: (.commit.message | split("\n") | .[0]), author: .commit.author.name}]'
```

### Step 3: Classify Changed Files by Relevance

Use these mappings to assign HIGH / MEDIUM / LOW relevance to each changed file:

#### yayolande/gota

| Relevance | Path Patterns |
|-----------|--------------|
| **HIGH** | `analyzer.go`, `parser/**`, `lexer/**`, `template.go` |
| **MEDIUM** | `go.mod`, `go.sum`, `*_test.go` |
| **LOW** | `.github/**`, `README.md`, `LICENSE` |

#### yayolande/go-template-lsp

| Relevance | Path Patterns |
|-----------|--------------|
| **HIGH** | `lsp/**`, `main.go` |
| **MEDIUM** | `go.mod`, `go.sum`, `*_test.go` |
| **LOW** | `.github/**`, `README.md`, `LICENSE`, `.gitignore` |

Files not matching any pattern → **MEDIUM** (unknown = worth reviewing).

### Step 4: Present Summary

Format the output as follows:

```
# Upstream Changes

## yayolande/gota (N new commits since YYYY-MM-DD)

### Commits
- `abc1234` Fix parser edge case — @author (2025-05-01)
- `def5678` Add range variable support — @author (2025-04-28)

### Changed Files

**HIGH relevance** (directly mapped into our template internals):
- analyzer.go
- parser/parse.go

**MEDIUM relevance** (may affect behavior):
- go.mod

**LOW relevance** (infrastructure/docs):
- README.md

**Assessment:** 2 high-relevance changes to parser — worth reviewing for potential incorporation.

---

(repeat for go-template-lsp)

---

## Overall Recommendation
(Summarize: how many repos have high-relevance changes, suggest priority order for review)
```

If a repo has **no new commits**, show:

```
## repo/name — No new commits since last check (YYYY-MM-DD)
```

### Step 5: Update Marker File

Build the new marker JSON with the HEAD SHA and current date for each repo.

- **First run:** Write the marker file automatically (tell the user it was created).
- **Subsequent runs:** Ask the user "Update the last-checked markers to current HEAD?" before writing.

Write to `.claude/skills/upstream/references/last-checked.json`:

```json
{
  "yayolande/gota": {
    "last_checked_sha": "<HEAD_SHA>",
    "last_checked_date": "<ISO_DATE>"
  },
  "yayolande/go-template-lsp": {
    "last_checked_sha": "<HEAD_SHA>",
    "last_checked_date": "<ISO_DATE>"
  }
}
```
