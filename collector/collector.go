package collector

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "net"
    "net/http"
    "net/url"
    "strings"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "asachs01/remote-cert-exporter/config"
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
        // ... other metric descriptors ...
    }
}

func (c *CertificateCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- c.certExpirySeconds
    // ... other descriptors ...
}

func (c *CertificateCollector) Collect(ch chan<- prometheus.Metric) {
    timeout := c.module.Timeout
    if timeout == 0 {
        timeout = 10 * time.Second
    }

    dialer := &net.Dialer{
        Timeout: timeout,
    }

    // Handle proxy if configured
    var transport *http.Transport
    if c.module.ProxyURL != "" {
        proxyURL, err := url.Parse(c.module.ProxyURL)
        if err == nil {
            transport = &http.Transport{
                Proxy: http.ProxyURL(proxyURL),
                DialContext: dialer.DialContext,
            }
        }
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
        ch <- prometheus.NewMetric(
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
        
        // Basic expiry information
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

        // Additional certificate details...
    }
} 