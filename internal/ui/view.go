package ui

import (
	"dev-metrics/internal/metrics" // Ajuste o import conforme seu module
	"fmt"
	"io"
)

func RenderReportTable(w io.Writer, report *metrics.FullReport) {
	for _, proj := range report.Projects {
		fmt.Fprintf(w, "\n%-12s : %-12s\n", "Projeto", proj.Name)
		fmt.Fprintln(w, "====================================================")
		fmt.Fprintf(w, "%-12s | %-12s | %-12s | %-10s\n", "Semana", "Total (min)", "Média (s)", "Builds")
		fmt.Fprintln(w, "----------------------------------------------------")

		for _, week := range proj.Weeks {
			totalMin := week.TotalDuration / 60
			fmt.Fprintf(w, "%-12s | %-12.2f | %-12.1f | %-10d\n",
				week.WeekLabel, totalMin, week.AvgDuration, week.Count)
		}

		fmt.Fprintln(w, "----------------------------------------------------")
		fmt.Fprintf(w, "%-12s | %-12.2f | %-12s | %-10d\n",
			"Total", proj.TotalDuration/60, "-", proj.TotalBuilds)
		fmt.Fprintln(w, "====================================================")
	}

	// Resumo Global
	fmt.Fprintf(w, "\nRelatório Geral: \n")
	fmt.Fprintln(w, "====================================================")
	fmt.Fprintf(w, "%-12s | %-12s | %-12s | %-12s\n", "", "Total (min)", "Média", "Builds")
	fmt.Fprintln(w, "----------------------------------------------------")
	fmt.Fprintf(w, "%-12s | %-12.1f | %-12s | %-10d\n",
		"", report.GlobalDuration/60, "-", report.GlobalBuilds)
	fmt.Fprintln(w, "====================================================")
}
