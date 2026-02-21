package metrics_test

import (
	"dev-metrics/internal/metrics"
	"io"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// mock reader that implements io.Reader
type mockReader struct {
	data []byte
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if len(m.data) == 0 {
		return 0, io.EOF
	}
	n = copy(p, m.data)
	m.data = m.data[n:]
	return n, nil
}

// mock writer that implements io.Writer
type mockWriter struct {
	data []byte
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	m.data = append(m.data, p...)
	return len(p), nil
}

type exportResult struct {
	metrics.ScanResult
	csvData []string // The CSV data that should have been written to the writer.
}

func TestExportCSVFromJSONL(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		r      io.Reader
		strict bool
		want   struct {
			metrics.ScanResult
			csvData []string // The CSV data that should have been written to the writer.
		}
		wantErr bool
	}{
		{
			name:   "one valid JSONL entry",
			r:      &mockReader{data: []byte(`{"project":"example-project","timestamp":"2024-01-01T00:00:00Z","user":"xpto","hostname":"localhost","os":"linux","branch":"b1","commit":"E2x40","command":"cmake","duration_sec":120.2,"return_code":0,"cpus":4,"status":"completed"}`)},
			strict: false,
			want: exportResult{
				ScanResult: metrics.ScanResult{
					Processed: 1,
					Skipped:   0,
				},
				csvData: []string{
					getCSVHeaderString(),
					"2024-01-01T00:00:00Z,xpto,localhost,linux,example-project,b1,E2x40,cmake,120.2,0,4,completed",
				},
			},
			wantErr: false,
		},
		{
			name:   "one valid JSONL entry with strict mode",
			r:      &mockReader{data: []byte(`{"project":"example-project","timestamp":"2024-01-01T00:00:00Z","user":"xpto","hostname":"localhost","os":"linux","branch":"b1","commit":"E2x40","command":"cmake","duration_sec":120.2,"return_code":0,"cpus":4,"status":"completed"}`)},
			strict: true,
			want: exportResult{
				ScanResult: metrics.ScanResult{
					Processed: 1,
					Skipped:   0,
				},
				csvData: []string{
					getCSVHeaderString(),
					"2024-01-01T00:00:00Z,xpto,localhost,linux,example-project,b1,E2x40,cmake,120.2,0,4,completed",
				},
			},
			wantErr: false,
		},
		{
			name:   "invalid JSONL entry should be skipped in non-strict mode",
			r:      &mockReader{data: []byte(`{"project": "example-project", "timestamp": "2024-01-01T00:00:00Z", "duration_sec": 120.2} INVALID_JSON`)},
			strict: false,
			want: exportResult{
				ScanResult: metrics.ScanResult{
					Processed: 0,
					Skipped:   1,
				},
				csvData: []string{getCSVHeaderString()},
			},
			wantErr: false,
		},
		{
			name:   "invalid JSONL entry should cause error in strict mode",
			r:      &mockReader{data: []byte(`{"project": "example-project", "timestamp": "2024-01-01T00:00:00Z", "duration_sec": 120.2} INVALID_JSON`)},
			strict: true,
			want: exportResult{
				ScanResult: metrics.ScanResult{
					Processed: 0,
					Skipped:   0, // In strict mode, we don't count skipped lines because we return an error immediately.
				},
				csvData: []string{getCSVHeaderString()},
			},
			wantErr: true,
		},
		{
			name: "Multiple valid JSONL entries should be processed correctly",
			r: &mockReader{data: []byte(`{"project": "A", "timestamp": "2024-01-01T10:00:00Z", "duration_sec": 5, "commit": "E2x40", "hostname": "localhost", "cpus": 4, "return_code": 0}
{"project": "A", "timestamp": "2024-01-02T12:00:00Z", "duration_sec": 5.2, "commit": "E2x40", "hostname": "localhost", "cpus": 4, "return_code": 0}
`)},
			strict: false,
			want: exportResult{
				ScanResult: metrics.ScanResult{
					Processed: 2,
					Skipped:   0,
				},
				csvData: []string{
					getCSVHeaderString(),
					"2024-01-01T10:00:00Z,,localhost,,A,,E2x40,,5,0,4,",
					"2024-01-02T12:00:00Z,,localhost,,A,,E2x40,,5.2,0,4,",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &mockWriter{}
			got, gotErr := metrics.ExportCSVFromJSONL(tt.r, writer, tt.strict)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExportCSVFromJSONL() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExportCSVFromJSONL() succeeded unexpectedly")
			}
			if got.Processed != tt.want.ScanResult.Processed || got.Skipped != tt.want.ScanResult.Skipped {
				t.Errorf("ExportCSVFromJSONL() = %v, want %v", got, tt.want)
			}
			csvLines := strings.Split(string(writer.data), "\n")
			// Remove linhas vazias do final
			if len(csvLines) > 0 && csvLines[len(csvLines)-1] == "" {
				csvLines = csvLines[:len(csvLines)-1]
			}
			if !reflect.DeepEqual(csvLines, tt.want.csvData) {
				t.Errorf("ExportCSVFromJSONL() CSV data = %v, want %v", csvLines, tt.want.csvData)
			}
		})
	}
}

func TestBuildMetricCSVRow(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		m    metrics.BuildMetric
		want []string
	}{
		{name: "Example test case", m: metrics.BuildMetric{
			Project:     "example-project",
			Timestamp:   "2024-01-01T00:00:00Z",
			User:        "xpto",
			Hostname:    "localhost",
			OS:          "linux",
			Branch:      "b1",
			Commit:      "E2x40",
			Command:     "cmake",
			DurationSec: 120.2,
			ReturnCode:  0,
			CPUs:        4,
			Status:      "completed",
		}, want: []string{"example-project", "2024-01-01T00:00:00Z", "xpto", "localhost", "linux", "b1", "E2x40", "cmake", "120.2", "0", "4", "completed"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := metrics.BuildMetricCSVRow(tt.m)
			if !equivalentStringSlices(t, got, tt.want) {
				t.Errorf("BuildMetricCSVRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCSVHeader(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want []string
	}{
		{name: "CSV header should have correct columns", want: []string{"timestamp", "user", "hostname", "os", "project", "branch", "commit", "command", "duration_sec", "returncode", "cpus", "status"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := metrics.CSVHeader()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CSVHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper para verificar se dois slices de strings tem os mesmos elementos.
func equivalentStringSlices(t *testing.T, a, b []string) bool {
	t.Helper()
	if len(a) != len(b) {
		return false
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i] < b[j]
	})
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func getCSVHeaderString() string {
	return "timestamp,user,hostname,os,project,branch,commit,command,duration_sec,returncode,cpus,status"
}
