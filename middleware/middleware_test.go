package middleware

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/testutil"
)

func TestInstrumentHandler(t *testing.T) {
    // Reset the metrics
    prometheus.Unregister(requestDuration)
    prometheus.Unregister(requestsTotal)
    prometheus.MustRegister(requestDuration)
    prometheus.MustRegister(requestsTotal)

    // Create a test handler
    handler := InstrumentHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _, err := w.Write([]byte("OK"))
        if err != nil {
            t.Errorf("Failed to write response: %v", err)
        }
    }))

    // Create a test request
    req := httptest.NewRequest("GET", "/test", nil)
    rr := httptest.NewRecorder()

    // Serve the request
    handler.ServeHTTP(rr, req)

    // Check response
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Check metrics
    if testutil.CollectAndCount(requestsTotal) == 0 {
        t.Error("No requests were recorded")
    }
    if testutil.CollectAndCount(requestDuration) == 0 {
        t.Error("No request duration was recorded")
    }
}

func TestErrorHandler(t *testing.T) {
    // Test panic recovery
    handler := ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        panic("test panic")
    }))

    req := httptest.NewRequest("GET", "/test", nil)
    rr := httptest.NewRecorder()

    // This should not panic
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusInternalServerError {
        t.Errorf("handler returned wrong status code: got %v want %v", 
            status, http.StatusInternalServerError)
    }

    if !strings.Contains(rr.Body.String(), "Internal Server Error") {
        t.Error("Expected 'Internal Server Error' in response")
    }
} 