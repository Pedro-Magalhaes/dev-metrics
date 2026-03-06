package metrics_test

import (
	"dev-metrics/internal/metrics"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestScanJSONL(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		r       io.Reader
		strict  bool
		fn      func(metrics.BuildMetric) error
		want    metrics.ScanResult
		wantErr bool
	}{
		{
			name:    "Json Vazio",
			r:       strings.NewReader(``),
			strict:  false,
			fn:      func(metrics.BuildMetric) error { return nil },
			want:    metrics.ScanResult{Processed: 0, Skipped: 0},
			wantErr: false,
		},
		{
			name:    "Jsonl com uma linha válida (strict=false)",
			r:       strings.NewReader(`{"project":"A","timestamp":"2026-01-01T00:00:00Z"}`),
			strict:  false,
			fn:      func(metrics.BuildMetric) error { return nil },
			want:    metrics.ScanResult{Processed: 1, Skipped: 0},
			wantErr: false,
		},
		{
			name:    "Jsonl com uma linha válida (strict=false) mas com erro na função de callback",
			r:       strings.NewReader(`{"project":"A","timestamp":"2026-01-01T00:00:00Z"}`),
			strict:  false,
			fn:      func(metrics.BuildMetric) error { return errors.New("callback error") },
			want:    metrics.ScanResult{}, // O resultado é irrelevante, pois esperamos um erro de callback
			wantErr: true,
		},
		{
			name: "Jsonl com múltiplas linhas válidas (strict=false)",
			r: strings.NewReader(`{"project":"A","timestamp":"2026-01-01T00:00:00Z"}
			{"project":"B","timestamp":"2026-01-02T00:00:00Z"}`),
			strict:  false,
			fn:      func(metrics.BuildMetric) error { return nil },
			want:    metrics.ScanResult{Processed: 2, Skipped: 0},
			wantErr: false,
		},
		{
			name: "Jsonl com múltiplas linhas, uma inválida (strict=false)",
			r: strings.NewReader(`{"project":"A","timestamp":"2026-01-01T00:00:00Z"}
			INVALID
			{"project":"B","timestamp":"2026-01-02T00:00:00Z"}`),
			strict:  false,
			fn:      func(metrics.BuildMetric) error { return nil },
			want:    metrics.ScanResult{Processed: 2, Skipped: 1},
			wantErr: false,
		},
		{
			name: "Jsonl com múltiplas linhas válidas (strict=true)",
			r: strings.NewReader(`{"project":"A","timestamp":"2026-01-01T00:00:00Z"}
			{"project":"B","timestamp":"2026-01-02T00:00:00Z"}`),
			strict:  true,
			fn:      func(metrics.BuildMetric) error { return nil },
			want:    metrics.ScanResult{Processed: 2, Skipped: 0},
			wantErr: false,
		},
		{
			name:    "Jsonl com uma linha inválida (strict=true)",
			r:       strings.NewReader(`INVALID`),
			strict:  true,
			fn:      func(metrics.BuildMetric) error { return nil },
			want:    metrics.ScanResult{Processed: 0, Skipped: 1},
			wantErr: true,
		},
		{
			name: "Jsonl com múltiplas linhas, uma inválida (strict=true)",
			r: strings.NewReader(`{"project":"A","timestamp":"2026-01-01T00:00:00Z"}
			INVALID
			{"project":"B","timestamp":"2026-01-02T00:00:00Z"}`),
			strict:  true,
			fn:      func(metrics.BuildMetric) error { return nil },
			want:    metrics.ScanResult{Processed: 1, Skipped: 1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := metrics.ScanJSONL(tt.r, tt.strict, tt.fn)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ScanJSONL() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ScanJSONL() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ScanJSONL() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: Refazer esse teste ruim para duas funções sem uso: JSONLLineError.Error() e JSONLLineError.Unwrap().
func TestJSONLLineError_Error(t *testing.T) {
	tests := []struct {
		name       string // description of this test case
		wantString string
		wantError  error
		instance   metrics.JSONLLineError
	}{
		{
			name:       "Error with short line",
			wantString: `jsonl parse error on line 1: x (raw="INVALID")`,
			wantError:  errors.New("x"),
			instance:   metrics.JSONLLineError{Line: 1, Err: errors.New("x"), Raw: "INVALID"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.instance
			got := e.Error()
			if got != tt.wantString {
				t.Errorf("Error() = %v, want %v", got, tt.wantString)
			}
			if e.Unwrap().Error() != tt.wantError.Error() {
				t.Errorf("Error() = %v, want error %v", e.Unwrap(), tt.wantError)
			}
		})
	}
}
