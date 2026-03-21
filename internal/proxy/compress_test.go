package proxy

import (
	"bytes"
	"io"
	"testing"
)

var testPayload = []byte(`{"model":"claude-sonnet-4-20250514","max_tokens":1024,"messages":[{"role":"user","content":"Hello, how are you doing today? I would like to discuss the implementation of a reverse proxy server that compresses HTTP traffic for developers on metered connections."}]}`)

func TestCompressGzipRoundtrip(t *testing.T) {
	original := testPayload
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

func TestCompressBrotliRoundtrip(t *testing.T) {
	original := testPayload
	compressed, err := compressBrotli(original)
	if err != nil {
		t.Fatalf("compressBrotli() error: %v", err)
	}

	if len(compressed) >= len(original) {
		t.Errorf("compressed size (%d) should be smaller than original (%d)", len(compressed), len(original))
	}

	decompressed, err := decompressBrotli(compressed)
	if err != nil {
		t.Fatalf("decompressBrotli() error: %v", err)
	}

	if string(decompressed) != string(original) {
		t.Errorf("roundtrip failed: got %q, want %q", string(decompressed), string(original))
	}
}

func TestCompressBrotliEmpty(t *testing.T) {
	compressed, err := compressBrotli([]byte{})
	if err != nil {
		t.Fatalf("compressBrotli() error: %v", err)
	}

	decompressed, err := decompressBrotli(compressed)
	if err != nil {
		t.Fatalf("decompressBrotli() error: %v", err)
	}

	if len(decompressed) != 0 {
		t.Errorf("expected empty result, got %d bytes", len(decompressed))
	}
}

func TestCompressDispatch(t *testing.T) {
	original := testPayload
	tests := []struct {
		name      string
		algorithm string
		wantEnc   string
	}{
		{name: "brotli", algorithm: "brotli", wantEnc: "br"},
		{name: "gzip", algorithm: "gzip", wantEnc: "gzip"},
		{name: "default falls back to brotli", algorithm: "unknown", wantEnc: "br"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := compress(original, tt.algorithm)
			if err != nil {
				t.Fatalf("compress() error: %v", err)
			}
			if len(compressed) >= len(original) {
				t.Errorf("compressed size (%d) should be smaller original (%d)", len(compressed), len(original))
			}
			if enc := contentEncoding(tt.algorithm); enc != tt.wantEnc {
				t.Errorf("contentEncoding(%q) = %q, want %q", tt.algorithm, enc, tt.wantEnc)
			}
		})
	}
}

func TestBrotliBeatGzipOnJson(t *testing.T) {
	original := testPayload

	gzipped, _ := compressGzip(original)
	brotlied, _ := compressBrotli(original)

	t.Logf("original: %d bytes, gzip: %d bytes, brotli: %d bytes", len(original), len(gzipped), len(brotlied))

	if len(brotlied) >= len(gzipped) {
		t.Errorf("brotli (%d) should produce smaller output than gzip (%d) on JSON", len(brotlied), len(gzipped))
	}
}

func TestCountingReader(t *testing.T) {
	data := []byte("hello, sandgrouse proxy!")
	cr := &countingReader{reader: bytes.NewReader(data)}

	buf, err := io.ReadAll(cr)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if string(buf) != string(data) {
		t.Errorf("got %q, want %q", string(buf), string(data))
	}
	if cr.bytesRead != int64(len(data)) {
		t.Errorf("bytesRead = %d, want %d", cr.bytesRead, len(data))
	}
}

func TestDecompressReaderGzip(t *testing.T) {
	original := []byte(`{"model":"claude-sonnet-4-20250514","content":"test response data"}`)
	compressed, err := compressGzip(original)
	if err != nil {
		t.Fatalf("compressGzip error: %v", err)
	}

	reader, decompressed, err := decompressReader(bytes.NewReader(compressed), "gzip")
	if err != nil {
		t.Fatalf("decompressReader error: %v", err)
	}
	defer reader.Close()

	if !decompressed {
		t.Error("expected decompressed=true for gzip")
	}

	result, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if string(result) != string(original) {
		t.Errorf("got %q, want %q", string(result), string(original))
	}
}

func TestDecompressReaderBrotli(t *testing.T) {
	original := []byte(`{"model":"claude-sonnet-4-20250514","content":"test response data"}`)
	compressed, err := compressBrotli(original)
	if err != nil {
		t.Fatalf("compressBrotli error: %v", err)
	}

	reader, decompressed, err := decompressReader(bytes.NewReader(compressed), "br")
	if err != nil {
		t.Fatalf("decompressReader error: %v", err)
	}
	defer reader.Close()

	if !decompressed {
		t.Error("expected decompressed=true for br")
	}

	result, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if string(result) != string(original) {
		t.Errorf("got %q, want %q", string(result), string(original))
	}
}

func TestDecompressReaderNoEncoding(t *testing.T) {
	data := []byte(`{"raw":"uncompressed data"}`)

	reader, decompressed, err := decompressReader(bytes.NewReader(data), "")
	if err != nil {
		t.Fatalf("decompressReader error: %v", err)
	}
	defer reader.Close()

	if decompressed {
		t.Error("expected decompressed=false for empty encoding")
	}

	result, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll error: %v", err)
	}
	if string(result) != string(data) {
		t.Errorf("got %q, want %q", string(result), string(data))
	}
}
