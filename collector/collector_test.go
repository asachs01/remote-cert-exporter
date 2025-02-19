package collector

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/asachs01/remote-cert-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

func TestCertificateCollector(t *testing.T) {
	// Create test certificates
	cert, priv, err := generateTestCert()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	// Start test TLS server
	listener, err := tls.Listen("tcp", "localhost:0", &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  priv,
		}},
	})
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			t.Errorf("Failed to close listener: %v", err)
		}
	}()

	// Accept connections in background
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			if err := conn.Close(); err != nil {
				t.Errorf("Failed to close connection: %v", err)
			}
		}
	}()

	// Test the collector
	module := &config.Module{
		Timeout: 5 * time.Second,
		Port:    listener.Addr().(*net.TCPAddr).Port,
	}

	collector := NewCertificateCollector("localhost", module)

	// Test metric collection
	ch := make(chan prometheus.Metric, 10)
	done := make(chan bool)

	go func() {
		collector.Collect(ch)
		close(ch)
		done <- true
	}()

	// Wait for collection to complete or timeout
	select {
	case <-done:
		// Success
	case <-time.After(10 * time.Second):
		t.Fatal("Collection timed out")
	}

	// Verify metrics were collected
	metrics := make([]prometheus.Metric, 0)
	for m := range ch {
		metrics = append(metrics, m)
	}

	if len(metrics) == 0 {
		t.Error("No metrics were collected")
	}
}

func generateTestCert() (*x509.Certificate, *rsa.PrivateKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "test.local",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, priv, nil
}
