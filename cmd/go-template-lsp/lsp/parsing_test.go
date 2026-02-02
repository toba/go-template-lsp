package lsp

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestSendToLspClient_ConcurrentRace(t *testing.T) {
	const goroutines = 50
	const msgsPerGoroutine = 20

	var buf bytes.Buffer
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(goroutines)

	for i := range goroutines {
		go func() {
			defer wg.Done()
			for j := range msgsPerGoroutine {
				msg := []byte(
					`{"jsonrpc":"2.0","id":` + strconv.Itoa(
						i*msgsPerGoroutine+j,
					) + `,"result":null}`,
				)
				mu.Lock()
				SendToLspClient(&buf, msg)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Verify all messages can be decoded
	data := buf.Bytes()
	decoded := 0

	for len(data) > 0 {
		idx := bytes.Index(data, []byte(HeaderDelimiter))
		if idx == -1 {
			break
		}

		contentLength, err := getHeaderContentLength(data[:idx])
		if err != nil {
			t.Fatalf("failed to parse header at message %d: %v", decoded, err)
		}

		bodyStart := idx + len(HeaderDelimiter)
		if bodyStart+contentLength > len(data) {
			t.Fatalf(
				"truncated message %d: need %d bytes, have %d",
				decoded,
				contentLength,
				len(data)-bodyStart,
			)
		}

		body := data[bodyStart : bodyStart+contentLength]

		var msg map[string]any
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Fatalf("invalid JSON in message %d: %v\nbody: %q", decoded, err, body)
		}

		data = data[bodyStart+contentLength:]
		decoded++
	}

	expected := goroutines * msgsPerGoroutine
	if decoded != expected {
		t.Errorf("decoded %d messages, want %d", decoded, expected)
	}
}

func TestSendToLspClient_ConcurrentFrameIntegrity(t *testing.T) {
	const goroutines = 30
	const msgsPerGoroutine = 15

	pr, pw := io.Pipe()
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(goroutines)

	for i := range goroutines {
		go func() {
			defer wg.Done()
			for j := range msgsPerGoroutine {
				msg := []byte(
					`{"jsonrpc":"2.0","id":` + strconv.Itoa(
						i*msgsPerGoroutine+j,
					) + `,"result":"ok"}`,
				)
				mu.Lock()
				SendToLspClient(pw, msg)
				mu.Unlock()
			}
		}()
	}

	go func() {
		wg.Wait()
		_ = pw.Close()
	}()

	total := goroutines * msgsPerGoroutine
	decoded := 0
	scanner := ReceiveInput(pr)

	for scanner.Scan() {
		body := scanner.Bytes()
		if len(body) == 0 {
			continue
		}

		if !strings.HasPrefix(string(body), "{") {
			t.Fatalf("message %d doesn't start with '{': %q", decoded, body)
		}

		var msg map[string]any
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Fatalf("invalid JSON in message %d: %v\nbody: %q", decoded, err, body)
		}

		if msg["jsonrpc"] != "2.0" {
			t.Fatalf("message %d has wrong jsonrpc: %v", decoded, msg["jsonrpc"])
		}

		decoded++
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner error: %v", err)
	}

	if decoded != total {
		t.Errorf("decoded %d messages, want %d", decoded, total)
	}
}
