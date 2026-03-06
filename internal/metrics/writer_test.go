package metrics

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

// Não rodar em paralelo pois estamos mockando variáveis globais
type fakeFile struct {
	content  string
	closed   bool
	closeErr error
}

func (f *fakeFile) WriteString(s string) (int, error) {
	f.content += s
	return len(s), nil
}

func (f *fakeFile) Close() error {
	f.closed = true
	return f.closeErr
}

func TestSave(t *testing.T) {
	origOpenFile := OpenFile
	origEnsureDir := EnsureDir
	defer func() {
		OpenFile = origOpenFile
		EnsureDir = origEnsureDir
	}()

	var fakeF *fakeFile
	OpenFile = func(name string, flag int, perm os.FileMode) (fileWriter, error) {
		fakeF = &fakeFile{}
		return fakeF, nil
	}

	calledEnsureDir := false
	EnsureDir = func(path string) error {
		calledEnsureDir = true
		if path == "/fail/dir" {
			return fmt.Errorf("fail dir")
		}
		return nil
	}

	tests := []struct {
		name     string
		m        BuildMetric
		filePath string
		wantErr  bool
		wantDir  string
	}{
		{
			name:     "Sucesso salva JSONL",
			m:        BuildMetric{Project: "A", Timestamp: "2026-01-01T00:00:00Z"},
			filePath: "/tmp/test.jsonl",
			wantErr:  false,
			wantDir:  "/tmp",
		},
		{
			name:     "Erro em ensureDir",
			m:        BuildMetric{Project: "B", Timestamp: "2026-01-02T00:00:00Z"},
			filePath: "/fail/dir/test.jsonl",
			wantErr:  true,
			wantDir:  "/fail/dir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calledEnsureDir = false
			fakeF = nil
			gotErr := Save(tt.m, tt.filePath)
			if !calledEnsureDir {
				t.Errorf("ensureDir não foi chamado")
			}
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Save() falhou: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Save() deveria falhar, mas não falhou")
			}
			if fakeF == nil || !fakeF.closed {
				t.Errorf("Arquivo não foi fechado corretamente")
			}
			// Verifica se o JSON foi salvo corretamente
			if tt.wantDir != "/fail/dir" && !contains(fakeF.content, tt.m.Project) {
				t.Errorf("Conteúdo salvo não contém o projeto: %s", tt.m.Project)
			}
		})
	}
}

func TestSave_CloseError(t *testing.T) {
	origOpenFile := OpenFile
	origEnsureDir := EnsureDir
	origLogWriter := log.Default().Writer()
	defer func() {
		OpenFile = origOpenFile
		EnsureDir = origEnsureDir
		log.SetOutput(origLogWriter)
	}()
	// Redireciona o log para um buffer
	var buf bytes.Buffer
	log.SetOutput(&buf)

	OpenFile = func(name string, flag int, perm os.FileMode) (fileWriter, error) {
		return &fakeFile{closeErr: fmt.Errorf("close error")}, nil
	}
	EnsureDir = func(path string) error { return nil }

	err := Save(BuildMetric{Project: "X", Timestamp: "2026-01-01T00:00:00Z"}, "/tmp/test.jsonl")
	if err != nil {
		t.Errorf("Save() retornou erro inesperado: %v", err)
	}

	// Verifica se o erro foi logado
	logOutput := buf.String()
	if !strings.Contains(logOutput, "close error") {
		t.Errorf("Esperava mensagem de erro no log, mas foi: %q", logOutput)
	}
}

// Helper
func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
