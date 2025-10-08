// Package middleware provides HTTP middleware functions
package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// compressWriter wraps http.ResponseWriter to provide gzip compression
type compressWriter struct {
	w          http.ResponseWriter
	zw         *gzip.Writer
	statusCode int
}

// newCompressWriter creates a new compressWriter
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns the header map
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write writes the data to the connection as part of an HTTP reply
func (c *compressWriter) Write(p []byte) (int, error) {
	if c.statusCode >= 300 || c.statusCode == 0 {
		// For error status codes or before status is set, write directly without compression
		return c.w.Write(p)
	}
	return c.zw.Write(p)
}

// WriteHeader sends an HTTP response header with the provided status code
func (c *compressWriter) WriteHeader(statusCode int) {
	c.statusCode = statusCode
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
		c.w.Header().Add("Vary", "Accept-Encoding")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes the compressWriter
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader wraps io.ReadCloser to provide gzip decompression
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// newCompressReader creates a new compressReader
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read reads up to len(p) bytes into p
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes the compressReader
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// GzipMiddleware provides gzip compression for HTTP responses
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}
		next.ServeHTTP(ow, r)
	})
}
