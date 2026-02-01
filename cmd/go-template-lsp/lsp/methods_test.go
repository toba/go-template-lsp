package lsp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessFormattingRequest_FallbackToDisk(t *testing.T) {
	// Create a temp file with unformatted content
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.gohtml")
	content := "<div>\n<p>hello</p>\n</div>\n"
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	fileUri := "file://" + filePath

	// Build a formatting request
	req := RequestMessage[DocumentFormattingParams]{
		JsonRpc: "2.0",
		Id:      1,
		Method:  "textDocument/formatting",
		Params: DocumentFormattingParams{
			TextDocument: TextDocumentIdentifier{Uri: fileUri},
			Options:      FormattingOptions{TabSize: 4, InsertSpaces: false},
		},
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	// Empty openFiles — file not tracked via didOpen
	openFiles := map[string]string{}

	response, _ := ProcessFormattingRequest(
		data,
		openFiles,
		InitializationOptions{},
	)

	if response == nil {
		t.Fatal(
			"expected non-nil response when file exists on disk but is not in textFromClient",
		)
	}

	var res ResponseMessage[[]TextEdit]
	if err := json.Unmarshal(response, &res); err != nil {
		t.Fatal(err)
	}

	if len(res.Result) == 0 {
		t.Fatal("expected at least one TextEdit in response")
	}

	if res.Result[0].NewText == "" {
		t.Fatal("expected non-empty NewText in edit")
	}
}

func TestProcessFormattingRequest_UsesOpenFilesNotDisk(t *testing.T) {
	// Simulate the race condition: disk has stale content, openFiles has latest.
	// The formatter should use openFiles content, not the stale disk content.
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.gohtml")
	staleContent := "<div>\n<p>stale</p>\n</div>\n"
	if err := os.WriteFile(filePath, []byte(staleContent), 0600); err != nil {
		t.Fatal(err)
	}

	fileUri := "file://" + filePath

	req := RequestMessage[DocumentFormattingParams]{
		JsonRpc: "2.0",
		Id:      1,
		Method:  "textDocument/formatting",
		Params: DocumentFormattingParams{
			TextDocument: TextDocumentIdentifier{Uri: fileUri},
			Options:      FormattingOptions{TabSize: 4, InsertSpaces: false},
		},
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	// openFiles has newer content that differs from disk
	latestContent := "<div>\n\t<p>latest</p>\n</div>\n"
	openFiles := map[string]string{
		fileUri: latestContent,
	}

	response, _ := ProcessFormattingRequest(
		data,
		openFiles,
		InitializationOptions{},
	)

	if response == nil {
		t.Fatal("expected non-nil response")
	}

	var res ResponseMessage[[]TextEdit]
	if err := json.Unmarshal(response, &res); err != nil {
		t.Fatal(err)
	}

	if len(res.Result) == 0 {
		t.Fatal("expected at least one TextEdit in response")
	}

	// The formatted output should be based on latestContent, not staleContent.
	// "latest" should appear in the output, "stale" should not.
	formatted := res.Result[0].NewText

	if !strings.Contains(formatted, "latest") {
		t.Fatalf(
			"expected formatted output to contain 'latest' (from openFiles), got: %s",
			formatted,
		)
	}
	if strings.Contains(formatted, "stale") {
		t.Fatalf(
			"formatted output contains 'stale' (from disk), should use openFiles content, got: %s",
			formatted,
		)
	}
}
