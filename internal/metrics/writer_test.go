package metrics

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type fakeFile struct {
	data []byte
}

func (f *fakeFile) WriteString(s string) (int, error) {
	f.data = append(f.data, []byte(s)...)
	return len(s), nil
}

func (f *fakeFile) Close() error { return nil }

func TestSaveWritesJSONL_WithMocks(t *testing.T) {
	// save originals and restore after
	origUserHome := userHomeDir
	origMkdirAll := mkdirAll
	origOpenFile := openFile
	defer func() {
		userHomeDir = origUserHome
		mkdirAll = origMkdirAll
		openFile = origOpenFile
	}()

	var mkdirPath string
	var openedName string
	var ff *fakeFile

	userHomeDir = func() (string, error) { return "/fake/home", nil }
	mkdirAll = func(path string, perm os.FileMode) error {
		mkdirPath = path
		return nil
	}
	openFile = func(name string, flag int, perm os.FileMode) (fileWriter, error) {
		openedName = name
		ff = &fakeFile{}
		return ff, nil
	}

	m := BuildMetric{User: "tester", Command: "make build", DurationSec: 1.23}

	if err := Save(m); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	wantDir := filepath.Join("/fake/home", ".local", "share", "build-metrics")
	if mkdirPath != wantDir {
		t.Fatalf("mkdir called with %q, want %q", mkdirPath, wantDir)
	}

	wantFile := filepath.Join(wantDir, "build_log.jsonl")
	if openedName != wantFile {
		t.Fatalf("openFile called with %q, want %q", openedName, wantFile)
	}

	// check written content
	if ff == nil {
		t.Fatalf("fake file not created")
	}
	// remove trailing newline
	raw := bytes.TrimSuffix(ff.data, []byte("\n"))
	var got BuildMetric
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal written data: %v", err)
	}
	if got.User != m.User || got.Command != m.Command || got.DurationSec != m.DurationSec {
		t.Fatalf("mismatch: got %+v want %+v", got, m)
	}
}

func TestSaveAppends_WithMocks(t *testing.T) {
	origUserHome := userHomeDir
	origMkdirAll := mkdirAll
	origOpenFile := openFile
	defer func() {
		userHomeDir = origUserHome
		mkdirAll = origMkdirAll
		openFile = origOpenFile
	}()

	userHomeDir = func() (string, error) { return "/fake/home", nil }
	mkdirAll = func(path string, perm os.FileMode) error { return nil }

	// return the same file object to simulate appending to same file
	ff := &fakeFile{}
	openFile = func(name string, flag int, perm os.FileMode) (fileWriter, error) {
		return ff, nil
	}

	m := BuildMetric{User: "tester", Command: "make build"}

	if err := Save(m); err != nil {
		t.Fatalf("first Save returned error: %v", err)
	}
	if err := Save(m); err != nil {
		t.Fatalf("second Save returned error: %v", err)
	}

	// split lines
	raw := bytes.TrimSuffix(ff.data, []byte("\n"))
	lines := bytes.Split(raw, []byte("\n"))
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines after two Save calls, got %d", len(lines))
	}

	for i, ln := range lines {
		var got BuildMetric
		if err := json.Unmarshal(ln, &got); err != nil {
			t.Fatalf("unmarshal line %d: %v", i, err)
		}
		if got.User != m.User || got.Command != m.Command {
			t.Fatalf("line %d mismatch: got %+v want %+v", i, got, m)
		}
	}
}
