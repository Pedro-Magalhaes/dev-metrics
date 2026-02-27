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
	OpenFile = func(name string, flag int, perm os.FileMode) (fileWriter, error) {
		return os.OpenFile(name, flag, perm)
	}
	EnsureDir = func(path string) error {
		return EnsureLogDir(path)
	}
)

// Save writes the metric to the specified file path in JSONL format
func Save(m BuildMetric, filePath string) error {
	logDir := filepath.Dir(filePath)
	if err := EnsureDir(logDir); err != nil {
		return err
	}

	// Usando append + JSONL obtemos escrita atômica sem problemas de concorrência, dispensando locks adicionais.
	// A atomicidade é garantida se a string for escrita em uma única chamada de WriteString com tamanho menor que o page size do sistema (geralmente 4KB).
	f, err := OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
