package util

import (
	"encoding/csv"
	"os"
	"sync"
)

type Writer interface {
	Write(record []string) error
	Flush()
	Error() error
}

type CsvWriter struct {
	mutex  sync.Mutex
	writer *csv.Writer
}

// NewCsvWriter creates a new CsvWriter.
func NewCsvWriter(fileName string) (*CsvWriter, error) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return &CsvWriter{writer: csv.NewWriter(csvFile)}, nil
}

// Write writes a CSV record to the file.
func (w *CsvWriter) Write(record []string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	err := w.writer.Write(record)
	if err != nil {
		return err
	}
	return nil
}

// Flush flushes any buffered data to the underlying file.
func (w *CsvWriter) Flush() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.writer.Flush()
}

// Error returns any error
func (w *CsvWriter) Error() error {
	return w.writer.Error()
}
