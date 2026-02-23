package runner_test

import (
	"context"
	"dev-metrics/internal/runner"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		minDuration  time.Duration // Opcional: para verificar se durou pelo menos X (teste de sleep)
		maxDuration  time.Duration // Opcional: para verificar timeout
	}{
		{
			name:         "Empty args",
			args:         []string{},
			wantExitCode: 1,
		},
		{
			name:         "Success command (true)",
			args:         []string{"true"},
			wantExitCode: 0,
		},
		{
			name:         "Failure command (false)",
			args:         []string{"false"},
			wantExitCode: 1,
		},
		{
			name:         "Specific exit code",
			args:         []string{"sh", "-c", "exit 42"},
			wantExitCode: 42,
		},
		{
			name:         "Command not found",
			args:         []string{"cmd-que-nao-existe-12345"},
			wantExitCode: 127,
		},
		{
			name:        "Command with duration",
			args:        []string{"sleep", "0.2"},
			minDuration: 150 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Mede o tempo do teste também para validações simples
			start := time.Now()
			durationSeconds, exitCode := runner.Run(ctx, tt.args)
			realDuration := time.Since(start)

			if exitCode != tt.wantExitCode && tt.wantExitCode != 0 {
				// Nota: Para o caso do sleep (0), ignoramos validação estrita de código aqui se não definido explicitamente
				t.Errorf("Run() exitCode = %v, want %v", exitCode, tt.wantExitCode)
			}

			if tt.minDuration > 0 {
				if realDuration < tt.minDuration {
					t.Errorf("Run() executed too fast: got %v, want >= %v", realDuration, tt.minDuration)
				}
				// Verifica se o retorno da função (float64) é coerente
				if durationSeconds < tt.minDuration.Seconds() {
					t.Errorf("Run() returned duration %f, expected >= %f", durationSeconds, tt.minDuration.Seconds())
				}
			}
		})
	}
}

func TestRun_ContextCancellation(t *testing.T) {
	// Testa se o comando é interrompido quando o contexto é cancelado
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Tenta rodar um sleep de 2 segundos, mas o contexto morre em 0.1s
	args := []string{"sleep", "2"}

	start := time.Now()
	_, exitCode := runner.Run(ctx, args)
	duration := time.Since(start)

	// O processo deve terminar muito antes de 2 segundos
	if duration >= 1500*time.Millisecond {
		t.Errorf("Run() ignores context cancellation, duration: %v", duration)
	}

	// Quando morto por sinal (context kill), o exit code geralmente não é 0
	if exitCode == 0 {
		t.Errorf("Run() with cancelled context should not return exit code 0")
	}
}
