package commands

import (
	metrics "dev-metrics/internal/metrics"
	"fmt"
	"os"
	"time"
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

	buildTime, err := time.Parse(time.RFC3339, metrics.BuildTime)
	if err != nil {
		fmt.Printf("Build Time: %s\n", metrics.BuildTime)
	} else {
		fmt.Printf("Build Time: %s\n", buildTime.Local().Format(time.RFC3339))
	}
	metrics.PrintResolvedLogPath(os.Stdout, "Arquivo de log: ", "")
	return nil
}

func init() {
	Register(&Info{})
}
