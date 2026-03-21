package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"net/url"
	"os"
	"strings"
	"sync"

	tmpl "github.com/toba/go-template-lsp/internal/template"
	"github.com/toba/go-template-lsp/internal/template/analyzer"
	"github.com/toba/go-template-lsp/internal/template/lexer"
	"github.com/toba/go-template-lsp/internal/template/parser"
	"github.com/toba/lsp/pathutil"
	"go.lsp.dev/protocol"
)

// handler implements server.Handler and optional handler interfaces
// for the Go template language server.
type handler struct {
	mu sync.Mutex

	rootPath  string
	initOpts  initializationOptions
	rawFiles  map[string][]byte
	parsed    map[string]*parser.GroupStatementNode
	analyzed  map[string]*analyzer.FileDefinition
	errsParse map[string][]lexer.Error
	errsAnal  map[string][]lexer.Error

	// filesOpenedByEditor tracks content of files currently open in the editor.
	filesOpenedByEditor map[string]string
}

// initializationOptions holds formatting settings from the client.
type initializationOptions struct {
	PrintWidth   int    `json:"printWidth"`
	AttrWrapMode string `json:"attrWrapMode"`
}

func newHandler() *handler {
	return &handler{
		rawFiles:            make(map[string][]byte),
		parsed:              make(map[string]*parser.GroupStatementNode),
		analyzed:            make(map[string]*analyzer.FileDefinition),
		errsParse:           make(map[string][]lexer.Error),
		errsAnal:            make(map[string][]lexer.Error),
		filesOpenedByEditor: make(map[string]string),
	}
}

// Initialize implements server.Handler.
func (h *handler) Initialize(
	_ context.Context,
	params *protocol.InitializeParams,
) (protocol.ServerCapabilities, error) {
	// Parse initialization options.
	if params.InitializationOptions != nil {
		raw, err := json.Marshal(params.InitializationOptions)
		if err == nil {
			_ = json.Unmarshal(raw, &h.initOpts)
		}
	}

	// Extract root URI from WorkspaceFolders or fallback to RootURI.
	var rootURI string
	if len(params.WorkspaceFolders) > 0 {
		rootURI = params.WorkspaceFolders[0].URI
	} else {
		rootURI = string(
			params.RootURI,
		) //nolint:staticcheck // fallback for older clients
	}
	if rootURI != "" {
		h.rootPath = pathutil.URIToFilePath(rootURI)
	}

	// Scan for custom template functions in Go source files.
	if h.rootPath != "" {
		customFuncs, err := analyzer.ScanWorkspaceForFuncMap(h.rootPath)
		if err != nil {
			slog.Warn("failed to scan for custom template functions",
				"error", err.Error())
		} else if len(customFuncs) > 0 {
			tmpl.SetWorkspaceCustomFunctions(customFuncs)
			funcNames := make([]string, 0, len(customFuncs))
			for name := range customFuncs {
				funcNames = append(funcNames, name)
			}
			slog.Info("discovered custom template functions",
				"functions", funcNames)
		}

		// Load workspace files.
		h.mu.Lock()
		h.rawFiles = pathutil.ConvertMapKeysToURI(
			tmpl.OpenProjectFiles(h.rootPath, TargetFileExtensions),
		)
		// Parse all workspace files.
		for uri, content := range h.rawFiles {
			tree, errs := tmpl.ParseSingleFile(content)
			h.parsed[uri] = tree
			h.errsParse[uri] = errs
		}
		// Run initial cross-file analysis.
		h.runAnalysis()
		h.mu.Unlock()
	}

	return protocol.ServerCapabilities{
		HoverProvider:              true,
		DefinitionProvider:         true,
		DocumentFormattingProvider: true,
		FoldingRangeProvider:       true,
		DocumentHighlightProvider:  true,
		DocumentLinkProvider: &protocol.DocumentLinkOptions{
			ResolveProvider: false,
		},
		SemanticTokensProvider: semanticTokensProviderOptions(),
	}, nil
}

