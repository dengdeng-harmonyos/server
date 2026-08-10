package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// rotatingFileWriter 是一个并发安全的按大小轮转文件写入器。
// app.log 为当前文件，app.log.1 为最近一次轮转文件。
type rotatingFileWriter struct {
	mu       sync.Mutex
	path     string
	maxSize  int64
	maxFiles int
	file     *os.File
	fileSize int64
}

func newRotatingFileWriter(path string, maxSize int64, maxFiles int) (*rotatingFileWriter, error) {
	if path == "" {
		return nil, fmt.Errorf("log file path is empty")
	}
	if maxSize <= 0 {
		return nil, fmt.Errorf("log max size must be positive")
	}
	if maxFiles <= 0 {
		return nil, fmt.Errorf("log max files must be positive")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	w := &rotatingFileWriter{
		path:     path,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}
	if err := w.open(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *rotatingFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.fileSize > 0 && w.fileSize+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := w.file.Write(p)
	w.fileSize += int64(n)
	return n, err
}

func (w *rotatingFileWriter) open() error {
	file, err := os.OpenFile(w.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("stat log file: %w", err)
	}

	w.file = file
	w.fileSize = stat.Size()
	return nil
}

func (w *rotatingFileWriter) rotate() error {
	if err := w.file.Close(); err != nil {
		return fmt.Errorf("close log file before rotation: %w", err)
	}

	oldest := fmt.Sprintf("%s.%d", w.path, w.maxFiles)
	if err := os.Remove(oldest); err != nil && !os.IsNotExist(err) {
		_ = w.open()
		return fmt.Errorf("remove oldest log file: %w", err)
	}

	for i := w.maxFiles - 1; i >= 1; i-- {
		source := fmt.Sprintf("%s.%d", w.path, i)
		target := fmt.Sprintf("%s.%d", w.path, i+1)
		if err := os.Rename(source, target); err != nil && !os.IsNotExist(err) {
			_ = w.open()
			return fmt.Errorf("rotate log file %s: %w", source, err)
		}
	}

	if err := os.Rename(w.path, w.path+".1"); err != nil && !os.IsNotExist(err) {
		_ = w.open()
		return fmt.Errorf("rotate current log file: %w", err)
	}

	if err := w.open(); err != nil {
		return fmt.Errorf("create new log file after rotation: %w", err)
	}
	return nil
}
