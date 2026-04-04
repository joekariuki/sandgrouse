package proxy

import (
	"net/http"
	"strings"
)

// Provider represents an upstream LLM API provider.
type Provider struct {
	Name             string
	BaseURL          string
	CompressRequests bool
}

var providers = map[string]Provider{
	"anthropic": {Name: "anthropic", BaseURL: "https://api.anthropic.com", CompressRequests: false},
	"openai":    {Name: "openai", BaseURL: "https://api.openai.com", CompressRequests: false},
	"gemini":    {Name: "gemini", BaseURL: "https://generativelanguage.googleapis.com", CompressRequests: false},
}

// detectProvider identifies the LLM provider from request headers.
func detectProvider(r *http.Request) (Provider, bool) {
	if r.Header.Get("anthropic-version") != "" {
		return providers["anthropic"], true
	}
	if r.Header.Get("x-goog-api-key") != "" {
		return providers["gemini"], true
	}
	if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		return providers["openai"], true
	}
	return Provider{}, false
}
