#!/bin/bash
set -e

# Default values
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/remote-cert-exporter"
SERVICE_USER="remote-cert-exporter"
SERVICE_GROUP="remote-cert-exporter"
LOG_DIR="/var/log/remote-cert-exporter"

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

# Copy binary
cp remote-cert-exporter "$INSTALL_DIR/"
chmod 755 "$INSTALL_DIR/remote-cert-exporter"

# Copy default config if it doesn't exist
if [ ! -f "$CONFIG_DIR/remote-cert-exporter.yml" ]; then
    cp remote-cert-exporter.yml "$CONFIG_DIR/"
fi

# Set permissions
chown -R "$SERVICE_USER:$SERVICE_GROUP" "$CONFIG_DIR"
chown -R "$SERVICE_USER:$SERVICE_GROUP" "$LOG_DIR"
chmod 644 "$CONFIG_DIR/remote-cert-exporter.yml"

# Install systemd service
cp scripts/remote-cert-exporter.service /etc/systemd/system/
chmod 644 /etc/systemd/system/remote-cert-exporter.service

# Reload systemd
systemctl daemon-reload

echo "Installation complete!"
echo "1. Edit the configuration file at $CONFIG_DIR/remote-cert-exporter.yml"
echo "2. Start the service: systemctl start remote-cert-exporter"
echo "3. Enable at boot: systemctl enable remote-cert-exporter" 