package metrics_test

import (
	"dev-metrics/internal/metrics"
	"errors"
	"testing"
)

// Função de teste para GetLogFilePath
func MockGetenvBuilder(returnString string) func(string) string {
	return func(key string) string {
		return returnString
	}
}

// Função de teste para GetLogFilePath
func MockHomeDirFuncBuilder(returnString string) func() (string, error) {
	return func() (string, error) {
		if returnString == "error" {
			return "", errors.New("failed to get home dir")
		}
		return returnString, nil
	}
}

func TestGetLogFilePath(t *testing.T) {
	tests := []struct {
		name       string
		definedEnv string
		homeDir    string
		override   string
		want       string
		wantErr    bool
	}{
		{
			name:       "Override path provided",
			definedEnv: "/path/from/env.jsonl",
			homeDir:    "/home",
			override:   "/custom/path/log.jsonl",
			want:       "/custom/path/log.jsonl",
			wantErr:    false,
		},
		{
			name:       "Environment variable path used",
			definedEnv: "/path/from/env.jsonl",
			homeDir:    "/home",
			override:   "",
			want:       "/path/from/env.jsonl",
			wantErr:    false,
		},
		{
			name:       "Default path used when no override or env var",
			override:   "",
			definedEnv: "",
			homeDir:    "/home",
			want:       "/home/.local/share/build-metrics/build_log.jsonl",
			wantErr:    false,
		},
		{
			name:       "Default path used when no override or env var",
			override:   "",
			definedEnv: "",
			homeDir:    "error", // Simula erro ao obter diretório home
			want:       "/home/.local/share/build-metrics/build_log.jsonl",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics.EnvGetter = MockGetenvBuilder(tt.definedEnv)
			metrics.HomeDirGetter = MockHomeDirFuncBuilder(tt.homeDir)
			got, gotErr := metrics.GetLogFilePath(tt.override)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetLogFilePath() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetLogFilePath() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("GetLogFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
