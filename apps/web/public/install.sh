#!/usr/bin/env sh
set -eu

REPO="Createitv/agc-cli"
BINARY="agc"
INSTALL_DIR="${AGC_INSTALL_DIR:-$HOME/.local/bin}"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"

case "$os" in
  darwin|linux) ;;
  *) echo "unsupported OS: $os" >&2; exit 1 ;;
esac

case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "unsupported architecture: $arch" >&2; exit 1 ;;
esac

tag="${AGC_VERSION:-latest}"
if [ "$tag" = "latest" ]; then
  api_url="https://api.github.com/repos/$REPO/releases/latest"
  tag="$(curl -fsSL "$api_url" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n 1)"
fi

if [ -z "$tag" ]; then
  echo "could not resolve latest agc-cli release" >&2
  exit 1
fi

version="${tag#v}"
archive="agc-cli_${version}_${os}_${arch}.tar.gz"
base_url="https://github.com/$REPO/releases/download/$tag"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

curl -fsSLo "$tmpdir/$archive" "$base_url/$archive"
curl -fsSLo "$tmpdir/checksums.txt" "$base_url/checksums.txt"

(cd "$tmpdir" && grep "  $archive\$" checksums.txt | shasum -a 256 -c -)
tar -xzf "$tmpdir/$archive" -C "$tmpdir" "$BINARY"

mkdir -p "$INSTALL_DIR"
install -m 0755 "$tmpdir/$BINARY" "$INSTALL_DIR/$BINARY"

echo "agc installed to $INSTALL_DIR/$BINARY"
echo "Add $INSTALL_DIR to PATH if needed."
