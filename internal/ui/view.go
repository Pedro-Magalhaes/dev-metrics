package ui

import (
	"dev-metrics/internal/metrics" // Ajuste o import conforme seu module
	"fmt"
	"io"
)

func RenderReportTable(w io.Writer, report *metrics.FullReport, totalUnit metrics.DurationUnit) {
	avgHeader := "Média (auto)"
	totalHeader := "Total"
	if totalUnit != metrics.DurationAuto {
		totalHeader = fmt.Sprintf("Total (%s)", metrics.DurationUnitLabel(totalUnit))
	}

	// Mostrar intervalo do relatório se fornecido
	if !report.Since.IsZero() || !report.Until.IsZero() {
		sinceStr, untilStr := "-", "agora"
		if !report.Since.IsZero() {
			sinceStr = report.Since.Format("2006-01-02")
		}
		if !report.Until.IsZero() {
			untilStr = report.Until.Format("2006-01-02")
		}
		fmt.Fprintf(w, "Período: Desde %s  até %s\n", sinceStr, untilStr)
		fmt.Fprintln(w, "====================================================")
	}

	for _, proj := range report.Projects {
		fmt.Fprintf(w, "\n%-12s : %-12s\n", "Projeto", proj.Name)
		fmt.Fprintln(w, "====================================================")
		fmt.Fprintf(w, "%-12s | %-12s | %-12s | %-10s\n", "Semana", totalHeader, avgHeader, "Builds")
		fmt.Fprintln(w, "----------------------------------------------------")

		for _, week := range proj.Weeks {
			totalStr := metrics.FormatDuration(week.TotalDuration, totalUnit, totalUnit == metrics.DurationAuto)
			avgStr := metrics.FormatDuration(week.AvgDuration, metrics.DurationAuto, true)
			fmt.Fprintf(w, "%-12s | %-12s | %-12s | %-10d\n",
				week.WeekLabel, totalStr, avgStr, week.Count)
		}

		fmt.Fprintln(w, "----------------------------------------------------")
		totalStr := metrics.FormatDuration(proj.TotalDuration, totalUnit, totalUnit == metrics.DurationAuto)
		fmt.Fprintf(w, "%-12s | %-12s | %-12s | %-10d\n",
			"Total", totalStr, "-", proj.TotalBuilds)
		fmt.Fprintln(w, "====================================================")
	}

	// Resumo Global
	fmt.Fprintf(w, "\nRelatório Geral: \n")
	fmt.Fprintln(w, "====================================================")
	fmt.Fprintf(w, "%-12s | %-12s | %-12s | %-12s\n", "", totalHeader, avgHeader, "Builds")
	fmt.Fprintln(w, "----------------------------------------------------")
	globalTotalStr := metrics.FormatDuration(report.GlobalDuration, totalUnit, totalUnit == metrics.DurationAuto)
	fmt.Fprintf(w, "%-12s | %-12s | %-12s | %-10d\n",
		"", globalTotalStr, "-", report.GlobalBuilds)
	fmt.Fprintln(w, "====================================================")
}
