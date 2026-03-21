package dash

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/joekariuki/sandgrouse/internal/proxy"
)

//go:embed static/index.html
var indexHTML []byte

// ProxySource abstracts the proxy data the dashboard needs.
type ProxySource interface {
	Stats() *proxy.Stats
	Uptime() time.Duration
	RequestLog() *proxy.RequestLog
}

// Dashboard is the web dashboard server.
type Dashboard struct {
	listenAddr string
	source     ProxySource
	httpServer *http.Server
}

// New creates a new Dashboard server.
func New(addr string, source ProxySource) *Dashboard {
	return &Dashboard{
		listenAddr: addr,
		source:     source,
	}
}

// Start begins serving the dashboard.
func (d *Dashboard) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/stats", d.handleStats)
	mux.HandleFunc("/api/requests/recent", d.handleRecentRequests)
	mux.HandleFunc("/api/events", d.handleSSE)
	mux.HandleFunc("/", d.handleIndex)

	d.httpServer = &http.Server{
		Addr:    d.listenAddr,
		Handler: mux,
	}

	log.Printf("dashboard at http://localhost%s", d.listenAddr)
	err := d.httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Shutdown gracefully shuts down the dashboard server.
func (d *Dashboard) Shutdown(ctx context.Context) error {
	if d.httpServer == nil {
		return nil
	}
	return d.httpServer.Shutdown(ctx)
}
