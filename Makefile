# Variáveis
REPO_BIN_DIR := dist
CMD_DIR := ./cmd
INTERNAL_PKG := dev-metrics/internal/metrics

# Instalação (sem sudo)
PREFIX ?= $(HOME)/.local
BIN_DIR ?= $(PREFIX)/bin
INSTALL ?= install

# Nomes dos binários
BMT := bmt

# Captura informações do ambiente
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# Monta os flags de injeção
# O formato é -X pacote.variavel=valor
LDFLAGS := -X $(INTERNAL_PKG).Version=$(VERSION) \
           -X $(INTERNAL_PKG).GitCommit=$(GIT_COMMIT) \
           -X $(INTERNAL_PKG).BuildTime=$(BUILD_TIME) \
           -s -w

.PHONY: all build clean setup install uninstall path-hint

# Alvo padrão: limpa, prepara e compila tudo
all: setup build

# Cria a pasta dist se não existir
setup:
	@mkdir -p $(REPO_BIN_DIR)

# Compila os dois binários
build: setup
	@echo "Compilando $(BMT)..."
	go build -ldflags="$(LDFLAGS)" -o $(REPO_BIN_DIR)/$(BMT) $(CMD_DIR)/$(BMT)
	
	@echo "Build concluído! Binários disponíveis em ./$(REPO_BIN_DIR)"

# Remove a pasta dist
clean:
	@echo "Limpando arquivos..."
	rm -rf $(REPO_BIN_DIR)
	@echo "Pasta $(REPO_BIN_DIR) removida."

# Atalho para rodar o build-meter (exemplo)
# Use como: make run ARGS="cmake --version"
run: build
	./$(REPO_BIN_DIR)/$(BMT) $(ARGS)

install: build
	@echo "Instalando binários em $(BIN_DIR)"
	@mkdir -p "$(BIN_DIR)"
	@$(INSTALL) -m 0755 "$(REPO_BIN_DIR)/$(BMT)" "$(BIN_DIR)/$(BMT)"
	@echo "OK: $(BMT), instalados em $(BIN_DIR)"
	@$(MAKE) --no-print-directory path-hint

uninstall:
	@echo "Removendo binários de $(BIN_DIR)"
	@rm -f "$(BIN_DIR)/$(BMT)"
	@echo "OK: removido"

path-hint:
	@echo "Adicione ao seu shell rc (bash/zsh):"
	@echo "  export PATH=\"$(BIN_DIR):\$$PATH\""