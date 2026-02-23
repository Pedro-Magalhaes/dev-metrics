package commands

import (
	metrics "dev-metrics/internal/metrics"
	"flag"
	"fmt"
	"io"
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
	return "Executa e mede um comando (ex: mtx export -out path/to/new.csv -log path/to/file.jsonl)"
}

func (c *ExportCommand) Run(args []string) error {
	c.ensureDefaults()

	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	fs.SetOutput(c.Out)

	logOverride := fs.String("log", "", "Caminho do arquivo JSONL de log (ou use BUILD_METRICS_LOG)")
	outPath := fs.String("out", "-", "Caminho do arquivo CSV de saída (ou '-' para stdout)")
	strict := fs.Bool("strict", false, "Falha ao encontrar linhas inválidas no JSONL")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Uso: export [-out path] [-log path] \n")
		fs.PrintDefaults()
		metrics.PrintResolvedLogPath(fs.Output(), "Arquivo de log: ", fs.Lookup("log").Value.String())
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	logPath, err := metrics.GetLogFilePath(*logOverride)
	if err != nil {
		return fmt.Errorf("erro ao resolver log: %v\n", err)
	}

	in, err := c.MetricsOpener(logPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir log %s: %v\n", logPath, err)
	}
	defer in.Close()

	var out io.Writer
	if *outPath == "-" {
		out = os.Stdout
	} else {
		if err := metrics.EnsureLogDir(filepath.Dir(*outPath)); err != nil {
			return fmt.Errorf("erro ao criar diretório de saída: %v\n", err)
		}
		f, err := c.FileCreator(*outPath)
		if err != nil {
			return fmt.Errorf("erro ao criar %s: %v\n", *outPath, err)
		}
		defer f.Close()
		out = f
	}

	res, err := c.MetricsSaver(in, out, *strict)
	if err != nil {
		return fmt.Errorf("erro ao exportar: %v\n", err)
	}

	if *outPath != "-" {
		fmt.Fprintf(c.Err, "exportado: %d linhas (puladas: %d) -> %s\n", res.Processed, res.Skipped, *outPath)
	} else {
		fmt.Fprintf(c.Err, "exportado: %d linhas (puladas: %d)\n", res.Processed, res.Skipped)
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

func init() {
	Register(&ExportCommand{})
}
