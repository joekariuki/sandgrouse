package proxy

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// hopByHopHeaders are headers that must not be forwarded by proxies (RFC 9110 §7.6.1).
var hopByHopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

// Server is the sandgrouse proxy server.
type Server struct {
	ListenAddr string
	Algorithm  string // "gzip" or "brotli" (default: "brotli")
	client     *http.Client
	stats      *Stats
	requestLog *RequestLog
	httpServer *http.Server
	startedAt  time.Time
}

// Start begins listening for HTTP requests and forwarding to upstream APIs.
func (s *Server) Start() error {
	s.client = &http.Client{}
	if s.stats == nil {
		s.stats = &Stats{}
	}
	s.requestLog = NewRequestLog(50)
	if s.Algorithm == "" {
		s.Algorithm = "brotli"
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if s.stats.totalRequests.Load() > 0 {
				log.Printf("[stats] %s", s.stats.Summary())
			}
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleProxy)

	s.httpServer = &http.Server{
		Addr:    s.ListenAddr,
		Handler: mux,
	}
	s.startedAt = time.Now()

	log.Printf("sandgrouse proxy listening on %s (compression: %s)", s.ListenAddr, s.Algorithm)
	err := s.httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Shutdown gracefully shuts down the proxy server, draining in-flight requests.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Uptime returns how long the server has been running.
func (s *Server) Uptime() time.Duration {
	return time.Since(s.startedAt)
}

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	reqStart := time.Now()
	provider, ok := detectProvider(r)

	if !ok {
		http.Error(w, "unknown provider: set anthropic-version or Authorization header", http.StatusBadRequest)
		return
	}

	// Build upstream URL
	upstream, err := url.Parse(provider.BaseURL)
	if err != nil {
		http.Error(w, "invalid upstream URL", http.StatusInternalServerError)
		return
	}
	upstream.Path = r.URL.Path
	upstream.RawQuery = r.URL.RawQuery

	// Read the original request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}

	// Compress the request body
	originalSize := len(body)
	var outBody []byte
	outBody = body
	compressed := false
	if originalSize > 0 && provider.CompressRequests {
		outBody, err = compress(body, s.Algorithm)
		if err != nil {
			log.Printf("compression failed, sending uncompressed: %v", err)
			outBody = body
		} else {
			compressed = true
		}
	}

	// Create outgoing request with compressed body
	outReq, err := http.NewRequestWithContext(r.Context(), r.Method, upstream.String(), bytes.NewReader(outBody))
	if err != nil {
		http.Error(w, "failed to create upstream request", http.StatusInternalServerError)
		return
	}

	// Copy all headers from the original request
	for key, values := range r.Header {
		for _, v := range values {
			outReq.Header.Add(key, v)
		}
	}

	// Strip hop by hop headers that must not be forwarded
	for _, h := range hopByHopHeaders {
		outReq.Header.Del(h)
	}

	// Set compression headers
	if compressed {
		outReq.Header.Set("Content-Encoding", contentEncoding(s.Algorithm))
		outReq.ContentLength = int64(len(outBody))
	}

	// Request compressed responses from upstream
	outReq.Header.Set("Accept-Encoding", "gzip, br")

	log.Printf("%s %s -> %s | request: %s",
		r.Method, r.URL.Path, upstream.String(),
		FormatBytes(int64(originalSize)))

	// Record bandwidth stats
	s.stats.RecordRequest(int64(originalSize), int64(len(outBody)))

	// Send request to upstream
	resp, err := s.client.Do(outReq)
	if err != nil {
		log.Printf("upstream error: %v", err)
		http.Error(w, "upstream request failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Wrap response body in counting reader to measure wire bytes
	wireCounter := &countingReader{reader: resp.Body}

	// Get the content encoding from upstream response
	respEncoding := resp.Header.Get("Content-Encoding")

	// Copy response headers back to client
	for key, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}

	// If upstream sent compressed response, decompress before forwarding to client
	// (localhost leg is free, no need to keep it compressed)
	if respEncoding != "" {
		w.Header().Del("Content-Encoding")
		w.Header().Del("Content-Length") // length changes after decompression
	}

	w.WriteHeader(resp.StatusCode)

	// Stream SSE responses with immediate flushing, buffer everything else
	var respWire, respOrig int64
	if isSSE(resp) {
		streamResponse(wireCounter, w)
		respWire = wireCounter.bytesRead
		respOrig = wireCounter.bytesRead
	} else {
		decompReader, wasCompressed, err := decompressReader(wireCounter, respEncoding)
		if err != nil {
			log.Printf("response decompression failed: %v", err)
			io.Copy(w, wireCounter)
			respWire = wireCounter.bytesRead
			respOrig = wireCounter.bytesRead
		} else {
			defer decompReader.Close()
			written, _ := io.Copy(w, decompReader)
			respWire = wireCounter.bytesRead
			respOrig = written
			if wasCompressed {
				log.Printf("  response: %s on wire -> %s decompressed (%.0f%% saved)",
					FormatBytes(respWire), FormatBytes(respOrig),
					compressionRatio(int(respOrig), int(respWire)))
			}
		}
	}

	s.stats.RecordResponse(respWire, respOrig)
	if s.requestLog == nil {
		return
	}
	s.requestLog.Add(RequestEvent{
		Timestamp:    reqStart,
		Method:       r.Method,
		Path:         r.URL.Path,
		Provider:     provider.Name,
		RequestBytes: int64(originalSize),
		ResponseWire: respWire,
		ResponseOrig: respOrig,
	})
}

// Stats returns the server's bandwidth statistics.
func (s *Server) Stats() *Stats {
	return s.stats
}

// SetStats sets pre-loaded stats on the server (must be called before Start).
func (s *Server) SetStats(stats *Stats) {
	s.stats = stats
}

// RequestLog returns the server's request event log.
func (s *Server) RequestLog() *RequestLog {
	return s.requestLog
}

// compressionRatio calculate the percentage reduction.
func compressionRatio(original, compressed int) float64 {
	if original == 0 {
		return 0
	}
	return (1 - float64(compressed)/float64(original)) * 100
}

// isSSE returns true if the response is a server-sent events stream.
func isSSE(resp *http.Response) bool {
	ct := resp.Header.Get("Content-Type")
	return strings.HasPrefix(ct, "text/event-stream")
}

// streamResponse copies an SSE stream to the client, flushing after each chunk.
func streamResponse(src io.Reader, dst http.ResponseWriter) {
	flusher, ok := dst.(http.Flusher)
	if !ok {
		// ResponseWriter doesn't support flushing, fall back to buffered copy
		io.Copy(dst, src)
		return
	}

	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			dst.Write(buf[:n])
			flusher.Flush()
		}
		if err != nil {
			return
		}
	}
}
