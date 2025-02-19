#!/bin/bash
set -e

# Default values
VERSION=${VERSION:-"latest"}
INSTALL_DIR=${INSTALL_DIR:-"/usr/local/bin"}
CONFIG_DIR=${CONFIG_DIR:-"/etc/remote-cert-exporter"}
SERVICE_USER=${SERVICE_USER:-"remote-cert-exporter"}
SERVICE_GROUP=${SERVICE_GROUP:-"remote-cert-exporter"}
LOG_DIR=${LOG_DIR:-"/var/log/remote-cert-exporter"}
GITHUB_REPO="asachs01/remote-cert-exporter"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo_info() { echo -e "${BLUE}INFO: ${NC}$1"; }
echo_success() { echo -e "${GREEN}SUCCESS: ${NC}$1"; }
echo_error() { echo -e "${RED}ERROR: ${NC}$1"; }

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo_error "This script must be run as root"
    exit 1
fi

# Create temporary directory
TMP_DIR=$(mktemp -d)
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

# Determine system architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    *) echo_error "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Determine OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case $OS in
    linux) ;;
    *) echo_error "Unsupported operating system: $OS"; exit 1 ;;
esac

# Download latest release if version not specified
if [ "$VERSION" = "latest" ]; then
    VERSION=$(curl -s https://api.github.com/repos/${GITHUB_REPO}/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")')
fi
VERSION=${VERSION#v} # Remove 'v' prefix if present

echo_info "Installing remote-cert-exporter version ${VERSION}"

# Download and extract release
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/remote-cert-exporter_${OS}_${ARCH}.tar.gz"
echo_info "Downloading from ${DOWNLOAD_URL}"
curl -L -o "$TMP_DIR/remote-cert-exporter.tar.gz" "$DOWNLOAD_URL"
tar xzf "$TMP_DIR/remote-cert-exporter.tar.gz" -C "$TMP_DIR"

# Create user and group if they don't exist
if ! getent group "$SERVICE_GROUP" >/dev/null; then
    groupadd --system "$SERVICE_GROUP"
fi

if ! getent passwd "$SERVICE_USER" >/dev/null; then
    useradd --system \
        --gid "$SERVICE_GROUP" \
        --no-create-home \
        --home-dir /nonexistent \
        --shell /sbin/nologin \
        --comment "Remote Certificate Exporter Service User" \
        "$SERVICE_USER"
fi

# Create necessary directories
mkdir -p "$INSTALL_DIR"
mkdir -p "$CONFIG_DIR"
mkdir -p "$LOG_DIR"

# Install binary
cp "$TMP_DIR/remote-cert-exporter" "$INSTALL_DIR/"
chmod 755 "$INSTALL_DIR/remote-cert-exporter"

# Download and install systemd service
curl -L -o "/etc/systemd/system/remote-cert-exporter.service" \
    "https://raw.githubusercontent.com/${GITHUB_REPO}/v${VERSION}/scripts/remote-cert-exporter.service"
chmod 644 "/etc/systemd/system/remote-cert-exporter.service"

# Download default config if it doesn't exist
if [ ! -f "$CONFIG_DIR/remote-cert-exporter.yml" ]; then
    curl -L -o "$CONFIG_DIR/remote-cert-exporter.yml" \
        "https://raw.githubusercontent.com/${GITHUB_REPO}/v${VERSION}/remote-cert-exporter.yml"
fi

# Set permissions
chown -R "$SERVICE_USER:$SERVICE_GROUP" "$CONFIG_DIR"
chown -R "$SERVICE_USER:$SERVICE_GROUP" "$LOG_DIR"
chmod 644 "$CONFIG_DIR/remote-cert-exporter.yml"

# Reload systemd
systemctl daemon-reload

echo_success "Installation complete!"
echo "
Next steps:
1. Edit the configuration file at $CONFIG_DIR/remote-cert-exporter.yml
2. Start the service: systemctl start remote-cert-exporter
3. Enable at boot: systemctl enable remote-cert-exporter

To verify installation:
- systemctl status remote-cert-exporter
- curl http://localhost:9117/health
" 