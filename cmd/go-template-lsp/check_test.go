package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCheck_ValidDir(t *testing.T) {
	dir := filepath.Join("testdata", "check", "valid")
	code := runCheck([]string{dir}, false)
	if code != 0 {
		t.Errorf("expected exit code 0 for valid templates, got %d", code)
	}
}

func TestRunCheck_ErrorsDir(t *testing.T) {
	dir := filepath.Join("testdata", "check", "errors")
	code := runCheck([]string{dir}, false)
	if code != 1 {
		t.Errorf("expected exit code 1 for templates with errors, got %d", code)
	}
}

func TestRunCheck_ErrorsTextOutput(t *testing.T) {
	dir := filepath.Join("testdata", "check", "errors")

	// Capture stdout.
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	code := runCheck([]string{dir}, false)

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = old

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	output := string(out)

	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}

	// Verify the output contains the expected file:line:col: message format.
	if !strings.Contains(output, "broken.html:") {
		t.Errorf("expected output to mention broken.html, got: %s", output)
	}
	// Each line should follow the pattern file:line:col: message.
	for line := range strings.SplitSeq(strings.TrimSpace(output), "\n") {
		parts := strings.SplitN(line, ":", 4)
		if len(parts) < 4 {
			t.Errorf("unexpected output format: %s", line)
		}
	}
}

func TestRunCheck_ErrorsJSONOutput(t *testing.T) {
	dir := filepath.Join("testdata", "check", "errors")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	code := runCheck([]string{dir}, true)

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = old

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	output := string(out)

	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}

	var diagnostics []checkDiagnostic
	if err := json.Unmarshal([]byte(output), &diagnostics); err != nil {
		t.Fatalf("JSON output is not valid: %v\noutput: %s", err, output)
	}
	if len(diagnostics) == 0 {
		t.Fatal("expected at least one diagnostic in JSON output")
	}
	for _, d := range diagnostics {
		if d.File == "" || d.Message == "" || d.Line < 1 || d.Column < 1 {
			t.Errorf("diagnostic has missing fields: %+v", d)
		}
	}
}

func TestRunCheck_NonexistentDir(t *testing.T) {
	// Capture stderr to avoid noise.
	oldErr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	code := runCheck([]string{"/nonexistent/path/that/does/not/exist"}, false)

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stderr = oldErr

	if code != 2 {
		t.Errorf("expected exit code 2 for nonexistent directory, got %d", code)
	}
}

func TestRunCheck_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	code := runCheck([]string{dir}, false)
	if code != 0 {
		t.Errorf("expected exit code 0 for empty directory, got %d", code)
	}
}

func TestRunCheck_EmptyDirJSON(t *testing.T) {
	dir := t.TempDir()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	code := runCheck([]string{dir}, true)

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = old

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	output := strings.TrimSpace(string(out))

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if output != "[]" {
		t.Errorf("expected empty JSON array, got: %s", output)
	}
}

func TestRunCheck_DefaultDir(t *testing.T) {
	// With no args, runCheck should use the current directory and not crash.
	code := runCheck(nil, false)
	// The exit code depends on whether the CWD has templates with errors,
	// but it should not be 2 (tool failure).
	if code == 2 {
		t.Errorf("expected exit code 0 or 1, got 2 (tool failure)")
	}
}
