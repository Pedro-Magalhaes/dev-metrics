package metrics

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type fakeFile struct {
	content string
	closed  bool
}

func (f *fakeFile) WriteString(s string) (int, error) {
	f.content += s
	return len(s), nil
}

func (f *fakeFile) Close() error {
	f.closed = true
	return nil
}

func TestSave(t *testing.T) {
	// Mocks

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

// Helper
func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
