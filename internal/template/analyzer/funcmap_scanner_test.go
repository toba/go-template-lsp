package analyzer_test

import (
	"path/filepath"
	"testing"

	"github.com/toba/go-template-lsp/internal/template/analyzer"
	"github.com/toba/go-template-lsp/internal/template/testutil"
)

func TestScanWorkspaceForFuncMap(t *testing.T) {
	// Test file with template.FuncMap composite literal
	testFile1 := `package main

import (
	"strings"
	"text/template"
)

var funcs = template.FuncMap{
	"lower": strings.ToLower,
	"upper": strings.ToUpper,
	"custom": func(s string) string { return s },
}
`
	// Test file with html/template import
	testFile2 := `package main

import (
	"html/template"
)

var htmlFuncs = template.FuncMap{
	"safe": func(s string) template.HTML { return template.HTML(s) },
	"dict": func(values ...any) map[string]any { return nil },
}
`
	tmpDir := testutil.TempDir(t, map[string]string{
		"funcs1.go": testFile1,
		"funcs2.go": testFile2,
	})

	// Scan the workspace
	funcs, err := analyzer.ScanWorkspaceForFuncMap(tmpDir)
	if err != nil {
		t.Fatalf("ScanWorkspaceForFuncMap failed: %v", err)
	}

	// Verify expected functions were found
	expectedFuncs := []string{"lower", "upper", "custom", "safe", "dict"}
	for _, name := range expectedFuncs {
		if _, ok := funcs[name]; !ok {
			t.Errorf("expected function %q not found", name)
		}
	}

	// Verify the count
	if len(funcs) != len(expectedFuncs) {
		t.Errorf("expected %d functions, got %d", len(expectedFuncs), len(funcs))
		for name := range funcs {
			t.Logf("found function: %s", name)
		}
	}
}

func TestScanWorkspaceForFuncMap_NoTemplateImport(t *testing.T) {
	// Test file without template import
	testFile := `package main

import "fmt"

func main() {
	fmt.Println("hello")
}
`
	tmpDir := testutil.TempDir(t, map[string]string{
		"main.go": testFile,
	})

	funcs, err := analyzer.ScanWorkspaceForFuncMap(tmpDir)
	if err != nil {
		t.Fatalf("ScanWorkspaceForFuncMap failed: %v", err)
	}

	if len(funcs) != 0 {
		t.Errorf("expected no functions, got %d", len(funcs))
	}
}

func TestScanWorkspaceForFuncMap_SkipsVendor(t *testing.T) {
	// Test file in vendor (should be skipped)
	testFile := `package vendor

import "text/template"

var funcs = template.FuncMap{
	"vendorFunc": func() {},
}
`
	tmpDir := testutil.TempDir(t, map[string]string{
		"vendor/vendor.go": testFile,
	})

	funcs, err := analyzer.ScanWorkspaceForFuncMap(tmpDir)
	if err != nil {
		t.Fatalf("ScanWorkspaceForFuncMap failed: %v", err)
	}

	if _, ok := funcs["vendorFunc"]; ok {
		t.Error("vendor function should have been skipped")
	}
}

func TestScanWorkspaceForFuncMap_SkipsTestFiles(t *testing.T) {
	// Test file that is a test (should be skipped)
	testFile := `package main

import "text/template"

var funcs = template.FuncMap{
	"testFunc": func() {},
}
`
	tmpDir := testutil.TempDir(t, map[string]string{
		"funcs_test.go": testFile,
	})

	funcs, err := analyzer.ScanWorkspaceForFuncMap(tmpDir)
	if err != nil {
		t.Fatalf("ScanWorkspaceForFuncMap failed: %v", err)
	}

	if _, ok := funcs["testFunc"]; ok {
		t.Error("test file function should have been skipped")
	}
}

func TestScanWorkspaceForFuncMap_AliasedImport(t *testing.T) {
	// Test file with aliased import
	testFile := `package main

import (
	tmpl "text/template"
)

var funcs = tmpl.FuncMap{
	"aliased": func() {},
}
`
	tmpDir := testutil.TempDir(t, map[string]string{
		"aliased.go": testFile,
	})

	funcs, err := analyzer.ScanWorkspaceForFuncMap(tmpDir)
	if err != nil {
		t.Fatalf("ScanWorkspaceForFuncMap failed: %v", err)
	}

	if _, ok := funcs["aliased"]; !ok {
		t.Error("aliased import function should have been found")
	}
}

