package commands_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"dev-metrics/internal/commands"
)

func TestReportCommand_Aliases(t *testing.T) {
	out := bytes.Buffer{}
	c := &commands.ReportCommand{Out: &out}
	//
	//nolint
	c.Run([]string{"-sssss", "invalid"}) // garante que ensureDefaults é chamado
	got := c.Aliases()

	if len(got) != 0 {
		t.Errorf("Aliases() = %v, want empty", got)
		return
	}
}

func TestReportCommand_Run(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		fileHandler func(name string) (io.ReadCloser, error)
		wantErr     bool
	}{
		{
			name: "File Open Error",
			args: []string{"-log", "nonexistent.log"},
			fileHandler: func(name string) (io.ReadCloser, error) {
				return nil, os.ErrNotExist
			},
			wantErr: true,
		},
		{
			name: "Empty Log File",
			args: []string{"-log", "testdata/sample.log"},
			fileHandler: func(name string) (io.ReadCloser, error) {
				// retorna mock de arquivo
				return &mockReadCloser{Reader: bytes.NewBufferString(""), closeFunc: nil}, nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			fakeFileHandler := func(name string) (io.ReadCloser, error) {
				if tt.fileHandler != nil {
					return tt.fileHandler(name)
				}
				return &mockReadCloser{Reader: bytes.NewBufferString(""), closeFunc: nil}, nil
			}

			c := &commands.ReportCommand{Out: &buf, FileOpener: fakeFileHandler}
			gotErr := c.Run(tt.args)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Run() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Run() succeeded unexpectedly")
			}
		})
	}
}
