package proxy

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/andybalholm/brotli"
)

// compressGzip compresses data using gzip.
func compressGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// decompressGzip decompresses gzip data.
func decompressGzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

// compressBrotli compresses data using brotli.
func compressBrotli(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := brotli.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// decompressBrotli decompresses brotli data.
func decompressBrotli(data []byte) ([]byte, error) {
	r := brotli.NewReader(bytes.NewReader(data))
	return io.ReadAll(r)
}

// compress compresses data using the specified algorithm.
func compress(data []byte, algorithm string) ([]byte, error) {
	switch algorithm {
	case "brotli":
		return compressBrotli(data)
	case "gzip":
		return compressGzip(data)
	default:
		return compressBrotli(data)
	}
}

// contentEncoding returns the Content-Encoding header value for the algorithm.
func contentEncoding(algorithm string) string {
	switch algorithm {
	case "brotli":
		return "br"
	case "gzip":
		return "gzip"
	default:
		return "br"
	}
}

// countingReader wraps an io.Reader and counts bytes read through it.
type countingReader struct {
	reader    io.Reader
	bytesRead int64
}

func (cr *countingReader) Read(p []byte) (int, error) {
	n, err := cr.reader.Read(p)
	cr.bytesRead += int64(n)
	return n, err
}

// decompressReader wraps a reader with decompression based on Content-Encoding.
// Returns the decompressing reader and a boolean indicating whether decompression was applied.
// The caller should close the returned ReadCloser when done.
func decompressReader(r io.Reader, encoding string) (io.ReadCloser, bool, error) {
	switch encoding {
	case "gzip":
		gr, err := gzip.NewReader(r)
		if err != nil {
			return nil, false, err
		}
		return gr, true, nil
	case "br":
		return io.NopCloser(brotli.NewReader(r)), true, nil
	default:
		return io.NopCloser(r), false, nil
	}
}
