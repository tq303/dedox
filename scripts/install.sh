#!/bin/sh
set -e

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

[ "$ARCH" = "x86_64" ] && ARCH="amd64"
[ "$ARCH" = "aarch64" ] && ARCH="arm64"

if [ "$OS" != "darwin" ] && [ "$OS" != "linux" ]; then
    echo "unsupported OS: $OS"
    exit 1
fi

URL="https://github.com/tq303/dedox/releases/latest/download/ddx-${OS}-${ARCH}"
DEST="/usr/local/bin/ddx"

echo "Installing ddx for ${OS}/${ARCH}..."
curl -sL "$URL" -o "$DEST" && chmod +x "$DEST"
echo "Installed to $DEST"
