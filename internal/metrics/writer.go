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
	userHomeDir = os.UserHomeDir
	mkdirAll    = os.MkdirAll
	openFile    = func(name string, flag int, perm os.FileMode) (fileWriter, error) {
		return os.OpenFile(name, flag, perm)
	}
)

// Save grava a métrica no arquivo ~/.local/share/build-metrics/build_log.jsonl
func Save(m BuildMetric) error {
	homeDir, err := userHomeDir()
	if err != nil {
		return err
	}

	logDir := filepath.Join(homeDir, ".local", "share", "build-metrics")
	if err := mkdirAll(logDir, 0755); err != nil {
		return err
	}

	logFile := filepath.Join(logDir, "build_log.jsonl")
	f, err := openFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
