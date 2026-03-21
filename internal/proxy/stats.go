package proxy

import (
	"fmt"
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
		reqs, formatBytes(respWire), formatBytes(respOrig), formatBytes(saved), pct)
}

// formatBytes formats a byte count as a human-readable string.
func formatBytes(b int64) string {
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
