package rwfile

import (
	"bufio"
	"fmt"
	"os"
)

const (
	filePermission = 0o600
)

type FileWriter struct {
	File *os.File
}

func (c *FileWriter) Close() error {
	err := c.File.Close()
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
		File: file,
	}, nil
}

type FileReader struct {
	file   *os.File
	Reader *bufio.Scanner
}

func NewFileReader(filename string) (*FileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, filePermission)
	if err != nil {
		return nil, fmt.Errorf("restorer.go NewFileReader: %w", err)
	}

	return &FileReader{
		file:   file,
		Reader: bufio.NewScanner(file),
	}, nil
}

func (c *FileReader) Close() error {
	err := c.file.Close()
	if err != nil {
		return fmt.Errorf("restorer.go func (c *FileReader) Close(): %w", err)
	}
	return nil
}
