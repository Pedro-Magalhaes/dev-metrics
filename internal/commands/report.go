package commands

import (
	"flag"
	"fmt"
	"os"

	"dev-metrics/internal/metrics"
	"dev-metrics/internal/ui"
)

type ReportCommand struct{}

func (c *ReportCommand) Name() string { return "report" }
func (c *ReportCommand) Description() string {
	return "Gera um relatório a partir dos dados coletados de builds"
}

func (c *ReportCommand) Run(args []string) error {
	// 1. Configuração e abertura de arquivo
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	logFlag := fs.String("log", "", "Caminho do arquivo de log")
	fs.Parse(args)

	logPath, err := metrics.PrintResolvedLogPath(os.Stdout, "Usando arquivo de log: ", *logFlag)
	if err != nil {
		return fmt.Errorf("Erro ao obter path do arquivo de log: %v\n", err)
	}

	file, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("Erro ao abrir log (%s): %v", logPath, err)
	}
	defer file.Close()

	reportData, err := metrics.GenerateReport(file)
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