// Diagnostics implements server.Handler.
// Called when a document is opened or changed.
func (h *handler) Diagnostics(
	_ context.Context,
	uri protocol.DocumentURI,
	content string,
) ([]protocol.Diagnostic, error) {
	uriStr := string(uri)

	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if file is inside workspace.
	if h.rootPath != "" &&
		!isFileInsideWorkspace(uriStr, h.rootPath, TargetFileExtensions) {
		return nil, nil
	}

	// Update file content.
	h.rawFiles[uriStr] = []byte(content)
	h.filesOpenedByEditor[uriStr] = content

	// Parse the changed file.
	tree, parseErrs := tmpl.ParseSingleFile([]byte(content))
	h.parsed[uriStr] = tree
	h.errsParse[uriStr] = parseErrs

	// Run cross-file analysis for the workspace.
	h.runAnalysis()

	// Return diagnostics for the requested URI.
	return h.diagnosticsForURI(uriStr), nil
}

// Shutdown implements server.Handler.
func (h *handler) Shutdown(_ context.Context) error {
	return nil
}

// Hover implements server.HoverHandler.
func (h *handler) Hover(
	_ context.Context,
	params *protocol.HoverParams,
) (*protocol.Hover, error) { //nolint:unparam // error return required by interface
	uriStr := string(params.TextDocument.URI)

	h.mu.Lock()
	file := h.analyzed[uriStr]
	h.mu.Unlock()

	if file == nil {
		return nil, nil
	}

	position := lexer.Position{
		Line:      int(params.Position.Line),
		Character: int(params.Position.Character),
	}

	typeStr, reach := tmpl.Hover(file, position)
	if typeStr == "" {
		return nil, nil
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.Markdown,
			Value: typeStr,
		},
		Range: convertRangePtr(reach),
	}, nil
}

// Definition implements server.DefinitionHandler.
func (h *handler) Definition(
	_ context.Context,
	params *protocol.DefinitionParams,
) ([]protocol.Location, error) {
	uriStr := string(params.TextDocument.URI)

	h.mu.Lock()
	currentFile := h.analyzed[uriStr]
	rawFiles := maps.Clone(h.rawFiles)
	h.mu.Unlock()

	if currentFile == nil {
		return nil, nil
	}

	position := lexer.Position{
		Line:      int(params.Position.Line),
		Character: int(params.Position.Character),
	}

	fileNames, reaches, defErr := tmpl.GoToDefinition(currentFile, position)
	if defErr != nil {
		return nil, fmt.Errorf("go to definition: %w", defErr)
	}

	_ = rawFiles // used by the original code but not needed for location results

	var locations []protocol.Location
	for i := range fileNames {
		if fileNames[i] == "" {
			continue
		}
		locations = append(locations, protocol.Location{
			URI:   protocol.DocumentURI(fileNames[i]),
			Range: convertRange(reaches[i]),
		})
	}

	return locations, nil
}

// Formatting implements server.FormattingHandler.
func (h *handler) Formatting(
	_ context.Context,
	params *protocol.DocumentFormattingParams,
) ([]protocol.TextEdit, error) {
	uriStr := string(params.TextDocument.URI)

	h.mu.Lock()
	content, ok := h.filesOpenedByEditor[uriStr]
	initOpts := h.initOpts
	h.mu.Unlock()

	var fileContent []byte
	if ok {
		fileContent = []byte(content)
	} else {
		// Fallback to reading from disk.
		parsed, parseErr := url.Parse(uriStr)
		if parseErr != nil {
			return nil, fmt.Errorf("parse URI: %w", parseErr)
		}
		if parsed.Path == "" {
			return nil, nil
		}
		var readErr error
		fileContent, readErr = os.ReadFile(
			parsed.Path,
		) //nolint:gosec // path from editor URI
		if readErr != nil {
			return nil, fmt.Errorf("read file: %w", readErr)
		}
	}

	formatted := tmpl.Format(fileContent, tmpl.FormatOptions{
		TabSize:      int(params.Options.TabSize),
		InsertSpaces: params.Options.InsertSpaces,
		PrintWidth:   initOpts.PrintWidth,
		AttrWrapMode: initOpts.AttrWrapMode,
	})

	// Count lines in original content.
	lineCount := uint32(0)
	for _, b := range fileContent {
		if b == '\n' {
			lineCount++
		}
	}

	return []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: lineCount + 1, Character: 0},
			},
			NewText: string(formatted),
		},
	}, nil
}

