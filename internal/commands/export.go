package commands

import (
	metrics "dev-metrics/internal/metrics"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type ExportCommand struct {
	Out           io.Writer
	Err           io.Writer
	MetricsOpener func(string) (io.ReadCloser, error)
	MetricsSaver  func(io.Reader, io.Writer, bool) (metrics.ScanResult, error)
	FileCreator   func(string) (io.WriteCloser, error)
}

func (c *ExportCommand) Name() string { return "export" }
func (c *ExportCommand) Description() string {
	return "Exporta as métricas salvas em formato CSV"
}

func (c *ExportCommand) Run(args []string) error {
	c.ensureDefaults()
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	fs.SetOutput(c.Out)

	logOverride := fs.String("log", "", "Caminho do arquivo JSONL de log (ou use BUILD_METRICS_LOG)")
	outPath := fs.String("out", "-", "Caminho do arquivo CSV de saída (ou '-' para stdout)")
	strict := fs.Bool("strict", false, "Falha ao encontrar linhas inválidas no JSONL")

	fs.Usage = func() {
		//nolint:errcheck
		fmt.Fprintf(fs.Output(), "Uso: export [-out path] [-log path] \n")
		fs.PrintDefaults()
		//nolint:errcheck
		metrics.PrintResolvedLogPath(fs.Output(), "Arquivo de log: ", fs.Lookup("log").Value.String())
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	logPath, err := metrics.GetLogFilePath(*logOverride)
	if err != nil {
		return fmt.Errorf("erro ao resolver log: %v", err)
	}

	in, err := c.MetricsOpener(logPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir log %s: %v", logPath, err)
	}
	defer Close(in)

	var out io.Writer
	if *outPath == "-" {
		out = c.Out
	} else {
		if err := metrics.EnsureLogDir(filepath.Dir(*outPath)); err != nil {
			return fmt.Errorf("erro ao criar diretório de saída: %v", err)
		}
		f, err := c.FileCreator(*outPath)
		if err != nil {
			return fmt.Errorf("erro ao criar %s: %v", *outPath, err)
		}
		defer Close(f)
		out = f
	}

	res, err := c.MetricsSaver(in, out, *strict)
	if err != nil {
		return fmt.Errorf("erro ao exportar: %v", err)
	}

	if *outPath != "-" {
		// nolint:errcheck
		fmt.Fprintf(c.Err, "exportado: %d linhas (puladas: %d) -> %s", res.Processed, res.Skipped, *outPath)
	} else {
		// nolint:errcheck
		fmt.Fprintf(c.Err, "exportado: %d linhas (puladas: %d)", res.Processed, res.Skipped)
	}

	return nil
}

func (c *ExportCommand) ensureDefaults() {
	if c.Out == nil {
		c.Out = os.Stdout
	}
	if c.Err == nil {
		c.Err = os.Stderr
	}
	if c.MetricsOpener == nil {
		c.MetricsOpener = func(name string) (io.ReadCloser, error) {
			return os.Open(name)
		}
	}
	if c.MetricsSaver == nil {
		// Wrapper function to match signature if necessary, or direct assignment
		c.MetricsSaver = func(in io.Reader, out io.Writer, strict bool) (metrics.ScanResult, error) {
			return metrics.ExportCSVFromJSONL(in, out, strict)
		}
	}
	if c.FileCreator == nil {
		c.FileCreator = func(name string) (io.WriteCloser, error) {
			return os.Create(name)
		}
	}
}

func (c *ExportCommand) Aliases() []string {
	return []string{}
}

func Close(f io.Closer) {
	if err := f.Close(); err != nil {
		log.Default().Println(err)
	}
}

func init() {
	Register(&ExportCommand{})
}
