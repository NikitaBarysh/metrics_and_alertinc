package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

type compressWriter struct {
	rw http.ResponseWriter
	zw *gzip.Writer

	once         sync.Once
	isCompressed bool
}

func (c *compressWriter) checkContentType() bool {
	c.once.Do(func() {
		for _, v := range c.rw.Header()["Content-Type"] {
			if v == "text/html" || v == "application/json" {
				c.rw.Header().Set("Content-Encoding", "gzip")
				c.isCompressed = true
			}
		}
	})
	return c.isCompressed
}

func newCompressWriter(rw http.ResponseWriter) *compressWriter {
	return &compressWriter{
		rw: rw,
		zw: gzip.NewWriter(rw),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.rw.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	if c.checkContentType() {
		return c.zw.Write(p)
	}
	return c.rw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.rw.Header().Set("Content-Encoding", "gzip")
	}
	c.rw.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	if c.checkContentType() {
		return c.zw.Close()
	}
	return nil
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

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

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ow := rw
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")
		if supportGzip {
			cw := newCompressWriter(rw)
			ow = cw
			defer cw.Close()
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendGzip := strings.Contains(contentEncoding, "gzip")
		if sendGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}
		h.ServeHTTP(ow, r)
	})
}
