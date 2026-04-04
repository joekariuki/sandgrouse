package proxy

import (
	"net/http"
	"testing"
)

func TestDetectProvider(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		wantName string
		wantOK   bool
	}{
		{
			name:     "anthropic detected from anthropic-version header",
			headers:  map[string]string{"anthropic-version": "2023-06-01"},
			wantName: "anthropic",
			wantOK:   true,
		},
		{
			name:     "anthropic detected with both headers",
			headers:  map[string]string{"anthropic-version": "2023-06-01", "x-api-key": "sk-ant-test"},
			wantName: "anthropic",
			wantOK:   true,
		},
		{
			name:     "openai detected from bearer token",
			headers:  map[string]string{"Authorization": "Bearer sk-test123"},
			wantName: "openai",
			wantOK:   true,
		},
		{
			name:    "unknown provider returns false",
			headers: map[string]string{"X-Custom": "value"},
			wantOK:  false,
		},
		{
			name:    "empty headers returns false",
			headers: map[string]string{},
			wantOK:  false,
		},
		{
			name:    "authorization without bearer prefix returns false",
			headers: map[string]string{"Authorization": "Basic dXNlxjpwYXNz"},
			wantOK:  false,
		},
		{
			name:     "gemini detected from x-goog-api-key header",
			headers:  map[string]string{"x-goog-api-key": "AIzaSy-test123"},
			wantName: "gemini",
			wantOK:   true,
		},
		{
			name:     "gemini takes priority over bearer when both present",
			headers:  map[string]string{"x-goog-api-key": "AIzaSy-test123", "Authorization": "Bearer ya29.oauth-token"},
			wantName: "gemini",
			wantOK:   true,
		},
		{
			name:     "anthropic takes priority over gemini when both present",
			headers:  map[string]string{"anthropic-version": "2023-06-01", "x-goog-api-key": "AIzaSy-test123"},
			wantName: "anthropic",
			wantOK:   true,
		},
		{
			name:     "all three provider headers present, anthropic wins",
			headers:  map[string]string{"anthropic-version": "2023-06-01", "x-goog-api-key": "AIzaSy-test123", "Authorization": "Bearer sk-test"},
			wantName: "anthropic",
			wantOK:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/v1/messages", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			p, ok := detectProvider(req)
			if ok != tt.wantOK {
				t.Fatalf("detectProvider() ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && p.Name != tt.wantName {
				t.Errorf("detectProvider() name = %q, want %q", p.Name, tt.wantName)
			}
		})
	}
}
