package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"dev-metrics/internal/metrics"
)

type WeeklyStats struct {
	TotalDuration float64
	Count         int
}

func main() {
	// 1. Localizar o arquivo de log
	logFlag := flag.String("log", "", "Caminho do arquivo de log para analisar")
	flag.Parse()

	logPath, err := metrics.PrintResolvedLogPath(os.Stdout, "Usando arquivo de log: ", *logFlag)
	if err != nil {
		return
	}

	file, err := os.Open(logPath)
	if err != nil {
		fmt.Printf("Erro ao abrir log: %v\nCertifique-se de que já realizou algum build.\n", err)
		return
	}
	defer file.Close()

	// Mapa para agrupar: "Ano-Semana" -> Estatísticas
	report := make(map[string]*WeeklyStats)

	// 2. Ler linha por linha (JSONL)
	_, err = metrics.ScanJSONL(file, false, func(m metrics.BuildMetric) error {
		// Converter timestamp para o objeto Time do Go
		t, err := time.Parse(time.RFC3339, m.Timestamp)
		if err != nil {
			return nil
		}

		// Gerar chave da semana (ex: 2024-W32)
		year, week := t.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)

		if _, ok := report[weekKey]; !ok {
			report[weekKey] = &WeeklyStats{}
		}

		report[weekKey].TotalDuration += m.DurationSec
		report[weekKey].Count++
		return nil
	})
	if err != nil {
		fmt.Printf("Erro ao ler log: %v\n", err)
		return
	}

	// 3. Exibir Relatório
	fmt.Println("====================================================")
	fmt.Printf("%-12s | %-12s | %-12s | %-10s\n", "Semana", "Total (min)", "Média (s)", "Builds")
	fmt.Println("----------------------------------------------------")

	for week, stats := range report {
		totalMin := stats.TotalDuration / 60
		avgSec := stats.TotalDuration / float64(stats.Count)
		fmt.Printf("%-12s | %-12.2f | %-12.1f | %-10d\n", week, totalMin, avgSec, stats.Count)
	}
	fmt.Println("====================================================")
}
