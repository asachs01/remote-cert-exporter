[![Test](https://github.com/asachs01/remote-cert-exporter/actions/workflows/test.yml/badge.svg)](https://github.com/asachs01/remote-cert-exporter/actions/workflows/test.yml)

# Remote Certificate Exporter

A Prometheus exporter that collects SSL/TLS certificate metrics from remote endpoints. This exporter helps monitor certificate expiration dates and provides alerts before certificates expire.

## Features

- Monitors SSL/TLS certificates from multiple remote endpoints
- Provides metrics about certificate expiration time
- Prometheus metrics exposed on HTTP endpoint
- Health check endpoint
- Configurable through YAML configuration file

## Metrics

The exporter provides the following metrics:

- `ssl_certificate_expiry_seconds{host="example.com",issuer="...",subject="..."}`: Time remaining until certificate expiration (in seconds)
- `ssl_certificate_not_after_timestamp{host="example.com",issuer="...",subject="..."}`: Unix timestamp when the certificate expires
- `ssl_certificate_scrape_errors_total{host="example.com"}`: Counter of scraping errors by host

## Installation

### Building from source

```bash
# Clone the repository
git clone https://github.com/asachs01/remote-cert-exporter
cd remote-cert-exporter

# Build the binary
make build
```

### Using Docker

```bash
make docker
```

### Systemd Service Installation

The exporter can be installed as a systemd service using the provided convenience scripts:

```bash
# Install the service (requires root/sudo)
sudo ./scripts/install.sh

# The installer will:
# 1. Create a system user and group (remote-cert-exporter)
# 2. Install the binary to /usr/local/bin
# 3. Create config directory at /etc/remote-cert-exporter
# 4. Set up logging directory at /var/log/remote-cert-exporter
# 5. Install and configure the systemd service

# After installation:
sudo systemctl start remote-cert-exporter  # Start the service
sudo systemctl enable remote-cert-exporter # Enable at boot

# Check the status
sudo systemctl status remote-cert-exporter

# View logs
sudo journalctl -u remote-cert-exporter
```

To uninstall the service:
```bash
sudo ./scripts/uninstall.sh
```

The uninstaller will:
- Stop and disable the service
- Remove the binary and service files
- Remove the system user and group
- Preserve configuration and log files (can be removed manually if desired)

## Usage

### Running the exporter

```bash
./cert-exporter --config.file=cert_exporter.yml
```

Default ports and paths:
- Listen address: `:9117`
- Metrics path: `/metrics`
- Health check: `/health`
- Probe endpoint: `/probe`

### Configuration

Create a YAML configuration file (`cert_exporter.yml`) with the following structure:

```yaml
modules:
  default:
    prober: tcp  # or http
    timeout: 5s
    port: 443
    validate_chain: true
    insecure_skip_verify: false
    client_cert:  # Optional
      cert_file: "/path/to/cert"
      key_file: "/path/to/key"
```

### Command Line Flags

- `--web.listen-address`: Address to listen on for telemetry (default: ":9117")
- `--web.telemetry-path`: Path under which to expose metrics (default: "/metrics")
- `--config.file`: Path to configuration file (default: "cert_exporter.yml")

## Production Deployment

### Docker Deployment

The recommended way to run the exporter in production is using Docker:

```bash
docker run -d \
  --name cert-exporter \
  --restart=unless-stopped \
  -p 9117:9117 \
  -v /path/to/your/cert_exporter.yml:/etc/cert-exporter/cert_exporter.yml \
  cert-exporter --config.file=/etc/cert-exporter/cert_exporter.yml
```

### Kubernetes Deployment

1. Create a ConfigMap for your configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cert-exporter-config
data:
  cert_exporter.yml: |
    modules:
      default:
        prober: tcp
        timeout: 5s
        port: 443
        validate_chain: true
        insecure_skip_verify: false
```

2. Deploy the exporter:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cert-exporter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cert-exporter
  template:
    metadata:
      labels:
        app: cert-exporter
    spec:
      containers:
      - name: cert-exporter
        image: cert-exporter:latest
        args:
          - "--config.file=/etc/cert-exporter/cert_exporter.yml"
        ports:
        - containerPort: 9117
          name: http
        volumeMounts:
        - name: config
          mountPath: /etc/cert-exporter
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 200m
            memory: 256Mi
        livenessProbe:
          httpGet:
            path: /health
            port: http
        readinessProbe:
          httpGet:
            path: /health
            port: http
      volumes:
      - name: config
        configMap:
          name: cert-exporter-config

---
apiVersion: v1
kind: Service
metadata:
  name: cert-exporter
  labels:
    app: cert-exporter
spec:
  ports:
  - port: 9117
    name: http
  selector:
    app: cert-exporter
```

3. Add Prometheus scrape configuration:

```yaml
scrape_configs:
  - job_name: 'cert-exporter'
    static_configs:
      - targets: ['cert-exporter:9117']
```

### Monitoring and Alerting

Add the following Prometheus alerting rules to be notified of expiring certificates:

```yaml
groups:
- name: CertificateAlerts
  rules:
  - alert: CertificateExpiringSoon
    expr: ssl_certificate_expiry_seconds < (14 * 24 * 3600) # 14 days
    for: 1h
    labels:
      severity: warning
    annotations:
      summary: "Certificate expiring soon for {{ $labels.host }}"
      description: "SSL certificate for {{ $labels.host }} will expire in {{ $value | humanizeDuration }}"
  
  - alert: CertificateExpired
    expr: ssl_certificate_expiry_seconds <= 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "Certificate expired for {{ $labels.host }}"
      description: "SSL certificate for {{ $labels.host }} has expired"
```

## Development

For development setup, testing, and contribution guidelines, please see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Quick Start Guide

To test the exporter quickly, follow these steps:

1. Create a test configuration file `test-config.yml`:

```yaml
modules:
  default:
    prober: tcp
    timeout: 5s
    port: 443
    validate_chain: true
    insecure_skip_verify: false
```

2. Start the exporter:

```bash
./cert-exporter --config.file=test-config.yml
```

3. Test the exporter by querying a specific domain (e.g., google.com):

```bash
curl "http://localhost:9117/probe?target=google.com"
```

4. View all metrics:

```bash
curl http://localhost:9117/metrics
```

You should see metrics like these in the output:
```
# HELP ssl_certificate_expiry_seconds Number of seconds until the SSL certificate expires
# TYPE ssl_certificate_expiry_seconds gauge
ssl_certificate_expiry_seconds{host="google.com",issuer="GTS CA 1C3",subject="*.google.com"} 7776000
```

To verify everything is working:
1. The `/health` endpoint should return a 200 status code
2. The probe endpoint should return certificate metrics for the target
3. The expiry time should be a positive number (unless the certificate is expired)

### Common Issues

- If you get connection errors, verify that:
  - The target host is accessible
  - The port is correct (default 443)
  - Your network allows the connection
- If you see certificate chain errors, you might need to set `insecure_skip_verify: true` for testing

### Quick Install (Linux)

Install using our convenience script:

```bash
# Install latest version
curl -sSL https://raw.githubusercontent.com/asachs01/remote-cert-exporter/main/scripts/get.sh | sudo bash

# Install specific version
curl -sSL https://raw.githubusercontent.com/asachs01/remote-cert-exporter/main/scripts/get.sh | VERSION=0.1.0 sudo bash

# Install to different directory
curl -sSL https://raw.githubusercontent.com/asachs01/remote-cert-exporter/main/scripts/get.sh | INSTALL_DIR=/opt/remote-cert-exporter sudo bash
```

The script will:
1. Detect your system architecture
2. Download the appropriate release
3. Create a system user and group
4. Install the binary and systemd service
5. Set up default configuration
