package commands_test

import (
	"bytes"
	"dev-metrics/internal/commands"
	"dev-metrics/internal/metrics"
	"strings"
	"testing"
	"time"
)

func TestInfo_Run(t *testing.T) {
	// Backup original values
	origVersion := metrics.Version
	origGitCommit := metrics.GitCommit
	origBuildTime := metrics.BuildTime

	defer func() {
		metrics.Version = origVersion
		metrics.GitCommit = origGitCommit
		metrics.BuildTime = origBuildTime
	}()

	metrics.Version = "1.0.0"
	metrics.GitCommit = "abcdef123"

	tests := []struct {
		name          string
		buildTime     string
		expectTimeStr string // Expected string in output for time
	}{
		{
			name:      "Valid Build Time",
			buildTime: "2023-10-27T10:00:00Z",
		},
		{
			name:          "Invalid Build Time",
			buildTime:     "invalid-time",
			expectTimeStr: "Build Time: invalid-time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics.BuildTime = tt.buildTime

			var buf bytes.Buffer
			c := &commands.Info{Out: &buf}

			err := c.Run(nil)
			if err != nil {
				t.Fatalf("Run() unexpected error: %v", err)
			}

			got := buf.String()

			// Check common fields
			if !strings.Contains(got, "Build Metrics Tool") {
				t.Errorf("Output missing title")
			}
			if !strings.Contains(got, "Version: 1.0.0") {
				t.Errorf("Output missing version")
			}
			if !strings.Contains(got, "Commit: abcdef123") {
				t.Errorf("Output missing commit")
			}
			if !strings.Contains(got, "Arquivo de log: ") {
				t.Errorf("Output missing log file path")
			}

			// Check time
			if tt.expectTimeStr != "" {
				if !strings.Contains(got, tt.expectTimeStr) {
					t.Errorf("Output missing expected time string: %q", tt.expectTimeStr)
				}
			} else {
				// For valid time, it converts to local.
				parsed, _ := time.Parse(time.RFC3339, tt.buildTime)
				expected := parsed.Local().Format(time.RFC3339)
				if !strings.Contains(got, "Build Time: "+expected) {
					t.Errorf("Output missing formatted time: %q. Got: %q", expected, got)
				}
			}
		})
	}
}

func TestInfo_Metadata(t *testing.T) {
	c := &commands.Info{}
	if c.Name() != "info" {
		t.Errorf("Name() = %q, want %q", c.Name(), "info")
	}
	if c.Description() == "" {
		t.Errorf("Description() is empty")
	}
}
