package conf

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sergeikus/go-rest-template/pkg/auth"
	"github.com/sergeikus/go-rest-template/pkg/storage"
	"gopkg.in/yaml.v2"
)

// Conf represents server configuration files
type Conf struct {
	TLS           bool          `yaml:"tls"`
	TLSKeyPath    string        `yaml:"tlsKeyPath,omitempty"`
	TLSCertPath   string        `yaml:"tlsCertPath,omitempty"`
	Port          int           `yaml:"port"`
	Database      Database      `yaml:"database"`
	Authorization Authorization `yaml:"authorization"`
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
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database configuration validation failed: %v", err)
	}
	if err := c.Authorization.Validate(); err != nil {
		return fmt.Errorf("authorization configuration validation failed: %v", err)
	}
	return nil
}

// Database represents database configuration
type Database struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Name     string `yaml:"name,omitempty"`
}

// Validate performs validation of database configuration
func (d *Database) Validate() error {
	if len(d.Type) == 0 {
		return fmt.Errorf("database type musy be non-empty string")
	}
	if strings.ToLower(d.Type) != storage.DatabaseTypeInMemory && strings.ToLower(d.Type) != storage.DatabaseTypePostgre {
		return fmt.Errorf("unsupported database type: %s", d.Type)
	}
	if strings.ToLower(d.Type) != storage.DatabaseTypeInMemory {
		if len(d.Host) == 0 {
			return fmt.Errorf("host must be non-empty string")
		}
		if d.Port == 0 {
			return fmt.Errorf("port must be not 0")
		}
		if len(d.Username) == 0 {
			return fmt.Errorf("username must be non-empty string")
		}
		if len(d.Password) == 0 {
			return fmt.Errorf("password must be non-empty string")
		}
		if len(d.Name) == 0 {
			return fmt.Errorf("database name must be non-empty string")
		}
	}
	return nil
}

// Authorization represents a server authorization parameters
type Authorization struct {
	Type             string `yaml:"type,omitempty"`
	SessionDuration  int    `yaml:"sessionDuration"`
	PBKDF2Iterations int    `yaml:"pbkdf2Iterations"`
	PBKDF2KeyLenght  int    `yaml:"pbkdf2KeyLenght"`
}

// Validate performs authorization parameters validation
func (a *Authorization) Validate() error {
	if len(a.Type) == 0 {
		return fmt.Errorf("authorization type must be provided")
	}
	switch strings.ToLower(a.Type) {
	case auth.SSMType:
		if a.SessionDuration <= 0 {
			return fmt.Errorf("session duration in 'session' type must be greater than 0")
		}
	default:
		return fmt.Errorf("unknown authorization type: %s", a.Type)
	}
	if a.PBKDF2Iterations <= 0 {
		return fmt.Errorf("PBKDF2 hashing iterations count must be greater than 0")
	}
	if a.PBKDF2KeyLenght <= 0 {
		return fmt.Errorf("PBKDF2 key lenght must be greater than 0")
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
