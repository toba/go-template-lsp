# Go Template LSP

A Language Server Protocol (LSP) implementation for Go templates (`text/template` and `html/template`).

## Features

- **Diagnostics**: Real-time syntax error detection and semantic analysis (undefined variables, missing fields, unknown functions)
- **Hover**: Type information and documentation on hover over template variables and functions
- **Go to Definition**: Navigate to template definitions (`{{define "name"}}`)
- **Formatting**: Re-indent based on HTML and template nesting, with optional attribute wrapping
- **Folding Ranges**: Collapse template blocks (`{{if}}...{{end}}`, `{{range}}...{{end}}`) and comments
- **Semantic Tokens**: Enhanced syntax highlighting for keywords, variables, functions, strings, numbers, and operators
- **Document Highlight**: Highlight matching template keywords (e.g. click `{{if}}` to highlight its `{{else}}` and `{{end}}`)
- **Document Links**: Clickable links in template documents
- **Custom Function Discovery**: Automatically scans Go source files for `template.FuncMap` definitions so custom functions are recognized
- **CLI Linting**: `check` subcommand for CI/CD and command-line diagnostics

## Installation

### Download Binary

Download prebuilt binaries from [GitHub Releases](https://github.com/toba/go-template-lsp/releases).

### Install from Source

```bash
go install github.com/toba/go-template-lsp/cmd/go-template-lsp@latest
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/toba/go-template-lsp.git
cd go-template-lsp

# Build the binary
go build -o go-template-lsp ./cmd/go-template-lsp

# Or install to $GOBIN
go install ./cmd/go-template-lsp

# Run tests
go test ./...

# Run linter
golangci-lint run
```

## Editor Integration

### Zed

Use the [gozer](https://github.com/toba/gozer) Zed extension, which automatically downloads this LSP.

### Neovim (with nvim-lspconfig)

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

configs.gotmpl = {
  default_config = {
    cmd = { 'go-template-lsp' },
    filetypes = { 'gotmpl', 'gohtml', 'html' },
    root_dir = lspconfig.util.root_pattern('go.mod', '.git'),
  },
}

lspconfig.gotmpl.setup{}
```

### Other Editors

The LSP binary works with any editor that supports the Language Server Protocol. Configure your editor to run `go-template-lsp` for template files.

## Formatting

The formatter re-indents Go template files based on HTML tag and template action nesting. It respects the standard LSP `tabSize` and `insertSpaces` settings.

### Attribute Wrapping

When an HTML opening tag exceeds a configured line width, the formatter can wrap its attributes across multiple lines. Configure via `initializationOptions`:

```json
{
  "initializationOptions": {
    "printWidth": 120,
    "attrWrapMode": "overflow"
  }
}
```

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `printWidth` | `int` | `0` (disabled) | Maximum line width before wrapping attributes. Set to `0` to disable. |
| `attrWrapMode` | `string` | `"overflow"` | `"overflow"`: only wrap attributes that push past `printWidth`. `"all"`: wrap every attribute onto its own line. |

**`overflow` mode** keeps attributes on the first line as long as they fit, then wraps the rest:

```html
<button type="button"
   class="map-zoom-btn"
   title="Expand">
```

**`all` mode** puts every attribute on its own line:

```html
<button
   type="button"
   class="map-zoom-btn"
   title="Expand">
```

Continuation lines are indented one level deeper than the tag. Multi-line tags already in the source are joined back into a single line before re-wrapping.

## CLI Linting (`check`)

The `check` subcommand runs the same diagnostic pipeline as the LSP server and prints errors to stdout. This is useful for CI/CD pipelines, pre-commit hooks, and AI agents.

```bash
# Check the current directory
go-template-lsp check

# Check a specific directory
go-template-lsp check ./views

# JSON output
go-template-lsp check -json ./views

# Declare extra template function names (comma-separated)
go-template-lsp check -funcs asset,renderPage ./views
```

Custom template functions are auto-discovered by scanning Go source files for `template.FuncMap` definitions. The scan starts from the Go module root (`go.mod` directory), so FuncMap definitions anywhere in the module are found even when checking a subdirectory. The `-funcs` flag lets you declare additional function names that can't be auto-discovered (e.g. from external packages or generated code).

**Text output** (default):

```
views/user.html:8:22: missing matching '{{ end }}' statement
```

**JSON output** (`-json`):

```json
[
  {
    "file": "views/user.html",
    "line": 8,
    "column": 22,
    "endLine": 8,
    "endColumn": 24,
    "message": "missing matching '{{ end }}' statement"
  }
]
```

**Exit codes:**

| Code | Meaning |
|------|---------|
| `0` | No errors found |
| `1` | Template errors found |
| `2` | Tool failure (e.g. directory not found) |

## Supported File Extensions

| Extension | Type |
|-----------|------|
| `.gotmpl`, `.go.tmpl`, `.gtpl`, `.tpl`, `.tmpl` | Go text templates |
| `.gohtml`, `.go.html` | Go HTML templates |
| `.html` | HTML (with template detection) |

## Credits

### Template Parser and Analyzer

**[yayolande/gota](https://github.com/yayolande/gota)** (MIT License)

The template parsing and semantic analysis code in `internal/template` is derived from gota by yayolande.

### LSP Implementation

**[yayolande/go-template-lsp](https://github.com/yayolande/go-template-lsp)** (MIT License)

The LSP server architecture is based on this project by yayolande.

## License

MIT License - see [LICENSE](LICENSE) for details.
