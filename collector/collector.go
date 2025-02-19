package collector

import (
    "crypto/tls"
    "fmt"
    "net"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/asachs01/remote-cert-exporter/config"
)

type CertificateCollector struct {
    target string
    module *config.Module
    
    // Metrics
    certExpirySeconds *prometheus.Desc
    certNotAfterTimestamp *prometheus.Desc
    certChainLength *prometheus.Desc
    certSerial *prometheus.Desc
    certKeyUsage *prometheus.Desc
    certError *prometheus.Desc
}

func NewCertificateCollector(target string, module *config.Module) *CertificateCollector {
    return &CertificateCollector{
        target: target,
        module: module,
        certExpirySeconds: prometheus.NewDesc(
            "ssl_certificate_expiry_seconds",
            "Number of seconds until the SSL certificate expires",
            []string{"host", "issuer", "subject", "serial", "position"},
            nil,
        ),
        certNotAfterTimestamp: prometheus.NewDesc(
            "ssl_certificate_not_after_timestamp",
            "Timestamp when the SSL certificate expires",
            []string{"host", "issuer", "subject", "serial", "position"},
            nil,
        ),
        certChainLength: prometheus.NewDesc(
            "ssl_certificate_chain_length",
            "Number of certificates in the chain",
            []string{"host"},
            nil,
        ),
        certSerial: prometheus.NewDesc(
            "ssl_certificate_serial_number",
            "Serial number of the certificate",
            []string{"host", "issuer", "subject", "serial"},
            nil,
        ),
        certKeyUsage: prometheus.NewDesc(
            "ssl_certificate_key_usage",
            "Key usage of the certificate",
            []string{"host", "issuer", "subject", "serial", "usage"},
            nil,
        ),
        certError: prometheus.NewDesc(
            "ssl_certificate_error",
            "Error encountered while collecting certificate metrics",
            []string{"host", "error"},
            nil,
        ),
    }
}

func (c *CertificateCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- c.certExpirySeconds
    ch <- c.certNotAfterTimestamp
    ch <- c.certChainLength
    ch <- c.certSerial
    ch <- c.certKeyUsage
    ch <- c.certError
}

func (c *CertificateCollector) Collect(ch chan<- prometheus.Metric) {
    timeout := c.module.Timeout
    if timeout == 0 {
        timeout = 10 * time.Second
    }

    dialer := &net.Dialer{
        Timeout: timeout,
    }

    // Configure TLS
    tlsConfig := &tls.Config{
        InsecureSkipVerify: c.module.InsecureSkipVerify,
        ServerName:         c.target,
    }

    // Add client certificates if configured
    if c.module.ClientCert != nil {
        cert, err := tls.LoadX509KeyPair(
            c.module.ClientCert.CertFile,
            c.module.ClientCert.KeyFile,
        )
        if err == nil {
            tlsConfig.Certificates = []tls.Certificate{cert}
        }
    }

    // Connect and get certificates
    port := c.module.Port
    if port == 0 {
        port = 443
    }

    conn, err := tls.DialWithDialer(dialer, "tcp", 
        fmt.Sprintf("%s:%d", c.target, port), tlsConfig)
    if err != nil {
        ch <- prometheus.MustNewConstMetric(
            c.certError,
            prometheus.GaugeValue,
            1.0,
            c.target,
            err.Error(),
        )
        return
    }
    defer conn.Close()

    // Process certificates
    certs := conn.ConnectionState().PeerCertificates
    for i, cert := range certs {
        position := fmt.Sprintf("%d", i)
        
        ch <- prometheus.MustNewConstMetric(
            c.certExpirySeconds,
            prometheus.GaugeValue,
            time.Until(cert.NotAfter).Seconds(),
            c.target,
            cert.Issuer.CommonName,
            cert.Subject.CommonName,
            cert.SerialNumber.String(),
            position,
        )
    }
} 