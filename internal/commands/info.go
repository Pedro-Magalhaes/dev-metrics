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
	out := c.Out
	if out == nil {
		out = os.Stdout
	}
	fmt.Fprintln(out, "Build Metrics Tool")
	fmt.Fprintf(out, "Version: %s\n", metrics.Version)
	fmt.Fprintf(out, "Commit: %s\n", metrics.GitCommit)

	buildTime, err := time.Parse(time.RFC3339, metrics.BuildTime)
	if err != nil {
		fmt.Fprintf(out, "Build Time: %s\n", metrics.BuildTime)
	} else {
		fmt.Fprintf(out, "Build Time: %s\n", buildTime.Local().Format(time.RFC3339))
	}
	metrics.PrintResolvedLogPath(out, "Arquivo de log: ", "")
	return nil
}

func init() {
	Register(&Info{})
}
