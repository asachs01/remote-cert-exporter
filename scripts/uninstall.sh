#!/bin/bash
set -e

# Stop and disable service
systemctl stop remote-cert-exporter || true
systemctl disable remote-cert-exporter || true

# Remove files
rm -f /usr/local/bin/remote-cert-exporter
rm -f /etc/systemd/system/remote-cert-exporter.service

# Optionally remove config and logs (commented out for safety)
# rm -rf /etc/remote-cert-exporter
# rm -rf /var/log/remote-cert-exporter

# Remove user and group
userdel remote-cert-exporter || true
groupdel remote-cert-exporter || true

# Reload systemd
systemctl daemon-reload

echo "Uninstallation complete!" 