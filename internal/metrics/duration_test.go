package metrics

import "testing"

func TestParseDurationUnit(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    DurationUnit
		wantErr bool
	}{
		{"auto_default", "", DurationAuto, false},
		{"auto_explicit", "auto", DurationAuto, false},
		{"seconds_short", "s", DurationSeconds, false},
		{"seconds_long", "seconds", DurationSeconds, false},
		{"minutes_short", "min", DurationMinutes, false},
		{"minutes_long", "minutes", DurationMinutes, false},
		{"hours_short", "h", DurationHours, false},
		{"hours_long", "hours", DurationHours, false},
		{"invalid", "banana", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDurationUnit(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("invalid result: got=%q want=%q", got, tt.want)
			}
		})
	}
}

func TestAutoDurationUnit(t *testing.T) {
	tests := []struct {
		name    string
		seconds float64
		want    DurationUnit
	}{
		{"below_minute", 59, DurationSeconds},
		{"exact_minute", 60, DurationMinutes},
		{"below_hour", 3599, DurationMinutes},
		{"exact_hour", 3600, DurationHours},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AutoDurationUnit(tt.seconds)
			if got != tt.want {
				t.Fatalf("invalid result: got=%q want=%q", got, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name        string
		seconds     float64
		unit        DurationUnit
		includeUnit bool
		want        string
	}{
		{"auto_seconds_with_unit", 45, DurationAuto, true, "45.0 s"},
		{"auto_minutes_with_unit", 90, DurationAuto, true, "1min30s"},
		{"auto_hours_with_unit", 5400, DurationAuto, true, "1h30min"},
		{"minutes_no_unit", 179, DurationMinutes, false, "2min59s"},
		{"seconds_with_unit", 120, DurationSeconds, true, "120.0 s"},
		{"hours_no_unit", 7200, DurationHours, false, "2h00min"},
		{"auto_no_unit", 95, DurationAuto, false, "1min35s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDuration(tt.seconds, tt.unit, tt.includeUnit)
			if got != tt.want {
				t.Fatalf("invalid result: got=%q want=%q", got, tt.want)
			}
		})
	}
}
