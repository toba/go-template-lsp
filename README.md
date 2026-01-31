# Go Template LSP

A Language Server Protocol (LSP) implementation for Go templates (`text/template` and `html/template`).

## Features

- **Diagnostics**: Real-time syntax error detection as you type
- **Hover**: Type information and documentation on hover over template variables and functions
- **Go to Definition**: Navigate to template definitions (`{{define "name"}}`)
- **Folding Ranges**: Collapse template blocks (`{{if}}...{{end}}`, `{{range}}...{{end}}`) and comments
- **Semantic Tokens**: Enhanced syntax highlighting
- **Document Highlight**: Highlight matching template keywords

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
