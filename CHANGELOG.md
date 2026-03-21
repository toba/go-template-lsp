# Changelog

## Week of Mar 15 – Mar 21, 2026

### ✨ Features

- Adopt shared `github.com/toba/lsp` module; replace local transport, logging, and path utilities

### 🐛 Fixes

- Fix `go-template-lsp check` nil pointer dereference on `{{with index . "key"}}`

### 🗜️ Tweaks

- Migrate to `toba/lsp` server harness; replace hand-rolled main loop with `server.Server.Run`

## Week of Feb 8 – Feb 14, 2026

### ✨ Features

- Add `check` subcommand for CLI linting
- Implement `textDocument/documentLink` for hyphenated template names
- Support custom FuncMap declarations for diagnostics; auto-discover from Go source
- Go-to-definition for custom template functions

### 🐞 Fixes

- Fix FuncMap discovery scope in `check` subcommand; scan from module root
- Fix false positive type errors for `map[string]any` template data
- Fix false positives for dollar variables assigned from any-typed expressions
- Fix false positives for partially-inferred struct compatibility
- Fix `documentLinkProvider` serializing as boolean instead of object

### 🗜️ Tweaks

- Upgrade to Go 1.26; adopt `wg.Go()`, `errors.AsType`, extract `resolveFileNode` helper
- Replace hand-rolled insertion sorts with `slices.SortFunc`
- Update READMEs for `check` subcommand

## Week of Feb 1 – Feb 7, 2026

### ✨ Features

- Add attribute wrapping with `printWidth` and `attrWrapMode` support
- Keep first attribute on tag line in `all` wrap mode when aligned with continuation indent
- Hot-reload custom template functions when Go files change

### 🐞 Fixes

- Fix template actions incorrectly contributing to indentation level
- Fix mixed tabs/spaces and void element indentation
- Format returns null when file not in `textFromClient`; fall back to reading from disk
- Fix format-on-save race; `didChange` not processed before formatting request
- Fix concurrent stdout writes; add mutex protection
- Fix HTML/template block formatting independence; save/restore stack for template branches
- Fix false positive for negative number literals in lexer
- Fix inline template blocks corrupting `blockStack`
- Fix inline template blocks with multiple `else` clauses breaking indentation
- Fix `attrTokenRe` splitting attributes with template actions containing quotes
- Skip attribute wrapping when it won't help; don't wrap attributes that still exceed `printWidth`
- Pack attributes on continuation lines in `overflow` wrap mode

## Week of Jan 26 – Feb 1, 2026

### 🐞 Fixes

- Fix function names with keyword prefixes incorrectly flagged

## Week of Jan 19 – Jan 25, 2026

### 🐞 Fixes

- Fix multiline template comments with nested delimiters
