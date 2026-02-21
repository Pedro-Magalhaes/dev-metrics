package ui_test

import (
	"bytes"
	"dev-metrics/internal/metrics"
	"dev-metrics/internal/ui"
	"testing"
	"time"
)

func makeReportBasic() *metrics.FullReport {
	return &metrics.FullReport{
		Projects: []metrics.ProjectSummary{
			{
				Name:          "ProjetoA",
				Weeks:         []metrics.WeeklySummary{{WeekLabel: "2026-01", BuildStats: metrics.BuildStats{TotalDuration: 100, Count: 2}, AvgDuration: 50}},
				TotalDuration: 100,
				TotalBuilds:   2,
			},
		},
		GlobalDuration: 100,
		GlobalBuilds:   2,
		ReportOptions: metrics.ReportOptions{
			Since: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			Until: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}
}

func makeReportNoProjects() *metrics.FullReport {
	return &metrics.FullReport{
		Projects:       nil,
		GlobalDuration: 0,
		GlobalBuilds:   0,
	}
}

func makeReportWithoutOptions() *metrics.FullReport {
	return &metrics.FullReport{
		Projects: []metrics.ProjectSummary{
			{
				Name:          "ProjetoB",
				Weeks:         []metrics.WeeklySummary{{WeekLabel: "2026-02", BuildStats: metrics.BuildStats{TotalDuration: 200, Count: 2}, AvgDuration: 100}},
				TotalDuration: 200,
				TotalBuilds:   2,
			},
		},
		GlobalDuration: 200,
		GlobalBuilds:   2,
	}
}

func removeSpaces(s string) string {
	var buf bytes.Buffer
	for _, r := range s {
		if r != ' ' {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func containsIgnoreSpaces(haystack, needle string) bool {
	return bytes.Contains([]byte(removeSpaces(haystack)), []byte(removeSpaces(needle)))
}

func containsAllSnips(haystack string, snips []string) (bool, string) {
	for _, snip := range snips {
		if !containsIgnoreSpaces(haystack, snip) {
			return false, snip
		}
	}
	return true, ""
}

func TestRenderReportTable(t *testing.T) {

	tests := []struct {
		name      string
		report    *metrics.FullReport
		totalUnit metrics.DurationUnit
		wantSnips []string // trechos esperados na saída
	}{
		{
			name:      "Relatório básico com filtro since e until",
			report:    makeReportBasic(),
			totalUnit: metrics.DurationSeconds,
			wantSnips: []string{
				"Período: Desde 2026-01-01  até 2026-02-01",
				"Projeto : ProjetoA",
				"2026-01 | 100.0 | 50.0 s | 2 ",
				"Total | 100.0 | - | 2 ",
				"Relatório Geral:",
				"| 100.0 | - | 2 ",
			},
		},
		{
			name:      "Sem projetos",
			report:    makeReportNoProjects(),
			totalUnit: metrics.DurationSeconds,
			wantSnips: []string{
				"Relatório Geral:",
				"| 0.0           | -            | 0         ",
			},
		},
		{
			name:      "Sem options",
			report:    makeReportWithoutOptions(),
			totalUnit: metrics.DurationSeconds,
			wantSnips: []string{
				"Projeto : ProjetoB",
				"2026-02 | 200.0 | 1min40s | 2 ",
				"Total | 200.0 | - | 2 ",
				"Relatório Geral:",
				"| 200.0 | - | 2 ",
			},
		},
		{
			name:      "Sem options em minutos",
			report:    makeReportWithoutOptions(),
			totalUnit: metrics.DurationMinutes,
			wantSnips: []string{
				"Projeto : ProjetoB",
				"Semana | Total (min) | Média (auto) | Builds",
				"2026-02 | 3min20s | 1min40s | 2 ",
				"Total | 3min20s | - | 2 ",
				"Relatório Geral:",
				"| 3min20s | -| 2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			ui.RenderReportTable(&buf, tt.report, tt.totalUnit)
			got := buf.String()
			all, missing := containsAllSnips(got, tt.wantSnips)
			if !all {
				t.Errorf("RenderReportTable() output does not contain all expected snippets.\nGot:\n%s\nExpected snippets:\n%v\nMissing snippet: %s", got, tt.wantSnips, missing)
			}
		})
	}
}
