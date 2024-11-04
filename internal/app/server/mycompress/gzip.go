package mycompress

import (
	"compress/gzip"
	"io"
	"net/http"
)

type CompressWriter struct {
	W  http.ResponseWriter
	ZW *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		W:  w,
		ZW: gzip.NewWriter(w),
	}
}

func (c *CompressWriter) Header() http.Header {
	return c.W.Header()
}

func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.ZW.Write(p)
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.W.Header().Set("Content-Encoding", "gzip")
	}
	c.W.WriteHeader(statusCode)
}

func (c *CompressWriter) Close() error {
	return c.ZW.Close()
}

type CompressReader struct {
	R  io.ReadCloser
	ZR *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		R:  r,
		ZR: zr,
	}, nil
}

func (c CompressReader) Read(p []byte) (n int, err error) {
	return c.ZR.Read(p)
}

func (c *CompressReader) Close() error {
	if err := c.R.Close(); err != nil {
		return err
	}
	return c.ZR.Close()
}
