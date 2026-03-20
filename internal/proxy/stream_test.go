package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestStreamSSEResponse(t *testing.T) {
	// Fake upstream that sends SSE events with delays
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("upstream ResponseWriter does not support flushing")
		}

		events := []string{
			"event: message_start\ndata: {\"type\":\"message_start\"}\n\n",
			"event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"delta\":{\"text\":\"Hello\"}}\n\n",
			"event: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"delta\":{\"text\":\" world\"}}\n\n",
			"event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n",
		}

		for _, event := range events {
			fmt.Fprint(w, event)
			flusher.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}))
	defer upstream.Close()

	// Point provider at fake upstream
	originalProviders := providers
	providers = map[string]Provider{
		"anthropic": {Name: "anthropic", BaseURL: upstream.URL, CompressRequests: false},
	}
	defer func() { providers = originalProviders }()

	srv := &Server{
		ListenAddr: ":0",
		Algorithm:  "gzip",
		client:     &http.Client{},
		stats:      &Stats{},
	}

	req := httptest.NewRequest("POST", "/v1/messages", strings.NewReader(`{"stream":true}`))
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	srv.handleProxy(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/event-stream") {
		t.Errorf("Content-Type = %q, want text/event-stream", ct)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	// Verify all events came through
	expectedParts := []string{
		"message_start",
		"Hello",
		" world",
		"message_stop",
	}
	for _, part := range expectedParts {
		if !strings.Contains(bodyStr, part) {
			t.Errorf("reponse body missing %q\ngot: %s", part, bodyStr)
		}
	}
}

func TestIsSSE(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        bool
	}{
		{"sse stream", "text/event-stream", true},
		{"sse with charset", "text/event-stream; charset=utf-8", true},
		{"json response", "application/json", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{Header: http.Header{}}
			resp.Header.Set("Content-Type", tt.contentType)
			if got := isSSE(resp); got != tt.want {
				t.Errorf("isSSE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNonSSEResponseNotStreamed(t *testing.T) {
	// False upstream returning normal JSON
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"type":"message","content":[{"text":"Hello"}]}`))
	}))
	defer upstream.Close()

	originalProviders := providers
	providers = map[string]Provider{
		"anthropic": {Name: "anthropic", BaseURL: upstream.URL, CompressRequests: false},
	}
	defer func() { providers = originalProviders }()

	srv := &Server{
		ListenAddr: ":0",
		Algorithm:  "gzip",
		client:     &http.Client{},
		stats:      &Stats{},
	}

	req := httptest.NewRequest("POST", "/v1/messages", strings.NewReader(`{"stream":false}`))
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	srv.handleProxy(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Hello") {
		t.Errorf("non-SSE response body missing expected content: %s", string(body))
	}

}
