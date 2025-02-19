package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLoggers(t *testing.T) {
	// Capture stdout and stderr
	var outBuf, errBuf bytes.Buffer
	Info = log.New(&outBuf, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(&errBuf, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Test Info logger
	Info.Println("test info message")
	if !strings.Contains(outBuf.String(), "INFO: ") {
		t.Error("Info log doesn't contain INFO prefix")
	}
	if !strings.Contains(outBuf.String(), "test info message") {
		t.Error("Info log doesn't contain test message")
	}

	// Test Error logger
	Error.Println("test error message")
	if !strings.Contains(errBuf.String(), "ERROR: ") {
		t.Error("Error log doesn't contain ERROR prefix")
	}
	if !strings.Contains(errBuf.String(), "test error message") {
		t.Error("Error log doesn't contain test message")
	}
}
