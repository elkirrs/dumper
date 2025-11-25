#!/usr/bin/env sh
set -e

REPO="elkirrs/dumper"
BIN_NAME="dumper"

OS="$(uname -s)"
ARCH="$(uname -m)"

echo "Detected OS: $OS, ARCH: $ARCH"

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac


VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed -E 's/.*"tag_name": *"v?([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "❌ Cannot determine latest version from GitHub"
  exit 1
fi

echo "Version: $VERSION"

download_and_install() {
  OS_LOWER=$(printf '%s' "$OS" | tr '[:upper:]' '[:lower:]')

  URL="https://github.com/${REPO}/releases/download/v${VERSION}/${BIN_NAME}_${VERSION}_${OS_LOWER}_${ARCH}.tar.gz"
  echo "Downloading ${URL}"

  tmpdir="$(mktemp -d 2>/dev/null || mktemp -d -t dumper)"
  archive="${tmpdir}/${BIN_NAME}.tar.gz"

  if ! curl -fSL "$URL" -o "$archive"; then
    echo "Download failed: $URL" >&2
    rm -rf "$tmpdir"
    exit 1
  fi

  if ! tar -xzf "$archive" -C "$tmpdir"; then
    echo "Failed to extract archive" >&2
    rm -rf "$tmpdir"
    exit 1
  fi

  if [ ! -f "${tmpdir}/dumper" ]; then
    echo "Error: file 'dumper' not found in archive" >&2
    rm -rf "$tmpdir"
    exit 1
  fi

  if [ ! -f "${tmpdir}/config.yaml" ]; then
    echo "Error: file 'config.yaml' not found in archive" >&2
    rm -rf "$tmpdir"
    exit 1
  fi

  dest_bin="."
  if [ -w "./" ]; then
    cp "${tmpdir}/dumper" "$dest_bin"
    chmod +x "$dest_bin"
  else
    echo "Installing dumper to ${dest_bin} (using sudo)"
    sudo cp "${tmpdir}/dumper" "$dest_bin"
    sudo chmod +x "$dest_bin"
  fi

  dest_cfg="."
  if [ ! -d "$dest_cfg" ]; then
    if ! mkdir -p "$dest_cfg" 2>/dev/null; then
      echo "Creating config directory with sudo: $dest_cfg"
      sudo mkdir -p "$dest_cfg"
    fi
  fi

  if [ -w "$dest_cfg" ]; then
    cp "${tmpdir}/config.yaml" "$dest_cfg/config.yaml"
  else
    echo "Installing config.yaml to ${dest_cfg} (using sudo)"
    sudo cp "${tmpdir}/config.yaml" "$dest_cfg/config.yaml"
  fi

  echo "Installed dumper to ${dest_bin}"
  echo "Installed config.yaml to ${dest_cfg}/config.yaml"

  rm -rf "$tmpdir"
}

case "$OS" in
  Linux*)
    download_and_install
    ;;
  Darwin*)
    download_and_install
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

echo "Installation of ${BIN_NAME} ${VERSION} completed!"