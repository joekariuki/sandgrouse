package proxy

import (
	"strings"
	"testing"
)

func TestStatsRecord(t *testing.T) {
	s := &Stats{}

	s.Record(1000, 300)
	s.Record(2000, 500)

	if got := s.totalRequests.Load(); got != 2 {
		t.Errorf("totalRequests = %d, want 2", got)
	}
	if got := s.originalBytes.Load(); got != 3000 {
		t.Errorf("originalBytes = %d, want 3000", got)
	}
	if got := s.compressedBytes.Load(); got != 800 {
		t.Errorf("compressedBytes = %d, want 800", got)
	}
}

func TestStatsSummary(t *testing.T) {
	tests := []struct {
		name      string
		records   [][2]int64 // {original, compressed}
		wantParts []string
	}{
		{
			name:      "no requests",
			records:   nil,
			wantParts: []string{"requests: 0", "no data compressed yet"},
		},
		{
			name:      "single request",
			records:   [][2]int64{{1000, 300}},
			wantParts: []string{"requests: 1", "70% reduction"},
		},
		{
			name:      "multiple requests",
			records:   [][2]int64{{1000, 300}, {2000, 500}},
			wantParts: []string{"requests: 2", "73% reduction"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stats{}
			for _, r := range tt.records {
				s.Record(r[0], r[1])
			}
			summary := s.Summary()
			for _, part := range tt.wantParts {
				if !strings.Contains(summary, part) {
					t.Errorf("Summary() = %q, want it to contain %q", summary, part)
				}
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := formatBytes(tt.input); got != tt.want {
				t.Errorf("formatBytes(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
