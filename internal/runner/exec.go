package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Run executa o comando e retorna a duração em segundos e o código de saída
func Run(ctx context.Context, args []string) (float64, int) {
	if len(args) == 0 {
		return 0, 1
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	// Conecta pipes para manter cores e interatividade
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime).Seconds()

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			// Caso o comando nem seja encontrado
			fmt.Fprintf(os.Stderr, "Erro ao iniciar processo: %v\n", err)
			exitCode = 127
		}
	}

	return duration, exitCode
}
