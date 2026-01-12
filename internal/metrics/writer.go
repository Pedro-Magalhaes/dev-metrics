package metrics

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// fileWriter is the subset of *os.File used by Save.
type fileWriter interface {
	WriteString(string) (int, error)
	Close() error
}

// package-level variables so tests can mock filesystem operations.
var (
	openFile = func(name string, flag int, perm os.FileMode) (fileWriter, error) {
		return os.OpenFile(name, flag, perm)
	}
)

// Save writes the metric to the specified file path in JSONL format
func Save(m BuildMetric, filePath string) error {
	logDir := filepath.Dir(filePath)
	if err := EnsureLogDir(logDir); err != nil {
		return err
	}

	f, err := openFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = f.WriteString(string(jsonData) + "\n")
	return err
}