// runAnalysis performs cross-file semantic analysis on all parsed files.
// Caller must hold h.mu.
func (h *handler) runAnalysis() {
	if len(h.parsed) == 0 {
		return
	}

	chainedFiles := tmpl.DefinitionAnalysisWithinWorkspace(h.parsed)

	for _, fileAnalyzed := range chainedFiles {
		h.analyzed[fileAnalyzed.FileName] = fileAnalyzed.File
		h.errsAnal[fileAnalyzed.FileName] = fileAnalyzed.Errs
	}
}

// diagnosticsForURI collects parse + analysis errors for a URI.
// Caller must hold h.mu.
func (h *handler) diagnosticsForURI(uri string) []protocol.Diagnostic {
	diags := make(
		[]protocol.Diagnostic,
		0,
		len(h.errsParse[uri])+len(h.errsAnal[uri]),
	)

	for _, err := range h.errsParse[uri] {
		diags = append(diags, protocol.Diagnostic{
			Range:    convertRange(err.GetRange()),
			Severity: protocol.DiagnosticSeverityError,
			Message:  err.GetError(),
		})
	}

	for _, err := range h.errsAnal[uri] {
		diags = append(diags, protocol.Diagnostic{
			Range:    convertRange(err.GetRange()),
			Severity: protocol.DiagnosticSeverityError,
			Message:  err.GetError(),
		})
	}

	return diags
}

// safeUint32 converts an int to uint32, clamping negatives to 0.
func safeUint32(v int) uint32 {
	if v < 0 {
		return 0
	}
	return uint32(v) //nolint:gosec // bounds checked above
}

// convertRange converts a lexer.Range to a protocol.Range.
func convertRange(r lexer.Range) protocol.Range {
	if r.IsEmpty() {
		return protocol.Range{}
	}
	return protocol.Range{
		Start: protocol.Position{
			Line:      safeUint32(r.Start.Line),
			Character: safeUint32(r.Start.Character),
		},
		End: protocol.Position{
			Line:      safeUint32(r.End.Line),
			Character: safeUint32(r.End.Character),
		},
	}
}

// convertRangePtr converts a lexer.Range to a *protocol.Range.
func convertRangePtr(r lexer.Range) *protocol.Range {
	pr := convertRange(r)
	return &pr
}

var semanticTokenTypes = []string{
	"keyword",
	"variable",
	"function",
	"property",
	"string",
	"number",
	"operator",
	"comment",
}

var semanticTokenModifiers = []string{
	"declaration",
	"definition",
	"readonly",
}

// semanticTokensProviderOptions returns the options struct for the
// SemanticTokensProvider capability. The go.lsp.dev/protocol library
// uses interface{} for this field, so we build it as a map.
func semanticTokensProviderOptions() any {
	return map[string]any{
		"legend": map[string]any{
			"tokenTypes":     semanticTokenTypes,
			"tokenModifiers": semanticTokenModifiers,
		},
		"full":  true,
		"range": false,
	}
}

// isFileInsideWorkspace checks if a file is inside the workspace and has allowed extension.
func isFileInsideWorkspace(uri, rootPath string, allowedFileExtensions []string) bool {
	rootURI := pathutil.FilePathToURI(rootPath)
	if !strings.HasPrefix(uri, rootURI) {
		return false
	}
	return tmpl.HasFileExtension(uri, allowedFileExtensions)
}
