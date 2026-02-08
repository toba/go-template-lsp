package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"

	tmpl "github.com/toba/go-template-lsp/internal/template"
	"github.com/toba/go-template-lsp/internal/template/analyzer"
	"github.com/toba/go-template-lsp/internal/template/parser"
)

// checkDiagnostic represents a single diagnostic finding (1-indexed positions).
type checkDiagnostic struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	EndLine int    `json:"endLine"`
	EndCol  int    `json:"endColumn"`
	Message string `json:"message"`
}

// runCheck runs the diagnostic pipeline on the given directory and outputs results.
// Returns exit code: 0 = clean, 1 = errors found, 2 = tool failure.
func runCheck(args []string, jsonOutput bool, extraFuncs []string) (exitCode int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "go-template-lsp check: internal error: %v\n", r)
			exitCode = 2
		}
	}()

	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "go-template-lsp check: %v\n", err)
		return 2
	}

	info, err := os.Stat(absDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "go-template-lsp check: %v\n", err)
		return 2
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "go-template-lsp check: not a directory: %s\n", absDir)
		return 2
	}

	// Discover custom template functions from Go source files.
	// Scan from the Go module root so FuncMaps defined outside the template
	// directory are discovered.
	customFuncs, _ := analyzer.ScanWorkspaceForFuncMap(analyzer.FindModuleRoot(absDir))
	if len(extraFuncs) > 0 {
		if customFuncs == nil {
			customFuncs = make(map[string]*analyzer.FunctionDefinition)
		}
		maps.Copy(customFuncs, analyzer.FunctionsFromNames(extraFuncs))
	}
	if len(customFuncs) > 0 {
		tmpl.SetWorkspaceCustomFunctions(customFuncs)
	}

	// Open all template files in the directory.
	rawFiles := tmpl.OpenProjectFiles(absDir, TargetFileExtensions)
	if len(rawFiles) == 0 {
		if jsonOutput {
			fmt.Println("[]")
		}
		return 0
	}

	// Parse each file individually to keep errors associated with filenames.
	parsedFiles := make(map[string]*parser.GroupStatementNode, len(rawFiles))
	parseErrors := make(map[string][]tmpl.Error, len(rawFiles))
	for filePath, content := range rawFiles {
		tree, errs := tmpl.ParseSingleFile(content)
		parsedFiles[filePath] = tree
		if len(errs) > 0 {
			parseErrors[filePath] = errs
		}
	}

	// Run cross-file semantic analysis.
	analysisResults := tmpl.DefinitionAnalysisWithinWorkspace(parsedFiles)

	diagnostics := collectDiagnostics(absDir, parseErrors, analysisResults)

	if len(diagnostics) == 0 {
		if jsonOutput {
			fmt.Println("[]")
		}
		return 0
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(diagnostics); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"go-template-lsp check: failed to encode JSON: %v\n",
				err,
			)
			return 2
		}
	} else {
		for _, d := range diagnostics {
			fmt.Printf("%s:%d:%d: %s\n", d.File, d.Line, d.Column, d.Message)
		}
	}

	return 1
}

// collectDiagnostics merges parse errors and analysis errors into a sorted slice
// of checkDiagnostic with 1-indexed positions and file paths relative to rootDir.
func collectDiagnostics(
	rootDir string,
	parseErrors map[string][]tmpl.Error,
	analysisResults []tmpl.FileAnalysisAndError,
) []checkDiagnostic {
	diagnostics := make([]checkDiagnostic, 0, len(parseErrors)+len(analysisResults))

	// Collect parse errors.
	for filePath, errs := range parseErrors {
		rel := relativePath(rootDir, filePath)
		for _, e := range errs {
			r := e.GetRange()
			diagnostics = append(diagnostics, checkDiagnostic{
				File:    rel,
				Line:    r.Start.Line + 1,
				Column:  r.Start.Character + 1,
				EndLine: r.End.Line + 1,
				EndCol:  r.End.Character + 1,
				Message: e.GetError(),
			})
		}
	}

	// Collect analysis errors.
	for _, result := range analysisResults {
		rel := relativePath(rootDir, result.FileName)
		for _, e := range result.Errs {
			r := e.GetRange()
			diagnostics = append(diagnostics, checkDiagnostic{
				File:    rel,
				Line:    r.Start.Line + 1,
				Column:  r.Start.Character + 1,
				EndLine: r.End.Line + 1,
				EndCol:  r.End.Character + 1,
				Message: e.GetError(),
			})
		}
	}

	// Sort by file, then line, then column.
	slices.SortFunc(diagnostics, func(a, b checkDiagnostic) int {
		if c := cmp.Compare(a.File, b.File); c != 0 {
			return c
		}
		if c := cmp.Compare(a.Line, b.Line); c != 0 {
			return c
		}
		return cmp.Compare(a.Column, b.Column)
	})

	return diagnostics
}

// relativePath returns filePath relative to rootDir, falling back to the absolute path.
func relativePath(rootDir, filePath string) string {
	rel, err := filepath.Rel(rootDir, filePath)
	if err != nil {
		return filePath
	}
	return rel
}
