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

## Development

For development setup, testing, and contribution guidelines, please see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
