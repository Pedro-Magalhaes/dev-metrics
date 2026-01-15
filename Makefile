# Variáveis
BINARY_DIR := dist
CMD_DIR := ./cmd
INTERNAL_PKG := dev-metrics/internal/metrics

# Instalação (sem sudo)
PREFIX ?= $(HOME)/.local
INSTALLBINDIR ?= $(PREFIX)/bin
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
	@mkdir -p $(BINARY_DIR)

# Compila os dois binários
build: setup
	@echo "Compilando $(BMT)..."
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_DIR)/$(BMT) $(CMD_DIR)/$(BMT)
	
	@echo "Build concluído! Binários disponíveis em ./$(BINARY_DIR)"

# Remove a pasta dist
clean:
	@echo "Limpando arquivos..."
	rm -rf $(BINARY_DIR)
	@echo "Pasta $(BINARY_DIR) removida."

# Atalho para rodar o build-meter (exemplo)
# Use como: make run ARGS="cmake --version"
run: build
	./$(BINARY_DIR)/$(BMT) $(ARGS)

install: build
	@echo "Instalando binários em $(INSTALLBINDIR)"
	@mkdir -p "$(INSTALLBINDIR)"
	@$(INSTALL) -m 0755 "$(BINARY_DIR)/$(BMT)" "$(INSTALLBINDIR)/$(BMT)"
	@echo "OK: $(BMT), instalados em $(INSTALLBINDIR)"
	@$(MAKE) --no-print-directory path-hint

uninstall:
	@echo "Removendo binários de $(INSTALLBINDIR)"
	@rm -f "$(INSTALLBINDIR)/$(BMT)"
	@echo "OK: removido"

path-hint:
	@echo "Adicione ao seu shell rc (bash/zsh):"
	@echo "  export PATH=\"$(INSTALLBINDIR):\$$PATH\""