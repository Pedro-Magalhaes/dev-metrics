package metrics

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Save grava a métrica no arquivo ~/.local/share/build-metrics/build_log.jsonl
func Save(m BuildMetric) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	logDir := filepath.Join(homeDir, ".local", "share", "build-metrics")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logFile := filepath.Join(logDir, "build_log.jsonl")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
