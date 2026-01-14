package commands

import (
	"flag"
	"fmt"
	"os"
	"time"

	"dev-metrics/internal/metrics"
)

type ReportCommand struct{}

func (c *ReportCommand) Name() string { return "report" }
func (c *ReportCommand) Description() string {
	return "Gera um relatório a partir dos dados coletados de builds"
}

func (c *ReportCommand) Run(args []string) error {
	// 1. Localizar o arquivo de log
	logFlag := flag.String("log", "", "Caminho do arquivo de log para analisar")
	flag.Parse()

	logPath, err := metrics.PrintResolvedLogPath(os.Stdout, "Usando arquivo de log: ", *logFlag)
	if err != nil {
		return fmt.Errorf("Erro ao obter path do arquivo de log: %v\n", err)
	}

	file, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("Erro ao abrir log: %v\nCertifique-se de que já realizou algum build.\n", err)

	}
	defer file.Close()

	// Mapa para agrupar: "Ano-Semana" -> Estatísticas por projeto
	reportByProject := make(map[string]map[string]*metrics.WeeklyStats)

	// 2. Ler linha por linha (JSONL)
	_, err = metrics.ScanJSONL(file, false, func(m metrics.BuildMetric) error {
		// Converter timestamp para o objeto Time do Go
		t, err := time.Parse(time.RFC3339, m.Timestamp)
		if err != nil {
			return fmt.Errorf("erro ao analisar timestamp: %v", err)
		}
		report := reportByProject[m.Project]
		if report == nil {
			report = make(map[string]*metrics.WeeklyStats)
			reportByProject[m.Project] = report
		}

		// Gerar chave da semana (ex: 2024-W32)
		year, week := t.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)

		if _, ok := report[weekKey]; !ok {
			report[weekKey] = &metrics.WeeklyStats{}
		}

		report[weekKey].TotalDuration += m.DurationSec
		report[weekKey].Count++
		return nil
	})
	if err != nil {
		return fmt.Errorf("Erro ao ler log: %v\n", err)

	}

	// 3. Exibir Relatório
	for projeto, reportFromProject := range reportByProject {
		fmt.Printf("\n%-12s : %-12s\n", "Projeto", projeto)
		fmt.Println("====================================================")
		fmt.Printf("%-12s | %-12s | %-12s | %-10s\n", "Semana", "Total (min)", "Média (s)", "Builds")
		fmt.Println("----------------------------------------------------")

		for week, stats := range reportFromProject {
			totalMin := stats.TotalDuration / 60
			avgSec := stats.TotalDuration / float64(stats.Count)
			fmt.Printf("%-12s | %-12.2f | %-12.1f | %-10d\n", week, totalMin, avgSec, stats.Count)
		}
		fmt.Println("====================================================")
	}
	return nil
}

func init() {
	Register(&ReportCommand{})
}
