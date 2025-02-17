# Remote Certificate Exporter

A Prometheus exporter that collects SSL/TLS certificate metrics from remote endpoints. This exporter helps monitor certificate expiration dates and provides alerts before certificates expire.

## Features

- Monitors SSL/TLS certificates from multiple remote endpoints
- Provides metrics about certificate expiration time
- Supports custom configuration for different monitoring scenarios
- Prometheus metrics exposed on HTTP endpoint
- Health check endpoint
- Configurable through YAML configuration file

## Metrics

The exporter provides the following metrics:

- `ssl_certificate_expiry_seconds`: Time remaining until certificate expiration (in seconds)
- `ssl_certificate_not_after_timestamp`: Unix timestamp when the certificate expires
- `ssl_certificate_scrape_errors_total`: Counter of scraping errors by host

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

Configuration options:
- `prober`: Protocol to use (tcp or http)
- `timeout`: Overall timeout for the check
- `port`: Default port if not specified in target
- `proxy_url`: Optional HTTP proxy
- `validate_chain`: Whether to validate the entire certificate chain
- `insecure_skip_verify`: Skip certificate validation (not recommended for production)
- `client_cert`: Optional client certificate configuration

### Command Line Flags

- `--web.listen-address`: Address to listen on for telemetry (default: ":9117")
- `--web.telemetry-path`: Path under which to expose metrics (default: "/metrics")
- `--config.file`: Path to configuration file (default: "cert_exporter.yml")

## Development

### Running Tests

```bash
make test
```

### Code Coverage

```bash
make coverage
```

### Linting

```bash
make lint
```

## License

MIT License

This project is open source and available under the MIT License. You are free to use, modify, and distribute this software, provided that you include the original copyright notice and attribution in any copies or substantial portions of the software.

## Contributing

1. Fork the repository
2. Create a new branch for your feature or bug fix (`git checkout -b feature/your-feature-name`)
3. Make your changes
4. Ensure all tests pass by running `make test`
5. Commit your changes
6. Push to your fork
7. Submit a Pull Request

Please ensure your PR description clearly describes the changes and the motivation for the changes.
