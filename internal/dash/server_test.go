package dash

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joekariuki/sandgrouse/internal/proxy"
)

type mockSource struct {
	stats      *proxy.Stats
	startedAt  time.Time
	requestLog *proxy.RequestLog
}

func (m *mockSource) Stats() *proxy.Stats        { return m.stats }
func (m *mockSource) Uptime() time.Duration       { return time.Since(m.startedAt) }
func (m *mockSource) RequestLog() *proxy.RequestLog { return m.requestLog }

func newMockSource() *mockSource {
	s := &proxy.Stats{}
	s.RecordRequest(1000, 1000)
	s.RecordResponse(300, 1000)
	return &mockSource{
		stats:      s,
		startedAt:  time.Now(),
		requestLog: proxy.NewRequestLog(50),
	}
}

func TestStatsEndpoint(t *testing.T) {
	src := newMockSource()
	d := New(":0", src)

	req := httptest.NewRequest("GET", "/api/stats", nil)
	w := httptest.NewRecorder()
	d.handleStats(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var stats statsResponse
	json.NewDecoder(resp.Body).Decode(&stats)

	if stats.AllTime.Requests != 1 {
		t.Errorf("AllTime.Requests = %d, want 1", stats.AllTime.Requests)
	}
	if stats.AllTime.SavingsBytes != 700 {
		t.Errorf("AllTime.SavingsBytes = %d, want 700", stats.AllTime.SavingsBytes)
	}
}

func TestRecentRequestsEndpoint(t *testing.T) {
	src := newMockSource()
	src.requestLog.Add(proxy.RequestEvent{
		Method:       "POST",
		Path:         "/v1/messages",
		Provider:     "anthropic",
		RequestBytes: 5000,
		ResponseWire: 300,
		ResponseOrig: 1000,
	})

	d := New(":0", src)
	req := httptest.NewRequest("GET", "/api/requests/recent", nil)
	w := httptest.NewRecorder()
	d.handleRecentRequests(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var events []proxy.RequestEvent
	json.Unmarshal(body, &events)

	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}
	if events[0].Provider != "anthropic" {
		t.Errorf("Provider = %q, want 'anthropic'", events[0].Provider)
	}
}

func TestIndexEndpoint(t *testing.T) {
	src := newMockSource()
	d := New(":0", src)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	d.handleIndex(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
}
