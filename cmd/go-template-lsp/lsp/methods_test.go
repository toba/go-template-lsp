package lsp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
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

	// Empty textFromClient — file not tracked via didOpen
	textFromClient := map[string][]byte{}
	mu := &sync.Mutex{}

	response, _ := ProcessFormattingRequest(
		data,
		textFromClient,
		mu,
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
