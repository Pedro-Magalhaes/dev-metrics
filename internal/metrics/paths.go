package metrics

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const EnvLogPath = "BUILD_METRICS_LOG"

// Usando monkey patching simples para facilitar testes
var EnvGetter = os.Getenv
var HomeDirGetter = os.UserHomeDir
var MkdirAll = os.MkdirAll

// GetLogFilePath retorna o caminho final baseado na prioridade:
// 1. override (se não for vazio)
// 2. Variável de ambiente BUILD_METRICS_LOG
// 3. Padrão ~/.local/share/...
func GetLogFilePath(override string) (string, error) {
	if override != "" {
		return override, nil
	}

	if envPath := EnvGetter(EnvLogPath); envPath != "" {
		return envPath, nil
	}

	homeDir, err := HomeDirGetter()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".local", "share", "build-metrics", "build_log.jsonl"), nil
}

// EnsureLogDir garante que o diretório informado exista.
//
// Importante: este método espera receber um caminho de diretório (ex.: filepath.Dir(logFilePath)).
func EnsureLogDir(dirPath string) error {
	return MkdirAll(dirPath, 0755)
}

// PrintResolvedLogPath resolves the log file path (using same priority as GetLogFilePath)
// and writes a descriptive line to the provided writer. It returns the resolved
// path or an error if resolution failed.
func PrintResolvedLogPath(w io.Writer, label string, override string) (string, error) {
	lp, err := GetLogFilePath(override)
	if err != nil {
		fmt.Fprintf(w, "%s(erro ao resolver: %v)\n", label, err)
		return "", err
	}
	fmt.Fprintf(w, "%s%s\n", label, lp)
	return lp, nil
}
