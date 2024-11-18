package rwfile

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

const (
	filePermission = 0o600
)

type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func (c *FileWriter) Close() error {
	err := c.file.Close()
	if err != nil {
		return fmt.Errorf("restorer.go Close: %w", err)
	}
	return nil
}

func NewFileWriter(filename string) (*FileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, filePermission)
	if err != nil {
		return nil, fmt.Errorf("restorer.go NewFileWriter: %w", err)
	}

	return &FileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

type FileReader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewFileReader(filename string) (*FileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, filePermission)
	if err != nil {
		return nil, fmt.Errorf("restorer.go NewFileReader: %w", err)
	}

	return &FileReader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *FileReader) Close() error {
	err := c.file.Close()
	if err != nil {
		return fmt.Errorf("restorer.go func (c *FileReader) Close(): %w", err)
	}
	return nil
}
