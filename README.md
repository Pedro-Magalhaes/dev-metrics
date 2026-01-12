
# dev-metrics

Ferramenta em Go para **medir a duração de um comando** (ex.: build, testes, lint) e **registrar métricas** localmente em formato **JSON Lines** (`.jsonl`).

O binário principal é o `measure-build`.

## O que este projeto faz

Ao executar:

- roda o comando que você passou (preservando saída, cores e interatividade);
- mede o tempo total de execução (em segundos);
- coleta metadados do ambiente e do Git (branch e commit);
- grava uma linha JSON por execução em um arquivo local;
- retorna o **mesmo exit code** do comando original.

## Estrutura do projeto

- `cmd/measure-build/main.go`: CLI principal.
- `cmd/analyze-metrics/main.go`: utilitário para analisar o arquivo JSONL e gerar relatórios por semana.
- `internal/runner/exec.go`: executor do comando e medição de duração (`runner.Run`).
- `internal/git/info.go`: coleta branch/commit (`git.GetInfo`).
- `internal/metrics/models.go`: modelo da métrica (`metrics.BuildMetric`).
- `internal/metrics/writer.go`: persistência em JSONL (`metrics.Save`).
- `internal/metrics/models.go`: modelo da métrica (`metrics.BuildMetric`).
- `internal/metrics/version.go`: metadados de versão embutidos via ldflags (versão, commit, build time).

## Como usar

### 1) Build do binário

Pelo terminal na raiz do projeto:

```bash
go build -o measure-build ./cmd/measure-build
```

### 2) Medir um comando (exemplos)

Medir um build:

```bash
./measure-build go build ./...
```

Medir testes:

```bash
./measure-build go test ./...
```

Medir qualquer comando:

```bash
./measure-build make build
```

## Onde as métricas são salvas

Por padrão as métricas são gravadas em:

- `~/.local/share/build-metrics/build_log.jsonl`

Cada execução adiciona **uma linha** (JSON) ao arquivo.

Você pode sobrescrever o caminho do arquivo de log com duas opções (ordem de prioridade):

1. Flag da CLI `-log <path>` — exemplo: `./measure-build -log /tmp/build_log.jsonl make build`
2. Variável de ambiente `BUILD_METRICS_LOG` — exemplo: `export BUILD_METRICS_LOG=/tmp/build_log.jsonl`

Internamente o binário resolve o caminho na seguinte ordem: flag `-log` (se fornecida), depois a variável de ambiente `BUILD_METRICS_LOG`, e por fim o caminho padrão em `~/.local/share/...`.

## Campos registrados (schema)

O objeto gravado segue o struct `metrics.BuildMetric`:

- `timestamp` (RFC3339)
- `user`
- `hostname`
- `os` (ex.: `linux`, `darwin`, `windows`)
- `project` (nome do projeto git ou `"unknown"`)
- `branch` (ou `"unknown"`)
- `commit` (hash curto, ou `"unknown"`)
- `command` (representação do array de args)
- `duration_sec`
- `returncode`
- `cpus`
- `status` (`success`/`failure`)

## Observações importantes

- Falhas ao gravar métricas **não interrompem** o comando (apenas logam em `stderr`). Veja `metrics.Save`.
- O comando executado herda `stdin/stdout/stderr` para manter interatividade e cores. Veja `runner.Run`.
- Se o binário `git` não estiver disponível ou você não estiver em um repositório Git, branch/commit ficam como `"unknown"`. Veja `git.GetInfo`.

---

## Dicas rápidas de Go (build e desenvolvimento)

- **Baixar/validar deps**:

	```bash
	go mod tidy
	```

- **Build reprodutível (sem cache de paths)**:

	```bash
	go build -trimpath ./cmd/measure-build
	```

- **Build com `Makefile`**:

	Este repositório contém um `Makefile` que prepara e compila os binários em `./dist`. Os alvos principais são:

	- `make` — prepara e compila os binários (padrão `all`).
	- `make build` — compila os binários e os coloca em `./dist`.
	- `make clean` — remove a pasta `dist`.
	- `make run ARGS="..."` — compila e executa o `build-meter` com os `ARGS` fornecidos.

	Exemplo de uso:

	```bash
	# compila tudo
	make
	```

- **Instalar no `$GOBIN`/`$GOPATH/bin`**:

	```bash
	go install ./cmd/measure-build
	```

	Depois:

	```bash
	measure-build go test ./...
	```

- **Cross-compile**:

	```bash
	GOOS=linux GOARCH=amd64 go build -o measure-build_linux_amd64 ./cmd/measure-build
	```

- **Checagens comuns antes de commitar**:

	```bash
	gofmt -w .
	go test ./...
	go vet ./...
	```

---

## Links: melhores práticas (Go)

- Effective Go: https://go.dev/doc/effective_go
- Go Code Review Comments: https://github.com/golang/go/wiki/CodeReviewComments
- Standard Go Project Layout (referência comum): https://github.com/golang-standards/project-layout
- Go Blog (boas práticas e novidades): https://go.dev/blog/
- `gofmt`: https://pkg.go.dev/cmd/gofmt
- `go vet`: https://pkg.go.dev/cmd/vet
- Go Modules Reference: https://go.dev/ref/mod

