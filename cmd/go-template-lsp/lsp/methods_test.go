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

func TestProcessInitializeRequest_CapabilitiesAreObjects(t *testing.T) {
	// LSP spec allows ServerCapabilities provider fields to be either a boolean
	// or an options object, but strict clients (e.g. Zed) only accept the object
	// form. Verify that provider fields that have an Options type serialize as
	// JSON objects, not bare booleans.
	req := RequestMessage[InitializeParams]{
		JsonRpc: JSONRPCVersion,
		Id:      1,
		Method:  MethodInitialize,
		Params: InitializeParams{
			RootUri: "file:///workspace",
		},
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	response, _, _ := ProcessInitializeRequest(data, "test-lsp", "0.0.1")

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(response, &raw); err != nil {
		t.Fatal(err)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(raw["result"], &result); err != nil {
		t.Fatal(err)
	}

	var caps map[string]json.RawMessage
	if err := json.Unmarshal(result["capabilities"], &caps); err != nil {
		t.Fatal(err)
	}

	// These capability fields must serialize as objects per the LSP spec's
	// *Options types. A bare `true`/`false` will cause deserialization
	// failures in strict clients.
	objectCapabilities := []string{
		"documentLinkProvider",
		"semanticTokensProvider",
	}

	for _, key := range objectCapabilities {
		raw, ok := caps[key]
		if !ok {
			continue // omitted is fine (omitempty)
		}

		trimmed := strings.TrimSpace(string(raw))
		if trimmed == "true" || trimmed == "false" {
			t.Errorf(
				"capability %q serialized as boolean %s, expected an object",
				key,
				trimmed,
			)
		}
		if trimmed[0] != '{' {
			t.Errorf("capability %q should be a JSON object, got: %s", key, trimmed)
		}
	}
}

func TestBuildRegisterFileWatcherRequest(t *testing.T) {
	data := BuildRegisterFileWatcherRequest(42)
	if data == nil {
		t.Fatal("expected non-nil result")
	}

	var req RequestMessage[RegistrationParams]
	if err := json.Unmarshal(data, &req); err != nil {
		t.Fatal(err)
	}

	if req.Method != MethodRegisterCapability {
		t.Errorf("expected method %q, got %q", MethodRegisterCapability, req.Method)
	}
	if int(req.Id) != 42 {
		t.Errorf("expected id 42, got %d", req.Id)
	}
	if len(req.Params.Registrations) != 1 {
		t.Fatalf("expected 1 registration, got %d", len(req.Params.Registrations))
	}

	reg := req.Params.Registrations[0]
	if reg.Id != "go-file-watcher" {
		t.Errorf("expected registration id %q, got %q", "go-file-watcher", reg.Id)
	}
	if reg.Method != MethodDidChangeWatchedFiles {
		t.Errorf("expected method %q, got %q", MethodDidChangeWatchedFiles, reg.Method)
	}

	// Verify RegisterOptions has the correct structure
	raw, err := json.Marshal(reg.RegisterOptions)
	if err != nil {
		t.Fatal(err)
	}
	var opts DidChangeWatchedFilesRegistrationOptions
	if err := json.Unmarshal(raw, &opts); err != nil {
		t.Fatal(err)
	}
	if len(opts.Watchers) != 1 {
		t.Fatalf("expected 1 watcher, got %d", len(opts.Watchers))
	}
	if opts.Watchers[0].GlobPattern != "**/*.go" {
		t.Errorf("expected glob %q, got %q", "**/*.go", opts.Watchers[0].GlobPattern)
	}
	if opts.Watchers[0].Kind != WatchKindCreate|WatchKindChange|WatchKindDelete {
		t.Errorf("expected kind %d, got %d",
			WatchKindCreate|WatchKindChange|WatchKindDelete, opts.Watchers[0].Kind)
	}
}

func TestProcessDidChangeWatchedFilesNotification(t *testing.T) {
	tests := []struct {
		name     string
		changes  []FileEvent
		expected bool
	}{
		{
			name: "regular Go file changed",
			changes: []FileEvent{
				{Uri: "file:///workspace/main.go", Type: FileChangeChanged},
			},
			expected: true,
		},
		{
			name: "test file changed",
			changes: []FileEvent{
				{Uri: "file:///workspace/main_test.go", Type: FileChangeChanged},
			},
			expected: false,
		},
		{
			name: "non-Go file changed",
			changes: []FileEvent{
				{Uri: "file:///workspace/template.html", Type: FileChangeChanged},
			},
			expected: false,
		},
		{
			name: "hidden Go file",
			changes: []FileEvent{
				{Uri: "file:///workspace/.hidden.go", Type: FileChangeChanged},
			},
			expected: false,
		},
		{
			name: "Go file created",
			changes: []FileEvent{
				{Uri: "file:///workspace/new.go", Type: FileChangeCreated},
			},
			expected: true,
		},
		{
			name: "Go file deleted",
			changes: []FileEvent{
				{Uri: "file:///workspace/old.go", Type: FileChangeDeleted},
			},
			expected: true,
		},
		{
			name: "mixed: only test files",
			changes: []FileEvent{
				{Uri: "file:///workspace/foo_test.go", Type: FileChangeChanged},
				{Uri: "file:///workspace/bar_test.go", Type: FileChangeChanged},
			},
			expected: false,
		},
		{
			name: "mixed: one regular Go file among test files",
			changes: []FileEvent{
				{Uri: "file:///workspace/foo_test.go", Type: FileChangeChanged},
				{Uri: "file:///workspace/bar.go", Type: FileChangeChanged},
			},
			expected: true,
		},
		{
			name:     "empty changes",
			changes:  []FileEvent{},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			notification := NotificationMessage[DidChangeWatchedFilesParams]{
				JsonRpc: JSONRPCVersion,
				Method:  MethodDidChangeWatchedFiles,
				Params:  DidChangeWatchedFilesParams{Changes: tc.changes},
			}
			data, err := json.Marshal(notification)
			if err != nil {
				t.Fatal(err)
			}

			result := ProcessDidChangeWatchedFilesNotification(data)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestProcessDidChangeWatchedFilesNotification_InvalidJSON(t *testing.T) {
	result := ProcessDidChangeWatchedFilesNotification([]byte(`{invalid json`))
	if result {
		t.Error("expected false for invalid JSON")
	}
}
