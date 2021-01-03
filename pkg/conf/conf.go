package conf

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Conf represents server configuration files
type Conf struct {
	TLS          bool   `yaml:"tls"`
	TLSKeyPath   string `yaml:"tlsKeyPath,omitempty"`
	TLSCertPath  string `yaml:"tlsCertPath,omitempty"`
	Port         int    `yaml:"port"`
	DatabaseType string `yaml:"databaseType"`
}

// Validate performs configuration validation
func (c *Conf) Validate() error {
	if c.TLS {
		if len(c.TLSKeyPath) == 0 {
			return fmt.Errorf("TLS key path must be provided")
		}
		if len(c.TLSCertPath) == 0 {
			return fmt.Errorf("TLS certificate path must be provided")
		}
	}
	if c.Port == 0 {
		return fmt.Errorf("port can't be 0 (verify that it's specified in the configuration)")
	}
	if len(c.DatabaseType) == 0 {
		return fmt.Errorf("database type must be provided")
	}
	return nil
}

// ReadConf reads configuration
func ReadConf(path string) (c Conf, err error) {
	if len(path) == 0 {
		return c, fmt.Errorf("configuration path must be non-empty string")
	}
	cb, err := ioutil.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("failed to read file: %v", err)
	}
	if err := yaml.Unmarshal(cb, &c); err != nil {
		return c, fmt.Errorf("failed to unmarshal configuration yaml: %v", err)
	}
	return c, nil
}
