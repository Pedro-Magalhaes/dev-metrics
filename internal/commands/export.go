package commands

import (
	metrics "dev-metrics/internal/metrics"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type ExportCommand struct{}

func (c *ExportCommand) Name() string { return "export" }
func (c *ExportCommand) Description() string {
	return "Executa e mede um comando (ex: mtx export -out path/to/new.csv -log path/to/file.jsonl)"
}

func (c *ExportCommand) Run(args []string) error {
	fs := flag.NewFlagSet("export", flag.ExitOnError)

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

	in, err := os.Open(logPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir log %s: %v\n", logPath, err)
	}
	defer in.Close()

	var out *os.File
	if *outPath == "-" {
		out = os.Stdout
	} else {
		if err := metrics.EnsureLogDir(filepath.Dir(*outPath)); err != nil {
			return fmt.Errorf("erro ao criar diretório de saída: %v\n", err)
		}
		f, err := os.Create(*outPath)
		if err != nil {
			return fmt.Errorf("erro ao criar %s: %v\n", *outPath, err)
		}
		defer f.Close()
		out = f
	}

	res, err := metrics.ExportCSVFromJSONL(in, out, *strict)
	if err != nil {
		return fmt.Errorf("erro ao exportar: %v\n", err)
	}

	if *outPath != "-" {
		fmt.Fprintf(os.Stderr, "exportado: %d linhas (puladas: %d) -> %s\n", res.Processed, res.Skipped, *outPath)
	} else {
		fmt.Fprintf(os.Stderr, "exportado: %d linhas (puladas: %d)\n", res.Processed, res.Skipped)
	}
	// DEBUG

	return nil
}

func init() {
	Register(&ExportCommand{})
}
