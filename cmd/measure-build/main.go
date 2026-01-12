package main

import (
	"flag"
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
	// Define a flag -log
	logFlag := flag.String("log", "", "Caminho customizado para o arquivo de log")
	flag.Parse()

	// Os argumentos após as flags são o comando real
	cmdArgs := flag.Args()

	if len(cmdArgs) < 1 {
		fmt.Println("Uso: measure-build [-log path] <comando> [args...]")
		os.Exit(1)
	}

	// Resolve o caminho do log
	logPath, err := metrics.GetLogFilePath(*logFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao resolver caminho do log: %v\n", err)
		os.Exit(1)
	}

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
	if err := metrics.Save(metric, logPath); err != nil {
		fmt.Fprintf(os.Stderr, "[Metrics Error] %v\n", err)
	}

	// 5. Sai com o código original
	os.Exit(exitCode)
}
