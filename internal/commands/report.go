package commands

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"dev-metrics/internal/metrics"
	"dev-metrics/internal/ui"
)

type ReportCommand struct {
	FileOpener func(name string) (io.ReadCloser, error)
	Out        io.Writer
}

func (c *ReportCommand) Name() string { return "report" }
func (c *ReportCommand) Description() string {
	return "Gera um relatório a partir dos dados coletados de builds"
}

func (c *ReportCommand) Run(args []string) error {
	c.ensureDefaults()
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	logFlag := fs.String("log", "", "Caminho do arquivo de log")
	sinceFlag := fs.String("since", "", "Data de início (YYYY-MM-DD) para filtrar o relatório")
	untilFlag := fs.String("until", "", "Data de fim (YYYY-MM-DD) para filtrar o relatório")
	unitFlag := fs.String("unit", "auto", "Unidade para os totais (auto|s|min|h)")
	fs.SetOutput(c.Out)
	fs.Usage = func() {
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "%s", `Exemplo:
  bmt report --log ./logs/dev-metrics.log --since 2024-01-01 --until 2024-01-31`)
	}
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	logPath, err := metrics.PrintResolvedLogPath(c.Out, "Usando arquivo de log: ", *logFlag)
	if err != nil {
		return fmt.Errorf("Erro ao obter path do arquivo de log: %v\n", err)
	}

	file, err := c.FileOpener(logPath)
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

	unit, err := metrics.ParseDurationUnit(*unitFlag)
	if err != nil {
		return err
	}

	// Geração do relatório
	reportData, err := metrics.GenerateReport(file, opts)
	if err != nil {
		return fmt.Errorf("Erro ao processar dados: %v", err)
	}

	// Passamos os.Stdout para que ele escreva no terminal
	ui.RenderReportTable(c.Out, reportData, unit)

	return nil
}

func (c *ReportCommand) Aliases() []string {
	return []string{}
}

func (c *ReportCommand) ensureDefaults() {
	if c.FileOpener == nil {
		c.FileOpener = func(name string) (io.ReadCloser, error) {
			return os.Open(name)
		}
	}
	if c.Out == nil {
		c.Out = os.Stdout
	}
}

func init() {
	Register(&ReportCommand{})
}
