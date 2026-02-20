package metrics

import "time"

// BuildMetric representa os dados coletados de uma execução de build
type BuildMetric struct {
	Timestamp   string  `json:"timestamp"`
	User        string  `json:"user"`
	Hostname    string  `json:"hostname"`
	OS          string  `json:"os"`
	Project     string  `json:"project"`
	Branch      string  `json:"branch"`
	Commit      string  `json:"commit"`
	Command     string  `json:"command"`
	DurationSec float64 `json:"duration_sec"`
	ReturnCode  int     `json:"returncode"`
	CPUs        int     `json:"cpus"`
	Status      string  `json:"status"`
}

// BuildStats armazena estatísticas agregadas por semana
type BuildStats struct {
	TotalDuration float64
	Count         int
}

// WeeklySummary representa uma linha da tabela do report (uma semana)
type WeeklySummary struct {
	WeekLabel   string  // ex: "2024-W32"
	AvgDuration float64 // Segundos
	BuildStats
}

// ProjectSummary representa um bloco de tabela do report (um projeto)
type ProjectSummary struct {
	Name          string
	Weeks         []WeeklySummary // Ordenar por semana
	TotalDuration float64
	TotalBuilds   int
}

// ReportOptions define filtros e opções para geração do relatório
type ReportOptions struct {
	Since time.Time // Desde quando olhar os dados. Se zero, olha desde o início.
	Until time.Time // Até quando olhar os dados. Se zero, olha até o último dado.
}

// FullReport contém todos os dados prontos para exibição
type FullReport struct {
	Projects       []ProjectSummary // Ordenar alfabeticamente
	GlobalDuration float64
	GlobalBuilds   int
	ReportOptions
}
