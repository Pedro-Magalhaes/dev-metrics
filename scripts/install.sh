#!/usr/bin/env sh
set -eu

REPO_DEFAULT="Pedro-Magalhaes/dev-metrics"
REPO="${REPO:-$REPO_DEFAULT}"
BINDIR="${BINDIR:-${HOME}/.local/bin}"
VERSION="${VERSION:-latest}"
BIN_NAME="bmt"

fail() {
  printf '❌ %s\n' "$*" >&2
  exit 1
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "dependência ausente: $1"
}

# Verificações de dependências
need_cmd uname
need_cmd mkdir
need_cmd rm
need_cmd mktemp
need_cmd tar

if command -v curl >/dev/null 2>&1; then
  DL="curl"
elif command -v wget >/dev/null 2>&1; then
  DL="wget"
else
  fail "é necessário curl ou wget"
fi

# Detectar OS e Arquitetura
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH_RAW="$(uname -m)"

case "$OS" in
  linux) : ;;
  *) fail "OS não suportado: $OS (apenas linux é suportado atualmente)" ;;
esac

case "$ARCH_RAW" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) fail "arquitetura não suportada: $ARCH_RAW" ;;
esac

fetch() {
  url="$1"
  out="$2"
  if [ "$DL" = "curl" ]; then
    curl -fsSL "$url" -o "$out"
  else
    wget -qO "$out" "$url"
  fi
}

# Resolver a Tag da Versão
if [ "$VERSION" = "latest" ]; then
  api_url="https://api.github.com/repos/${REPO}/releases/latest"
  tmp_json="$(mktemp)"
  fetch "$api_url" "$tmp_json"
  TAG="$(sed -n 's/^[[:space:]]*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]\+\)".*/\1/p' "$tmp_json" | head -n 1)"
  rm -f "$tmp_json"
  [ -n "$TAG" ] || fail "não foi possível resolver a versão mais recente em $api_url"
else
  TAG="$VERSION"
fi

BASE_URL="https://github.com/${REPO}/releases/download/${TAG}"
WORKDIR="$(mktemp -d)"
cleanup() { rm -rf "$WORKDIR"; }
trap cleanup EXIT INT TERM

mkdir -p "$BINDIR"

# 1. Download e Verificação de Checksums
CHECKSUMS="checksums_${TAG}.txt"
CHECKSUMS_PATH="${WORKDIR}/${CHECKSUMS}"
printf "🔍 Verificando integridade da versão %s...\n" "$TAG" >&2
fetch "${BASE_URL}/${CHECKSUMS}" "$CHECKSUMS_PATH" || fail "falha ao baixar checksums de ${BASE_URL}/${CHECKSUMS}"

sha256_file() {
  f="$1"
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$f" | awk '{print $1}'
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$f" | awk '{print $1}'
  else
    fail "necessário sha256sum ou shasum para verificação"
  fi
}

expected_sha_for() {
  name="$1"
  awk -v n="$name" '($2==n){print $1}' "$CHECKSUMS_PATH" | head -n 1
}

# 2. Download do Asset Único (bmt)
ASSET="${BIN_NAME}_${TAG}_linux_${ARCH}.tar.gz"
ASSET_PATH="${WORKDIR}/${ASSET}"

printf "📥 Baixando %s...\n" "$ASSET" >&2
fetch "${BASE_URL}/${ASSET}" "$ASSET_PATH" || fail "falha ao baixar ${BASE_URL}/${ASSET}"

# 3. Validar Checksum
EXP="$(expected_sha_for "$ASSET")"
[ -n "$EXP" ] || fail "checksum não encontrado para ${ASSET} em ${CHECKSUMS}"
GOT="$(sha256_file "$ASSET_PATH")"
[ "$GOT" = "$EXP" ] || fail "falha na verificação de checksum: esperado $EXP, obtido $GOT"

# 4. Extração e Instalação
printf "🚀 Instalando %s em %s...\n" "$BIN_NAME" "$BINDIR" >&2
tar -xzf "$ASSET_PATH" -C "$WORKDIR"
[ -f "${WORKDIR}/${BIN_NAME}" ] || fail "binário '${BIN_NAME}' não encontrado dentro do pacote"

if command -v install >/dev/null 2>&1; then
  install -m 0755 "${WORKDIR}/${BIN_NAME}" "${BINDIR}/${BIN_NAME}"
else
  cp "${WORKDIR}/${BIN_NAME}" "${BINDIR}/${BIN_NAME}"
  chmod 0755 "${BINDIR}/${BIN_NAME}"
fi

# 5. Feedback Final
printf "\n✅ %s instalado com sucesso!\n" "$BIN_NAME" >&2

if ! printf "%s" "$PATH" | grep -q "$BINDIR"; then
  printf "\n⚠️  Atenção: %s não está no seu PATH.\n" "$BINDIR" >&2
  printf "Adicione a seguinte linha ao seu ~/.bashrc ou ~/.zshrc:\n" >&2
  printf "  export PATH=\"%s:\$PATH\"\n" "$BINDIR" >&2
fi

printf "\nExperimente rodar: %s info\n" "$BIN_NAME info" >&2