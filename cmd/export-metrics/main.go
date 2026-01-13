package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"dev-metrics/internal/metrics"
)

func main() {
	logOverride := flag.String("log", "", "Caminho do arquivo JSONL de log (ou use BUILD_METRICS_LOG)")
	outPath := flag.String("out", "-", "Caminho do arquivo CSV de saída (ou '-' para stdout)")
	strict := flag.Bool("strict", false, "Falha ao encontrar linhas inválidas no JSONL")
	flag.Parse()

	logPath, err := metrics.GetLogFilePath(*logOverride)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao resolver log: %v\n", err)
		os.Exit(2)
	}

	in, err := os.Open(logPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao abrir log %s: %v\n", logPath, err)
		os.Exit(1)
	}
	defer in.Close()

	var out *os.File
	if *outPath == "-" {
		out = os.Stdout
	} else {
		if err := metrics.EnsureLogDir(filepath.Dir(*outPath)); err != nil {
			fmt.Fprintf(os.Stderr, "erro ao criar diretório de saída: %v\n", err)
			os.Exit(1)
		}
		f, err := os.Create(*outPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro ao criar %s: %v\n", *outPath, err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	res, err := metrics.ExportCSVFromJSONL(in, out, *strict)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao exportar: %v\n", err)
		os.Exit(1)
	}

	if *outPath != "-" {
		fmt.Fprintf(os.Stderr, "exportado: %d linhas (puladas: %d) -> %s\n", res.Processed, res.Skipped, *outPath)
	} else {
		fmt.Fprintf(os.Stderr, "exportado: %d linhas (puladas: %d)\n", res.Processed, res.Skipped)
	}
}
