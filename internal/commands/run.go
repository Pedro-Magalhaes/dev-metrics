package commands

import (
	"context"
	"dev-metrics/internal/git"
	metrics "dev-metrics/internal/metrics"
	"dev-metrics/internal/runner"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"runtime"
	"syscall"
	"time"
)

type ExecCommand struct{}

func (c *ExecCommand) Name() string { return "run" }
func (c *ExecCommand) Description() string {
	return "Executa e mede um comando (ex: mtx run cmake --build .)"
}

func (c *ExecCommand) Run(args []string) error {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	logFlag := fs.String("log", "", "Caminho customizado para o arquivo de log")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Uso: bmt run [-log path] <comando> [args...]\n")
		fs.PrintDefaults()
		metrics.PrintResolvedLogPath(fs.Output(), "Arquivo de log: ", fs.Lookup("log").Value.String())
	}
	fs.Parse(args)

	cmdArgs := fs.Args()

	if len(cmdArgs) < 1 {
		return errors.New("nenhum comando fornecido para execução")
	}

	// Resolve o caminho do log ("erro ao resolver caminho do log: %v", err))
	logPath, err := metrics.GetLogFilePath(*logFlag)
	if err != nil {
		return errors.New(fmt.Sprint("erro ao resolver caminho do log:", err))
	}

	// Cria um context que pode ser cancelado
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Captura SIGINT (Ctrl+C) e cancela o contexto
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		<-sigChan
		cancel()
	}()

	duration, exitCode := runner.Run(ctx, cmdArgs)

	// 2. Coleta metadados
	currUser, _ := user.Current()
	hostname, _ := os.Hostname()
	branch, commit, project := git.GetInfo()

	status := "success"
	if exitCode != 0 {
		status = "failure"
	}

	// Verifica se foi interrompido (contexto cancelado)
	if ctx.Err() == context.Canceled {
		status = "interrupted"
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
	// DEBUG

	return nil
}

func init() {
	Register(&ExecCommand{})
}
