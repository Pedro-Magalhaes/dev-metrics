
# dev-metrics

Ferramenta em Go para **medir a duração de um comando** (ex.: build, testes, lint) e **registrar métricas** localmente em formato **JSON Lines** (`.jsonl`).

O fluxo é simples:

1) `measure-build` executa um comando e registra uma métrica em JSONL.
2) `analyze-metrics` analisa o arquivo JSONL e imprime um relatório semanal.
3) `export-metrics` converte o JSONL para CSV.

## Quickstart

Compilar os binários em `./dist`:

```bash
make build
```

Medir um comando (vai gerar/atualizar o log local):

```bash
./dist/build-meter go test ./...
```

Exportar para CSV:

```bash
./dist/export-meter -out /tmp/build_metrics.csv
```

## Instalação (Linux, sem sudo)

Este projeto instala **3 binários**: `build-meter`, `analyze-meter`, `export-meter`.

### Opção A: compilar e instalar localmente (recomendado para dev)

Instala em `~/.local/bin` (padrão):

```bash
make install
```

Garanta que `~/.local/bin` esteja no seu `PATH` (bash/zsh):

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Dica: para imprimir a linha acima:

```bash
make path-hint
```

### Opção B: instalar via GitHub Releases (sem Go instalado)

Instala a **última release** publicada do GitHub em `~/.local/bin`:

```bash
curl -fsSL https://raw.githubusercontent.com/Pedro-Magalhaes/dev-metrics/main/scripts/install.sh | sh
```

Para instalar uma versão específica:

```bash
VERSION=v1.2.3 curl -fsSL https://raw.githubusercontent.com/Pedro-Magalhaes/dev-metrics/main/scripts/install.sh | sh
```

Requisitos do instalador: `curl` (ou `wget`), `tar`, `sha256sum`.

## Comandos

### `measure-build` (binário: `build-meter`)

Executa o comando informado, mede a duração e salva uma linha JSON no log.

Exemplos:

```bash
./dist/build-meter go build ./...
./dist/build-meter go test ./...
./dist/build-meter make build
```

Flags úteis:

- `-log <path>`: sobrescreve o caminho do arquivo de log.

### `analyze-metrics` (binário: `analyze-meter`)

Analisa o log JSONL e imprime um relatório agrupado por semana (ISO 8601).

```bash
./dist/analyze-meter
./dist/analyze-meter -log /tmp/build_log.jsonl
```

### `export-metrics` (binário: `export-meter`)

Converte o log JSONL para CSV.

- Escreve **CSV válido em stdout** (ou em arquivo com `-out`).
- Escreve mensagens de status em `stderr`.

Executar via binário:

```bash
./dist/export-meter -out -
./dist/export-meter -out /tmp/build_metrics.csv
./dist/export-meter -log /tmp/build_log.jsonl -out /tmp/build_metrics.csv
./dist/export-meter -strict -out /tmp/build_metrics.csv
```

Executar via `go run`:

```bash
go run ./cmd/export-metrics -out -
go run ./cmd/export-metrics -out /tmp/build_metrics.csv
```

Ajuda (flags disponíveis):

```bash
./dist/export-meter -h
```

Observações:

- O CSV sempre inclui a primeira linha com o header (nomes das colunas).
- Por padrão, linhas inválidas no JSONL são ignoradas; use `-strict` para falhar ao primeiro erro.
- Para capturar somente o CSV (sem mensagens em `stderr`): `... 2>/dev/null`.

## Log: local e configuração

Por padrão as métricas são gravadas em:

- `~/.local/share/build-metrics/build_log.jsonl`

Cada execução adiciona **uma linha** (JSON) ao arquivo.

Você pode sobrescrever o caminho do arquivo de log nesta ordem:

1. Flag da CLI `-log <path>`
2. Variável de ambiente `BUILD_METRICS_LOG`
3. Padrão `~/.local/share/build-metrics/build_log.jsonl`

Exemplos:

```bash
./dist/build-meter -log /tmp/build_log.jsonl make build
export BUILD_METRICS_LOG=/tmp/build_log.jsonl
./dist/analyze-meter
```

## Formatos e schema

### JSONL

O arquivo de log é JSON Lines (`.jsonl`): uma linha JSON por execução.

### CSV

O export usa um header fixo e gera uma linha por métrica.

### Campos registrados (schema)

O objeto gravado segue o struct `metrics.BuildMetric`:

- `timestamp` (RFC3339)
- `user`
- `hostname`
- `os` (ex.: `linux`, `darwin`, `windows`)
- `project` (nome do projeto git ou `"unknown"`)
- `branch` (ou `"unknown"`)
- `commit` (hash curto, ou `"unknown"`)
- `command`
- `duration_sec`
- `returncode`
- `cpus`
- `status` (`success`/`failure`)

## Estrutura do projeto

- `cmd/measure-build/main.go`: coleta e grava métricas.
- `cmd/analyze-metrics/main.go`: relatório semanal.
- `cmd/export-metrics/main.go`: export JSONL → CSV.
- `internal/runner/exec.go`: executor do comando e medição de duração (`runner.Run`).
- `internal/git/info.go`: coleta branch/commit (`git.GetInfo`).
- `internal/metrics/*`: modelo, persistência, paths e utilitários.

## Desenvolvimento

Checagens comuns:

```bash
gofmt -w .
go test ./...
go vet ./...
```

Baixar/validar deps:

```bash
go mod tidy
```

Build reprodutível:

```bash
go build -trimpath ./cmd/measure-build
```

Cross-compile (exemplo):

```bash
GOOS=linux GOARCH=amd64 go build -o build-meter_linux_amd64 ./cmd/measure-build
```

## Links

- Effective Go: https://go.dev/doc/effective_go
- Go Code Review Comments: https://github.com/golang/go/wiki/CodeReviewComments
- Go Blog: https://go.dev/blog/
- `gofmt`: https://pkg.go.dev/cmd/gofmt
- `go vet`: https://pkg.go.dev/cmd/vet
- Go Modules Reference: https://go.dev/ref/mod

