package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
)

// Stats tracks bandwidth savings across proxy requests.
type Stats struct {
	totalRequests          atomic.Int64
	requestOriginalBytes   atomic.Int64
	requestCompressedBytes atomic.Int64
	responseOriginalBytes  atomic.Int64
	responseWireBytes      atomic.Int64
}

// RecordRequest logs request body sizes.
func (s *Stats) RecordRequest(original, compressed int64) {
	s.totalRequests.Add(1)
	s.requestOriginalBytes.Add(original)
	s.requestCompressedBytes.Add(compressed)
}

// RecordResponse logs response body sizes.
func (s *Stats) RecordResponse(wireBytes, originalBytes int64) {
	s.responseWireBytes.Add(wireBytes)
	s.responseOriginalBytes.Add(originalBytes)
}

// Summary returns a formatted bandwidth savings summary.
func (s *Stats) Summary() string {
	reqs := s.totalRequests.Load()
	respWire := s.responseWireBytes.Load()
	respOrig := s.responseOriginalBytes.Load()
	if respOrig == 0 {
		return fmt.Sprintf("requests: %d | no response data tracked yet", reqs)
	}
	saved := respOrig - respWire
	pct := float64(saved) / float64(respOrig) * 100
	return fmt.Sprintf("requests: %d | responses: %s on wire, %s original (%s saved, %.0f%% reduction)",
		reqs, FormatBytes(respWire), FormatBytes(respOrig), FormatBytes(saved), pct)
}

// TotalRequests returns the total number of proxied requests.
func (s *Stats) TotalRequests() int64 {
	return s.totalRequests.Load()
}

// RequestOriginalBytes returns total original request body bytes.
func (s *Stats) RequestOriginalBytes() int64 {
	return s.requestOriginalBytes.Load()
}

// RequestCompressedBytes returns total compressed request body bytes.
func (s *Stats) RequestCompressedBytes() int64 {
	return s.requestCompressedBytes.Load()
}

// ResponseOriginalBytes returns total decompressed response body bytes.
func (s *Stats) ResponseOriginalBytes() int64 {
	return s.responseOriginalBytes.Load()
}

// ResponseWireBytes returns total on-wire response body bytes.
func (s *Stats) ResponseWireBytes() int64 {
	return s.responseWireBytes.Load()
}

// statsJSON is the JSON-serializable representation of Stats.
type statsJSON struct {
	TotalRequests          int64 `json:"total_requests"`
	RequestOriginalBytes   int64 `json:"request_original_bytes"`
	RequestCompressedBytes int64 `json:"request_compressed_bytes"`
	ResponseOriginalBytes  int64 `json:"response_original_bytes"`
	ResponseWireBytes      int64 `json:"response_wire_bytes"`
}

// SaveTo writes the current stats to a JSON file.
func (s *Stats) SaveTo(path string) error {
	data := statsJSON{
		TotalRequests:          s.totalRequests.Load(),
		RequestOriginalBytes:   s.requestOriginalBytes.Load(),
		RequestCompressedBytes: s.requestCompressedBytes.Load(),
		ResponseOriginalBytes:  s.responseOriginalBytes.Load(),
		ResponseWireBytes:      s.responseWireBytes.Load(),
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// LoadFrom loads stats from a JSON file. Returns nil if the file doesn't exist.
func (s *Stats) LoadFrom(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	var data statsJSON
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	s.totalRequests.Store(data.TotalRequests)
	s.requestOriginalBytes.Store(data.RequestOriginalBytes)
	s.requestCompressedBytes.Store(data.RequestCompressedBytes)
	s.responseOriginalBytes.Store(data.ResponseOriginalBytes)
	s.responseWireBytes.Store(data.ResponseWireBytes)
	return nil
}

// FormatBytes formats a byte count as a human-readable string.
func FormatBytes(b int64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
