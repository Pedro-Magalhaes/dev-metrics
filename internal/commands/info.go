package commands

import (
	metrics "dev-metrics/internal/metrics"
	"fmt"
	"io"
	"os"
	"time"
)

type Info struct {
	Out io.Writer
}

func (c *Info) Name() string { return "info" }
func (c *Info) Description() string {
	return "Exibe informações sobre a ferramenta de métricas de build"
}

func (c *Info) Run(args []string) error {
	c.ensureDefaults()
	fmt.Fprintln(c.Out, "Build Metrics Tool")
	fmt.Fprintf(c.Out, "Version: %s\n", metrics.Version)
	fmt.Fprintf(c.Out, "Commit: %s\n", metrics.GitCommit)

	buildTime, err := time.Parse(time.RFC3339, metrics.BuildTime)
	if err != nil {
		fmt.Fprintf(c.Out, "Build Time: %s\n", metrics.BuildTime)
	} else {
		fmt.Fprintf(c.Out, "Build Time: %s\n", buildTime.Local().Format(time.RFC3339))
	}

	path, err := metrics.GetLogFilePath("")
	if err != nil {
		fmt.Fprintf(c.Out, "Erro ao resolver caminho do log: %v\n", err)
		return nil // ignora erro pois já foi reportado
	}

	fmt.Fprintf(c.Out, "Arquivo de log: %s\n", path)

	// Obtém mais informações do arquivo
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(c.Out, "Erro ao obter informações do arquivo: %v\n", err)
		return nil // ignora erro pois já foi reportado
	}

	sizeInMb := float64(fileInfo.Size()) / (1024 * 1024)
	fmt.Fprintf(c.Out, "Tamanho do arquivo: %.2f Mb\n", sizeInMb)
	fmt.Fprintf(c.Out, "Última modificação: %s\n", fileInfo.ModTime().Format(time.RFC1123))

	return nil
}

func (c *Info) Aliases() []string {
	return []string{"version", "v"}
}

func (c *Info) ensureDefaults() {
	if c.Out == nil {
		c.Out = os.Stdout
	}
}

func init() {
	Register(&Info{})
}
