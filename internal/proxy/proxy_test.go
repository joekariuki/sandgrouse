package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleProxy(t *testing.T) {
	// Create a fake upstream that echoes back requests details
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Path", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		body, _ := io.ReadAll(r.Body)
		//  Decompress if the request body is gzip-compressed
		if r.Header.Get("Content-Encoding") == "gzip" {
			body, _ = decompressGzip(body)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer upstream.Close()

	// Override providers map to point at our fake upstream
	originalProviders := providers
	providers = map[string]Provider{
		"anthropic": {Name: "anthropic", BaseURL: upstream.URL},
	}
	defer func() { providers = originalProviders }()

	srv := &Server{ListenAddr: ":0", client: &http.Client{}, stats: &Stats{}}

	t.Run("forwards request to upstream and returns response", func(t *testing.T) {
		reqBody := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}]}`
		req := httptest.NewRequest("POST", "/v1/messages", strings.NewReader(reqBody))
		req.Header.Set("anthropic-version", "2023-06-01")
		req.Header.Set("x-api-key", "sk-ant-test")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		srv.handleProxy(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != reqBody {
			t.Errorf("body = %q, want %q", string(body), reqBody)
		}

		if got := resp.Header.Get("X-Test-Path"); got != "/v1/messages" {
			t.Errorf("upstream path = %q, want %q", got, "/v1/messages")
		}
	})

	t.Run("returns 400 for unknown provider", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/v1/messages", nil)
		req.Header.Set("X-Custom", "value")

		w := httptest.NewRecorder()
		srv.handleProxy(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)

		}
	})

	t.Run("preserves query parameters", func(t *testing.T) {
		upstream2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.URL.RawQuery))
		}))
		defer upstream2.Close()

		providers["anthropic"] = Provider{Name: "anthropic", BaseURL: upstream2.URL}

		req := httptest.NewRequest("GET", "/v1/models?limit=10&offset=0", nil)
		req.Header.Set("anthropic-version", "2023-06-01")

		w := httptest.NewRecorder()
		srv.handleProxy(w, req)

		body, _ := io.ReadAll(w.Result().Body)
		if string(body) != "limit=10&offset=0" {
			t.Errorf("query = %q, want %q", string(body), "limit=10&offset=0")
		}
	})
}

func TestHandleProxyWithCompressedResponse(t *testing.T) {
	// Create a fake upstream that returns gzip-compressed JSON
	originalBody := `{"id":"msg_123","type":"message","role":"assistant","content":[{"type":"text","text":"Hello! I'd be happy to help you with your reverse proxy implementation. The sandgrouse proxy is looking great so far with provider detection, compression, and streaming support all working correctly."}],"model":"claude-sonnet-4-20250514","stop_reason":"end_turn","usage":{"input_tokens":42,"output_tokens":128}}`
	compressedBody, err := compressGzip([]byte(originalBody))
	if err != nil {
		t.Fatalf("compressGzip error: %v", err)
	}

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the proxy requested compressed responses
		ae := r.Header.Get("Accept-Encoding")
		if ae == "" {
			t.Error("proxy did not send Accept-Encoding header")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(compressedBody)
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

	reqBody := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}]}`
	req := httptest.NewRequest("POST", "/v1/messages", strings.NewReader(reqBody))
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	srv.handleProxy(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Response should be decompressed (no Content-Encoding header to client)
	if ce := resp.Header.Get("Content-Encoding"); ce != "" {
		t.Errorf("Content-Encoding should be stripped, got %q", ce)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != originalBody {
		t.Errorf("body = %q, want %q", string(body), originalBody)
	}

	// Stats should show response savings
	wireBytes := srv.stats.responseWireBytes.Load()
	origBytes := srv.stats.responseOriginalBytes.Load()
	if wireBytes >= origBytes {
		t.Errorf("wire bytes (%d) should be less than original bytes (%d)", wireBytes, origBytes)
	}
}
