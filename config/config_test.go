package config

import (
	"os"
	"testing"
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
			tmpfile, err := os.CreateTemp("", "config-*.yml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			// Write config and check error
			if _, err := tmpfile.Write([]byte(tc.configYaml)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// Test LoadConfig
			_, err = LoadConfig(tmpfile.Name())
			if tc.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
} 