func TestFindModuleRoot(t *testing.T) {
	t.Run("finds go.mod in parent", func(t *testing.T) {
		// Create:  root/go.mod  and  root/sub/dir/
		tmpDir := testutil.TempDir(t, map[string]string{
			"go.mod":           "module example.com/test\n",
			"sub/dir/dummy.go": "package dummy\n",
		})
		subDir := filepath.Join(tmpDir, "sub", "dir")
		got := analyzer.FindModuleRoot(subDir)
		if got != tmpDir {
			t.Errorf("FindModuleRoot(%q) = %q, want %q", subDir, got, tmpDir)
		}
	})

	t.Run("returns dir when no go.mod exists", func(t *testing.T) {
		tmpDir := testutil.TempDir(t, map[string]string{
			"dummy.txt": "",
		})
		got := analyzer.FindModuleRoot(tmpDir)
		// Without go.mod anywhere, it should walk all the way to / and return /
		// We just verify it doesn't panic and returns a valid directory.
		if got == "" {
			t.Error("FindModuleRoot returned empty string")
		}
	})

	t.Run("finds go.mod in same directory", func(t *testing.T) {
		tmpDir := testutil.TempDir(t, map[string]string{
			"go.mod": "module example.com/test\n",
		})
		got := analyzer.FindModuleRoot(tmpDir)
		if got != tmpDir {
			t.Errorf("FindModuleRoot(%q) = %q, want %q", tmpDir, got, tmpDir)
		}
	})
}

func TestFunctionsFromNames(t *testing.T) {
	names := []string{"asset", "renderPage", "", "dict"}
	funcs := analyzer.FunctionsFromNames(names)

	if len(funcs) != 3 {
		t.Fatalf("expected 3 functions, got %d", len(funcs))
	}

	for _, name := range []string{"asset", "renderPage", "dict"} {
		fd, ok := funcs[name]
		if !ok {
			t.Errorf("expected function %q not found", name)
			continue
		}
		if fd.Name() != name {
			t.Errorf("expected Name() = %q, got %q", name, fd.Name())
		}
		if fd.FileName() != "" {
			t.Errorf("expected empty FileName(), got %q", fd.FileName())
		}
	}

	if _, ok := funcs[""]; ok {
		t.Error("empty name should have been skipped")
	}
}

func TestScanWorkspaceForFuncMap_CapturesSourcePositions(t *testing.T) {
	// Line numbers (1-indexed in Go source, 0-indexed in lexer.Range):
	// Line 1: package main
	// Line 2: (blank)
	// Line 3: import (
	// Line 4:     "text/template"
	// Line 5: )
	// Line 6: (blank)
	// Line 7: var funcs = template.FuncMap{
	// Line 8:     "alpha": func() {},
	// Line 9:     "beta":  func() {},
	// Line 10: }
	testFile := `package main

import (
	"text/template"
)

var funcs = template.FuncMap{
	"alpha": func() {},
	"beta":  func() {},
}
`
	tmpDir := testutil.TempDir(t, map[string]string{
		"funcs.go": testFile,
	})

	funcs, err := analyzer.ScanWorkspaceForFuncMap(tmpDir)
	if err != nil {
		t.Fatalf("ScanWorkspaceForFuncMap failed: %v", err)
	}

	expectedFilePath := filepath.Join(tmpDir, "funcs.go")

	tests := []struct {
		name      string
		startLine int
		startChar int
		endChar   int
	}{
		// "alpha" is on line 8 (0-indexed: 7), tab-indented, starts at column 2 (after opening quote)
		{"alpha", 7, 2, 7},
		// "beta" is on line 9 (0-indexed: 8)
		{"beta", 8, 2, 6},
	}

	for _, tc := range tests {
		fn, ok := funcs[tc.name]
		if !ok {
			t.Errorf("expected function %q not found", tc.name)
			continue
		}

		rng := fn.Range()
		if rng.IsEmpty() {
			t.Errorf("function %q has empty range", tc.name)
			continue
		}

		if fn.FileName() != expectedFilePath {
			t.Errorf(
				"function %q: expected file %q, got %q",
				tc.name,
				expectedFilePath,
				fn.FileName(),
			)
		}

		if rng.Start.Line != tc.startLine {
			t.Errorf(
				"function %q: expected start line %d, got %d",
				tc.name,
				tc.startLine,
				rng.Start.Line,
			)
		}

		if rng.Start.Character != tc.startChar {
			t.Errorf(
				"function %q: expected start char %d, got %d",
				tc.name,
				tc.startChar,
				rng.Start.Character,
			)
		}

		if rng.End.Character != tc.endChar {
			t.Errorf(
				"function %q: expected end char %d, got %d",
				tc.name,
				tc.endChar,
				rng.End.Character,
			)
		}
	}
}
