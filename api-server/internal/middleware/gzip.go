package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type compressResponseWriter struct {
	rw http.ResponseWriter
	zp *gzip.Writer
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// CompressResponseWriter

func newCompressResponseWriter(w http.ResponseWriter) *compressResponseWriter {
	return &compressResponseWriter{
		rw: w,
		zp: gzip.NewWriter(w),
	}
}

func (c *compressResponseWriter) Header() http.Header {
	return c.rw.Header()
}

func (c *compressResponseWriter) Write(p []byte) (int, error) {
	return c.zp.Write(p)
}

func (c *compressResponseWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.rw.Header().Set("Content-Encoding", "gzip")
	}

	c.rw.WriteHeader(statusCode)
}

func (c *compressResponseWriter) Close() error {
	return c.zp.Close()
}

// CompressReader

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

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func СompressionResponseRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := w

		// Write
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {

			cw := newCompressResponseWriter(w)
			response = cw
			defer cw.Close()
		}

		// Read
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				log.Printf("CompressReader error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		handler.ServeHTTP(response, r)
	})
}
