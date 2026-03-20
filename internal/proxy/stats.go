package proxy

import (
	"fmt"
	"sync/atomic"
)

// Stats tracks bandwidth savings across proxy requests.
type Stats struct {
	totalRequests   atomic.Int64
	originalBytes   atomic.Int64
	compressedBytes atomic.Int64
}

// Record logs a single request's compression result.
func (s *Stats) Record(original, compressed int64) {
	s.totalRequests.Add(1)
	s.originalBytes.Add(original)
	s.compressedBytes.Add(compressed)
}

// Summary returns a formatted bandwidth savings summary.
func (s *Stats) Summary() string {
	reqs := s.totalRequests.Load()
	orig := s.originalBytes.Load()
	comp := s.compressedBytes.Load()
	if orig == 0 {
		return fmt.Sprintf("requests: %d | no data compressed yet", reqs)
	}
	saved := orig - comp
	pct := float64(saved) / float64(orig) * 100
	return fmt.Sprintf("requests: %d | saved: %s of %s (%.0f%% reduction)",
		reqs, formatBytes(saved), formatBytes(orig), pct)
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
