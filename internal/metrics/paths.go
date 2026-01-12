package metrics

import (
	"os"
	"path/filepath"
)

const EnvLogPath = "BUILD_METRICS_LOG"

// GetLogFilePath retorna o caminho final baseado na prioridade:
// 1. override (se não for vazio)
// 2. Variável de ambiente BUILD_METRICS_LOG
// 3. Padrão ~/.local/share/...
func GetLogFilePath(override string) (string, error) {
	if override != "" {
		return override, nil
	}

	if envPath := os.Getenv(EnvLogPath); envPath != "" {
		return envPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".local", "share", "build-metrics", "build_log.jsonl"), nil
}

// EnsureLogDir agora recebe o caminho pretendido e garante que a pasta exista
func EnsureLogDir(targetPath string) error {
	dir := filepath.Dir(targetPath)
	return os.MkdirAll(dir, 0755)
}
