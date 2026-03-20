package proxy

import (
	"testing"
)

func TestCompressGzipRoundtrip(t *testing.T) {
	original := []byte(`{"model":"claude-sonnet-4-20250514","max_tokens":1024,"messages":[{"role":"user","content":"Hello, how are you doing today? I would like to discuss the implementation of a reverse proxy server that compresses HTTP traffic for developers on metered connections."}]}`)

	compressed, err := compressGzip(original)
	if err != nil {
		t.Fatalf("compressGzip() error: %v", err)
	}

	if len(compressed) >= len(original) {
		t.Errorf("compressed size (%d) should be smaller than original (%d)", len(compressed), len(original))
	}

	decompressed, err := decompressGzip(compressed)
	if err != nil {
		t.Fatalf("decompressGzip() error: %v", err)
	}

	if string(decompressed) != string(original) {
		t.Errorf("roundtrip failed: got %q, want %q", string(decompressed), string(original))
	}
}

func TestCompressGzipEmpty(t *testing.T) {
	compressed, err := compressGzip([]byte{})
	if err != nil {
		t.Fatalf("compressGzip() error: %v", err)
	}

	decompressed, err := decompressGzip(compressed)
	if err != nil {
		t.Fatalf("decompressGzip() error: %v", err)
	}

	if len(decompressed) != 0 {
		t.Errorf("expected empty result, got %d bytes", len(decompressed))
	}
}
