package commands

import (
	"context"
	"dev-metrics/internal/git"
	metrics "dev-metrics/internal/metrics"
	"dev-metrics/internal/runner"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"os/user"
	"runtime"
	"syscall"
	"time"
)

type ExecCommand struct {
	Out          io.Writer
	Err          io.Writer
	Runner       func(ctx context.Context, args []string) (float64, int)
	GitInfo      func() (string, string, string)
	MetricsSaver func(metrics.BuildMetric, string) error
	UserInfo     func() (*user.User, error)
	Hostname     func() (string, error)
}

func (c *ExecCommand) Name() string { return "run" }
func (c *ExecCommand) Description() string {
	return "Executa e mede um comando (ex: mtx run cmake --build .)"
}

func (c *ExecCommand) Run(args []string) error {
	// Defaults for dependencies
	c.ensureDefaults()

	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(c.Out)
	logFlag := fs.String("log", "", "Caminho customizado para o arquivo de log")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Uso: bmt run [-log path] <comando> [args...]\n")
		fs.PrintDefaults()
		// Note: PrintResolvedLogPath writes to fs.Output() internally if passed
		metrics.PrintResolvedLogPath(fs.Output(), "Arquivo de log: ", fs.Lookup("log").Value.String())
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	cmdArgs := fs.Args()

	if len(cmdArgs) < 1 {
		fs.Usage()
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

	duration, exitCode := c.Runner(ctx, cmdArgs)

	// 2. Coleta metadados
	currUser, _ := c.UserInfo()
	hostname, _ := c.Hostname()
	branch, commit, project := c.GitInfo()

	status := "success"
	if exitCode != 0 {
		status = "failure"
	}

	// Verifica se foi interrompido (contexto cancelado)
	if ctx.Err() == context.Canceled {
		status = "interrupted"
	}

	username := "unknown"
	if currUser != nil {
		username = currUser.Username
	}

	// 3. Monta a métrica
	metric := metrics.BuildMetric{
		Timestamp:   time.Now().Format(time.RFC3339),
		User:        username,
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
	if err := c.MetricsSaver(metric, logPath); err != nil {
		fmt.Fprintf(c.Err, "[Metrics Error] %v\n", err)
	}
	// DEBUG

	return nil
}

func (c *ExecCommand) Aliases() []string {
	return []string{"exec", "r"}
}

func (c *ExecCommand) ensureDefaults() {
	if c.Out == nil {
		c.Out = os.Stdout
	}
	if c.Err == nil {
		c.Err = os.Stderr
	}
	if c.Runner == nil {
		c.Runner = runner.Run
	}
	if c.GitInfo == nil {
		c.GitInfo = git.GetInfo
	}
	if c.MetricsSaver == nil {
		c.MetricsSaver = metrics.Save
	}
	if c.UserInfo == nil {
		c.UserInfo = user.Current
	}
	if c.Hostname == nil {
		c.Hostname = os.Hostname
	}
}

func init() {
	Register(&ExecCommand{})
}
