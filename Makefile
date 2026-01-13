# Variáveis
BINARY_DIR := dist
CMD_DIR := ./cmd
INTERNAL_PKG := dev-metrics/internal/metrics

# Nomes dos binários
BUILD_METER := build-meter
ANALYZE_METER := analyze-meter
EXPORT_METER := export-meter

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

.PHONY: all build clean setup

# Alvo padrão: limpa, prepara e compila tudo
all: setup build

# Cria a pasta dist se não existir
setup:
	@mkdir -p $(BINARY_DIR)

# Compila os dois binários
build: setup
	@echo "Compilando $(BUILD_METER)..."
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_DIR)/$(BUILD_METER) $(CMD_DIR)/measure-build
	
	@echo "Compilando $(ANALYZE_METER)..."
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_DIR)/$(ANALYZE_METER) $(CMD_DIR)/analyze-metrics
	
	@echo "Compilando $(EXPORT_METER)..."
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_DIR)/$(EXPORT_METER) $(CMD_DIR)/export-metrics
	@echo "Build concluído! Binários disponíveis em ./$(BINARY_DIR)"

# Remove a pasta dist
clean:
	@echo "Limpando arquivos..."
	rm -rf $(BINARY_DIR)
	@echo "Pasta $(BINARY_DIR) removida."

# Atalho para rodar o build-meter (exemplo)
# Use como: make run ARGS="cmake --version"
run: build
	./$(BINARY_DIR)/$(BUILD_METER) $(ARGS)