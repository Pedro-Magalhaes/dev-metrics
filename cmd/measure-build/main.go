package main

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"time"

	"dev-metrics/internal/git"
	"dev-metrics/internal/metrics"
	"dev-metrics/internal/runner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: measure-build <comando> [argumentos...]")
		os.Exit(1)
	}

	// 1. Executa o comando solicitado
	cmdArgs := os.Args[1:]
	duration, exitCode := runner.Run(cmdArgs)

	// 2. Coleta metadados
	currUser, _ := user.Current()
	hostname, _ := os.Hostname()
	branch, commit, project := git.GetInfo()

	status := "success"
	if exitCode != 0 {
		status = "failure"
	}

	// 3. Monta a métrica
	metric := metrics.BuildMetric{
		Timestamp:   time.Now().Format(time.RFC3339),
		User:        currUser.Username,
		Hostname:    hostname,
		OS:          runtime.GOOS,
		Project:     project,
		Branch:      branch,
		Commit:      commit,
		Command:     fmt.Sprintf("%v", cmdArgs),
		DurationSec: duration,
		ReturnCode:  exitCode,
		CPUs:        runtime.NumCPU(),
		Status:      status,
	}

	// 4. Salva (falha silenciosa para não atrapalhar o dev)
	if err := metrics.Save(metric); err != nil {
		fmt.Fprintf(os.Stderr, "[Metrics Error] %v\n", err)
	}

	// 5. Sai com o código original
	os.Exit(exitCode)
}
