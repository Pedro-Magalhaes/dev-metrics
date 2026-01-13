#!/usr/bin/env sh
set -eu

REPO_DEFAULT="Pedro-Magalhaes/dev-metrics"
REPO="${REPO:-$REPO_DEFAULT}"
BINDIR="${BINDIR:-${HOME}/.local/bin}"
VERSION="${VERSION:-latest}"

fail() {
  printf '%s\n' "$*" >&2
  exit 1
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "missing dependency: $1"
}

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
  fail "need curl or wget"
fi

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH_RAW="$(uname -m)"

case "$OS" in
  linux) : ;;
  *) fail "unsupported OS: $OS (only linux supported)" ;;
 esac

case "$ARCH_RAW" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) fail "unsupported arch: $ARCH_RAW (supported: x86_64, aarch64)" ;;
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

# Resolve tag
if [ "$VERSION" = "latest" ]; then
  api_url="https://api.github.com/repos/${REPO}/releases/latest"
  tmp_json="$(mktemp)"
  fetch "$api_url" "$tmp_json"
  # Extract tag_name without jq.
  TAG="$(
    sed -n 's/^[[:space:]]*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]\+\)".*/\1/p' "$tmp_json" | head -n 1
  )"
  rm -f "$tmp_json"
  [ -n "$TAG" ] || fail "could not resolve latest release tag from $api_url"
else
  TAG="$VERSION"
fi

BASE_URL="https://github.com/${REPO}/releases/download/${TAG}"

TOOLS="build-meter analyze-meter export-meter"
WORKDIR="$(mktemp -d)"
cleanup() { rm -rf "$WORKDIR"; }
trap cleanup EXIT INT TERM

mkdir -p "$BINDIR"

# Download checksums
CHECKSUMS="checksums_${TAG}.txt"
CHECKSUMS_PATH="${WORKDIR}/${CHECKSUMS}"
fetch "${BASE_URL}/${CHECKSUMS}" "$CHECKSUMS_PATH" || fail "failed to download ${BASE_URL}/${CHECKSUMS}"

sha256_file() {
  f="$1"
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$f" | awk '{print $1}'
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$f" | awk '{print $1}'
  else
    fail "need sha256sum (or shasum) for checksum verification"
  fi
}

expected_sha_for() {
  name="$1"
  # Format: <sha>  <filename>
  awk -v n="$name" '($2==n){print $1}' "$CHECKSUMS_PATH" | head -n 1
}

printf '%s\n' "Installing dev-metrics tools from ${REPO}@${TAG} for linux/${ARCH}..." >&2

for tool in $TOOLS; do
  asset="${tool}_${TAG}_linux_${ARCH}.tar.gz"
  asset_path="${WORKDIR}/${asset}"

  fetch "${BASE_URL}/${asset}" "$asset_path" || fail "failed to download ${BASE_URL}/${asset}"

  exp="$(expected_sha_for "$asset")"
  [ -n "$exp" ] || fail "checksum not found for ${asset} in ${CHECKSUMS}"

  got="$(sha256_file "$asset_path")"
  [ "$got" = "$exp" ] || fail "checksum mismatch for ${asset}"

  tar -xzf "$asset_path" -C "$WORKDIR"
  [ -f "${WORKDIR}/${tool}" ] || fail "expected binary '${tool}' inside ${asset}"
  chmod 0755 "${WORKDIR}/${tool}"

  if command -v install >/dev/null 2>&1; then
    install -m 0755 "${WORKDIR}/${tool}" "${BINDIR}/${tool}"
  else
    cp "${WORKDIR}/${tool}" "${BINDIR}/${tool}"
    chmod 0755 "${BINDIR}/${tool}"
  fi

  printf '%s\n' "Installed: ${BINDIR}/${tool}" >&2
done

printf '\n%s\n' "Add this to your shell rc (bash/zsh):" >&2
printf '%s\n' "export PATH=\"$HOME/.local/bin:\$PATH\"" >&2
printf '%s\n' "Then restart your shell and run: build-meter -h" >&2
