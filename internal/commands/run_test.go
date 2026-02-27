package commands_test

import (
	"bytes"
	"context"
	"dev-metrics/internal/commands"
	"dev-metrics/internal/metrics"
	"errors"
	"os/user"
	"strings"
	"testing"
)

func TestExecCommand_Run(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		mockExitCode   int
		mockDuration   float64
		mockSaveErr    error
		wantErr        bool
		wantStatus     string
		wantStderr     string
		expectedCmdStr string
	}{
		{
			name:           "Success",
			args:           []string{"echo", "hello"},
			mockExitCode:   0,
			mockDuration:   1.5,
			wantStatus:     "success",
			expectedCmdStr: "[echo hello]",
		},
		{
			name:           "Command Failure",
			args:           []string{"false"},
			mockExitCode:   1,
			mockDuration:   0.5,
			wantStatus:     "failure",
			expectedCmdStr: "[false]",
		},
		{
			name:           "Metrics Save Error",
			args:           []string{"echo", "hello"},
			mockExitCode:   0,
			mockSaveErr:    errors.New("disk full"),
			wantStatus:     "success", // Command succeeded, metric save failed but Run returns nil
			wantStderr:     "[Metrics Error] disk full",
			expectedCmdStr: "[echo hello]",
		},
		{
			name:    "No Args",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var savedMetric metrics.BuildMetric
			var stdout, stderr bytes.Buffer

			cmd := &commands.ExecCommand{
				Out: &stdout,
				Err: &stderr,
				Runner: func(ctx context.Context, args []string) (float64, int) {
					return tt.mockDuration, tt.mockExitCode
				},
				GitInfo: func() (string, string, string) {
					return "main", "1234567", "test-project"
				},
				MetricsSaver: func(m metrics.BuildMetric, filePath string) error {
					savedMetric = m
					return tt.mockSaveErr
				},
				UserInfo: func() (*user.User, error) {
					return &user.User{Username: "testuser"}, nil
				},
				Hostname: func() (string, error) {
					return "testhost", nil
				},
			}

			err := cmd.Run(tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Verify Metric Fields
			if savedMetric.Status != tt.wantStatus {
				t.Errorf("Metric.Status = %v, want %v", savedMetric.Status, tt.wantStatus)
			}
			if savedMetric.DurationSec != tt.mockDuration {
				t.Errorf("Metric.DurationSec = %v, want %v", savedMetric.DurationSec, tt.mockDuration)
			}
			if savedMetric.ReturnCode != tt.mockExitCode {
				t.Errorf("Metric.ReturnCode = %v, want %v", savedMetric.ReturnCode, tt.mockExitCode)
			}
			if savedMetric.User != "testuser" {
				t.Errorf("Metric.User = %v, want %v", savedMetric.User, "testuser")
			}
			if savedMetric.Project != "test-project" {
				t.Errorf("Metric.Project = %v, want %v", savedMetric.Project, "test-project")
			}
			if tt.expectedCmdStr != "" && savedMetric.Command != tt.expectedCmdStr {
				t.Errorf("Metric.Command = %v, want %v", savedMetric.Command, tt.expectedCmdStr)
			}

			// Verify Stderr for logging errors
			if tt.wantStderr != "" {
				if !strings.Contains(stderr.String(), tt.wantStderr) {
					t.Errorf("Stderr = %q, want substring %q", stderr.String(), tt.wantStderr)
				}
			}
		})
	}
}

func TestExecCommand_Metadata(t *testing.T) {
	c := &commands.ExecCommand{}
	if c.Name() != "run" {
		t.Errorf("Name() = %q, want %q", c.Name(), "run")
	}
	if c.Description() == "" {
		t.Errorf("Description() is empty")
	}
}

func TestExecCommand_Aliases(t *testing.T) {
	c := &commands.ExecCommand{}
	aliases := c.Aliases()
	expected := []string{"exec", "r"}

	if len(aliases) != len(expected) {
		t.Errorf("Aliases() = %v, want %v", aliases, expected)
		return
	}

	for i, alias := range expected {
		if aliases[i] != alias {
			t.Errorf("Aliases()[%d] = %q, want %q", i, aliases[i], alias)
		}
	}
}

func TestExecCommand_ensureDefaults(t *testing.T) {
	var stdout, stderr bytes.Buffer // mock para não usar os.Stdout e os.Stderr reais
	c := &commands.ExecCommand{
		Out: &stdout,
		Err: &stderr,
	}
	c.Run([]string{}) // Call Run to trigger ensureDefaults
	if c.Out == nil {
		t.Error("Out is not set")
	}
	if c.Err == nil {
		t.Error("Err is not set")
	}
	if c.Runner == nil {
		t.Error("Runner is not set")
	}
	if c.GitInfo == nil {
		t.Error("GitInfo is not set")
	}
	if c.MetricsSaver == nil {
		t.Error("MetricsSaver is not set")
	}
	if c.UserInfo == nil {
		t.Error("UserInfo is not set")
	}
	if c.Hostname == nil {
		t.Error("Hostname is not set")
	}
}
