package commands

import (
	"flag"
	"fmt"
	"os"
	"time"

	"dev-metrics/internal/metrics"
	"dev-metrics/internal/ui"
)

type ReportCommand struct{}

func (c *ReportCommand) Name() string { return "report" }
func (c *ReportCommand) Description() string {
	return "Gera um relatório a partir dos dados coletados de builds"
}

func (c *ReportCommand) Usage() string {
	return `report [--log PATH] [--since YYYY-MM-DD] [--until YYYY-MM-DD]

Opções:
  --log PATH            Caminho do arquivo de log (padrão: ./dev-metrics.log)
  --since YYYY-MM-DD    Data de início para filtrar o relatório
  --until YYYY-MM-DD    Data de fim para filtrar o relatório

Exemplo:
  dev-metrics report --log ./logs/dev-metrics.log --since 2024-01-01 --until 2024-01-31
`
}

func (c *ReportCommand) Run(args []string) error {
	// 1. Configuração e abertura de arquivo
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	logFlag := fs.String("log", "", "Caminho do arquivo de log")
	sinceFlag := fs.String("since", "", "Data de início (YYYY-MM-DD) para filtrar o relatório")
	untilFlag := fs.String("until", "", "Data de fim (YYYY-MM-DD) para filtrar o relatório")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	logPath, err := metrics.PrintResolvedLogPath(os.Stdout, "Usando arquivo de log: ", *logFlag)
	if err != nil {
		return fmt.Errorf("Erro ao obter path do arquivo de log: %v\n", err)
	}

	file, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("Erro ao abrir log (%s): %v", logPath, err)
	}
	defer file.Close()

	// Parse das opções
	opts := metrics.ReportOptions{}
	if *sinceFlag != "" {
		t, err := time.ParseInLocation("2006-01-02", *sinceFlag, time.Local)
		if err != nil {
			return fmt.Errorf("formato de data inválido para --since (use YYYY-MM-DD): %v", err)
		}
		opts.Since = t
	}

	if *untilFlag != "" {
		t, err := time.ParseInLocation("2006-01-02", *untilFlag, time.Local)
		if err != nil {
			return fmt.Errorf("formato de data inválido para --until (use YYYY-MM-DD): %v", err)
		}
		opts.Until = t
	}

	// Geração do relatório
	reportData, err := metrics.GenerateReport(file, opts)
	if err != nil {
		return fmt.Errorf("Erro ao processar dados: %v", err)
	}

	// Passamos os.Stdout para que ele escreva no terminal
	ui.RenderReportTable(os.Stdout, reportData)

	return nil
}

func init() {
	Register(&ReportCommand{})
}
