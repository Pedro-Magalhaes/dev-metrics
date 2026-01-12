
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
- `internal/runner/exec.go`: executor do comando e medição de duração (`runner.Run`).
- `internal/git/info.go`: coleta branch/commit (`git.GetInfo`).
- `internal/metrics/models.go`: modelo da métrica (`metrics.BuildMetric`).
- `internal/metrics/writer.go`: persistência em JSONL (`metrics.Save`).

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

As métricas são gravadas em:

- `~/.local/share/build-metrics/build_log.jsonl`

Cada execução adiciona **uma linha** (JSON) ao arquivo.

## Campos registrados (schema)

O objeto gravado segue o struct `metrics.BuildMetric`:

- `timestamp` (RFC3339)
- `user`
- `hostname`
- `os` (ex.: `linux`, `darwin`, `windows`)
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

