package commands

import (
	metrics "dev-metrics/internal/metrics"
	"fmt"
	"os"
)

type Info struct{}

func (c *Info) Name() string { return "info" }
func (c *Info) Description() string {
	return "Exibe informações sobre a ferramenta de métricas de build"
}

func (c *Info) Run(args []string) error {
	fmt.Println("Build Metrics Tool")
	fmt.Printf("Version %s\n", metrics.Version)
	fmt.Printf("Commit: %s\n", metrics.GitCommit)
	fmt.Printf("Build Time: %s\n", metrics.BuildTime)
	metrics.PrintResolvedLogPath(os.Stdout, "Arquivo de log: ", "")
	return nil
}

func init() {
	Register(&Info{})
}
