package config

import (
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name        string
		configYaml  string
		expectError bool
	}{
		{
			name: "valid config",
			configYaml: `
modules:
  default:
    prober: tcp
    timeout: 5s
    port: 443
`,
			expectError: false,
		},
		{
			name: "invalid timeout",
			configYaml: `
modules:
  default:
    prober: tcp
    timeout: invalid
`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: Write to temp file and test LoadConfig
		})
	}
} 