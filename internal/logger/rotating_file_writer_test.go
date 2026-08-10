package logger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRotatingFileWriterRotatesAndKeepsConfiguredFiles(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "app.log")
	w, err := newRotatingFileWriter(logPath, 10, 2)
	if err != nil {
		t.Fatalf("newRotatingFileWriter() error = %v", err)
	}

	for _, content := range []string{"12345678", "abcdefgh", "ABCDEFGH"} {
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatalf("Write() error = %v", err)
		}
	}

	assertFileContent(t, logPath, "ABCDEFGH")
	assertFileContent(t, logPath+".1", "abcdefgh")
	assertFileContent(t, logPath+".2", "12345678")
}

func assertFileContent(t *testing.T, path string, expected string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	if string(content) != expected {
		t.Fatalf("ReadFile(%q) = %q, want %q", path, content, expected)
	}
}
