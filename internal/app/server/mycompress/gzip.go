package mycompress

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

const (
	successfulMaxCode = 300
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
	n, err := c.ZW.Write(p)
	return n, fmt.Errorf("gzip.go func Write(): error write - %w", err)
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < successfulMaxCode {
		c.W.Header().Set("Content-Encoding", "gzip")
	}
	c.W.WriteHeader(statusCode)
}

func (c *CompressWriter) Close() error {
	return fmt.Errorf("gzip.go func Close(): error close - %w", c.ZW.Close())
}

type CompressReader struct {
	R  io.ReadCloser
	ZR *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("gzip.go func NewCompressReader(): error create reader - %w", err)
	}

	return &CompressReader{
		R:  r,
		ZR: zr,
	}, nil
}

func (c CompressReader) Read(p []byte) (n int, err error) {
	n, err = c.ZR.Read(p)
	return n, fmt.Errorf("gzip.go func Read(): error read - %w", err)
}

func (c *CompressReader) Close() error {
	if err := c.R.Close(); err != nil {
		return fmt.Errorf("gzip.go func Close(): %w", err)
	}
	return fmt.Errorf("gzip.go func Close(): %w", c.ZR.Close())
}
