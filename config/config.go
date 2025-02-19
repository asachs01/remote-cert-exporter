package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Modules map[string]*Module `yaml:"modules"`
	Targets []string          `yaml:"targets"`
}

type Module struct {
	Prober             string        `yaml:"prober"`           // tcp or http
	Timeout            time.Duration `yaml:"timeout"`          // Overall timeout for the check
	Port               int          `yaml:"port"`             // Default port if not specified in target
	ProxyURL           string        `yaml:"proxy_url"`        // Optional HTTP proxy
	ValidateChain      bool          `yaml:"validate_chain"`   // Whether to validate the entire chain
	InsecureSkipVerify bool          `yaml:"insecure_skip_verify"` // Skip certificate validation
	ClientCert         *ClientCert   `yaml:"client_cert"`      // Optional client certificate
}

type ClientCert struct {
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// LoadConfig loads the configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
} 