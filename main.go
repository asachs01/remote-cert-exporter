package main

import (
    "crypto/tls"
    "flag"
    "fmt"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "asachs01/remote-cert-exporter/middleware"
    "asachs01/remote-cert-exporter/logger"
    "asachs01/remote-cert-exporter/config"
)

var (
    listenAddress = flag.String("web.listen-address", ":9117", "Address to listen on for telemetry")
    metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")

    // Prometheus metrics
    certExpirySeconds = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "ssl_certificate_expiry_seconds",
            Help: "Number of seconds until the SSL certificate expires",
        },
        []string{"host", "issuer", "subject"},
    )
    
    certNotAfterTimestamp = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "ssl_certificate_not_after_timestamp",
            Help: "Timestamp when the SSL certificate expires",
        },
        []string{"host", "issuer", "subject"},
    )

    scrapeErrorsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "ssl_certificate_scrape_errors_total",
            Help: "Total number of scrape errors",
        },
        []string{"host"},
    )
)

func init() {
    prometheus.MustRegister(certExpirySeconds)
    prometheus.MustRegister(certNotAfterTimestamp)
    prometheus.MustRegister(scrapeErrorsTotal)
}

type Exporter struct {
    targets []string
}

func NewExporter(targets []string) *Exporter {
    return &Exporter{
        targets: targets,
    }
}

func (e *Exporter) scrapeTarget(target string) error {
    conf := &tls.Config{
        InsecureSkipVerify: false,
    }

    conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", target), conf)
    if err != nil {
        return fmt.Errorf("failed to connect: %v", err)
    }
    defer conn.Close()

    cert := conn.ConnectionState().PeerCertificates[0]
    
    // Calculate expiry time
    expirySeconds := time.Until(cert.NotAfter).Seconds()
    
    certExpirySeconds.With(prometheus.Labels{
        "host":    target,
        "issuer":  cert.Issuer.CommonName,
        "subject": cert.Subject.CommonName,
    }).Set(expirySeconds)

    certNotAfterTimestamp.With(prometheus.Labels{
        "host":    target,
        "issuer":  cert.Issuer.CommonName,
        "subject": cert.Subject.CommonName,
    }).Set(float64(cert.NotAfter.Unix()))

    return nil
}

func (e *Exporter) collectMetrics() {
    for _, target := range e.targets {
        if err := e.scrapeTarget(target); err != nil {
            scrapeErrorsTotal.With(prometheus.Labels{"host": target}).Inc()
        }
    }
}

func main() {
    configFile := flag.String("config.file", "cert_exporter.yml", "Path to configuration file")
    flag.Parse()

    config, err := config.LoadConfig(*configFile)
    if err != nil {
        logger.Error.Fatalf("Error loading config: %s", err)
    }

    exporter := NewExporter(config)

    // Create router with middleware
    router := http.NewServeMux()
    
    // Add handlers with middleware
    router.Handle(*metricsPath, middleware.InstrumentHandler(promhttp.Handler()))
    router.Handle("/probe", middleware.InstrumentHandler(http.HandlerFunc(exporter.probeHandler)))
    router.Handle("/health", middleware.InstrumentHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })))

    // Wrap entire router with error handling
    handler := middleware.ErrorHandler(router)

    // Start server
    logger.Info.Printf("Starting SSL certificate exporter on %s\n", *listenAddress)
    if err := http.ListenAndServe(*listenAddress, handler); err != nil {
        logger.Error.Fatalf("Error starting server: %s\n", err)
    }
} 