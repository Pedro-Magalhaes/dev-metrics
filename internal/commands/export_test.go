package commands_test

import (
	"bytes"
	"dev-metrics/internal/commands"
	"dev-metrics/internal/metrics"
	"errors"
	"io"
	"testing"
)

type mockReadCloser struct {
	io.Reader
	closeFunc func() error
}

func (m *mockReadCloser) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

type mockWriteCloser struct {
	io.Writer
	closeFunc func() error
}

func (m *mockWriteCloser) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func TestExportCommand_Run(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		mockOpenErr   error
		mockCreateErr error
		mockExportErr error
		exportResult  metrics.ScanResult
		wantErr       bool
		wantStderr    string
	}{
		{
			name:         "Success to Stdout",
			args:         []string{"-out", "-", "-log", "test.jsonl"},
			exportResult: metrics.ScanResult{Processed: 10, Skipped: 0},
			wantErr:      false,
			wantStderr:   "exportado: 10 linhas (puladas: 0)\n",
		},
		{
			name:         "Success to File",
			args:         []string{"-out", "output.csv", "-log", "test.jsonl"},
			exportResult: metrics.ScanResult{Processed: 5, Skipped: 2},
			wantErr:      false,
			wantStderr:   "exportado: 5 linhas (puladas: 2) -> output.csv\n",
		},
		{
			name:        "Open Error",
			args:        []string{"-log", "missing.jsonl"},
			mockOpenErr: errors.New("file not found"),
			wantErr:     true,
		},
		{
			name:          "Create Error",
			args:          []string{"-out", "/invalid/path.csv", "-log", "test.jsonl"},
			mockCreateErr: errors.New("permission denied"),
			wantErr:       true,
		},
		{
			name:          "Export Error",
			args:          []string{"-out", "-", "-log", "test.jsonl"},
			mockExportErr: errors.New("invalid json"),
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer

			c := &commands.ExportCommand{
				Out: &stdout,
				Err: &stderr,
				MetricsOpener: func(name string) (io.ReadCloser, error) {
					if tt.mockOpenErr != nil {
						return nil, tt.mockOpenErr
					}
					return &mockReadCloser{Reader: bytes.NewBufferString(""), closeFunc: nil}, nil
				},
				FileCreator: func(name string) (io.WriteCloser, error) {
					if tt.mockCreateErr != nil {
						return nil, tt.mockCreateErr
					}
					return &mockWriteCloser{Writer: &bytes.Buffer{}, closeFunc: nil}, nil
				},
				MetricsSaver: func(in io.Reader, out io.Writer, strict bool) (metrics.ScanResult, error) {
					if tt.mockExportErr != nil {
						return metrics.ScanResult{}, tt.mockExportErr
					}
					return tt.exportResult, nil
				},
			}

			err := c.Run(tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.wantStderr != "" {
				if stderr.String() != tt.wantStderr {
					t.Errorf("Stderr = %q, want %q", stderr.String(), tt.wantStderr)
				}
			}
		})
	}
}

func TestExportCommand_Metadata(t *testing.T) {
	c := &commands.ExportCommand{}
	if c.Name() != "export" {
		t.Errorf("Name() = %q, want %q", c.Name(), "export")
	}
	if c.Description() == "" {
		t.Errorf("Description() is empty")
	}
}
