package proxy

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestStatsRecordRequest(t *testing.T) {
	s := &Stats{}

	s.RecordRequest(1000, 1000) // no compression (APIs reject it)
	s.RecordRequest(2000, 2000)

	if got := s.totalRequests.Load(); got != 2 {
		t.Errorf("totalRequests = %d, want 2", got)
	}
	if got := s.requestOriginalBytes.Load(); got != 3000 {
		t.Errorf("requestOriginalBytes = %d, want 3000", got)
	}
}

func TestStatsRecordResponse(t *testing.T) {
	s := &Stats{}

	s.RecordResponse(300, 1000) // 300 bytes on wire, 1000 decompressed
	s.RecordResponse(500, 2000)

	if got := s.responseWireBytes.Load(); got != 800 {
		t.Errorf("responseWireBytes = %d, want 800", got)
	}
	if got := s.responseOriginalBytes.Load(); got != 3000 {
		t.Errorf("responseOriginalBytes = %d, want 3000", got)
	}
}

func TestStatsSummary(t *testing.T) {
	tests := []struct {
		name      string
		responses [][2]int64 // {wireBytes, originalBytes}
		wantParts []string
	}{
		{
			name:      "no responses",
			responses: nil,
			wantParts: []string{"requests: 0", "no response data tracked yet"},
		},
		{
			name:      "single response with savings",
			responses: [][2]int64{{300, 1000}},
			wantParts: []string{"requests: 0", "70% reduction"},
		},
		{
			name:      "multiple responses",
			responses: [][2]int64{{300, 1000}, {500, 2000}},
			wantParts: []string{"73% reduction"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stats{}
			for _, r := range tt.responses {
				s.RecordResponse(r[0], r[1])
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
			if got := FormatBytes(tt.input); got != tt.want {
				t.Errorf("FormatBytes(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestStatsSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stats.json")

	// Create stats with data
	s := &Stats{}
	s.RecordRequest(1000, 800)
	s.RecordRequest(2000, 1600)
	s.RecordResponse(300, 1000)
	s.RecordResponse(500, 2000)

	// Save
	if err := s.SaveTo(path); err != nil {
		t.Fatalf("SaveTo() error: %v", err)
	}

	// Load into fresh stats
	s2 := &Stats{}
	if err := s2.LoadFrom(path); err != nil {
		t.Fatalf("LoadFrom() error: %v", err)
	}

	// Verify all fields match
	if s2.TotalRequests() != 2 {
		t.Errorf("TotalRequests() = %d, want 2", s2.TotalRequests())
	}
	if s2.RequestOriginalBytes() != 3000 {
		t.Errorf("RequestOriginalBytes() = %d, want 3000", s2.RequestOriginalBytes())
	}
	if s2.RequestCompressedBytes() != 2400 {
		t.Errorf("RequestCompressedBytes() = %d, want 2400", s2.RequestCompressedBytes())
	}
	if s2.ResponseOriginalBytes() != 3000 {
		t.Errorf("ResponseOriginalBytes() = %d, want 3000", s2.ResponseOriginalBytes())
	}
	if s2.ResponseWireBytes() != 800 {
		t.Errorf("ResponseWireBytes() = %d, want 800", s2.ResponseWireBytes())
	}
}

func TestStatsLoadMissingFile(t *testing.T) {
	s := &Stats{}
	err := s.LoadFrom(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err != nil {
		t.Errorf("LoadFrom() should return nil for missing file, got: %v", err)
	}
	if s.TotalRequests() != 0 {
		t.Errorf("TotalRequests() = %d, want 0 after loading missing file", s.TotalRequests())
	}
}

func TestStatsGetters(t *testing.T) {
	s := &Stats{}
	s.RecordRequest(1000, 800)
	s.RecordRequest(2000, 1600)
	s.RecordResponse(300, 1000)
	s.RecordResponse(500, 2000)

	if got := s.TotalRequests(); got != 2 {
		t.Errorf("TotalRequests() = %d, want 2", got)
	}
	if got := s.RequestOriginalBytes(); got != 3000 {
		t.Errorf("RequestOriginalBytes() = %d, want 3000", got)
	}
	if got := s.RequestCompressedBytes(); got != 2400 {
		t.Errorf("RequestCompressedBytes() = %d, want 2400", got)
	}
	if got := s.ResponseOriginalBytes(); got != 3000 {
		t.Errorf("ResponseOriginalBytes() = %d, want 3000", got)
	}
	if got := s.ResponseWireBytes(); got != 800 {
		t.Errorf("ResponseWireBytes() = %d, want 800", got)
	}
}
