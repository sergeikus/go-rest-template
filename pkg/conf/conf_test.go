package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Conf_Validate(t *testing.T) {
	tt := []struct {
		name     string
		c        Conf
		fail     bool
		expected string
	}{
		{
			name: "Empty conf",
			c:    Conf{},
			fail: true,
		},
		{
			name: "Port not provided (tls disabled)",
			c: Conf{
				TLS: false,
			},
			fail:     true,
			expected: "port can't be 0 (verify that it's specified in the configuration)",
		},
		{
			name: "Database is not provided (tls disabled)",
			c: Conf{
				TLS:  false,
				Port: 8080,
			},
			fail:     true,
			expected: "database configuration validation failed:",
		},
		{
			name: "Valid conf (tls disabled)",
			c: Conf{
				TLS:      false,
				Port:     8080,
				Database: Database{Type: "in-memory"},
				Authorization: Authorization{
					Type:             "session",
					SessionDuration:  10,
					PBKDF2Iterations: 1,
					PBKDF2KeyLenght:  1,
				},
			},
			fail: false,
		},
		{
			name: "TLS key is not provided (tls enabled)",
			c: Conf{
				TLS: true,
			},
			fail:     true,
			expected: "TLS key path must be provided",
		},
		{
			name: "TLS certificate is not provided (tls enabled)",
			c: Conf{
				TLS:        true,
				TLSKeyPath: "test",
			},
			fail:     true,
			expected: "TLS certificate path must be provided",
		},
		{
			name: "Valid conf (tls enabled)",
			c: Conf{
				TLS:         true,
				TLSKeyPath:  "test",
				TLSCertPath: "test",
				Port:        8080,
				Database:    Database{Type: "in-memory"},
				Authorization: Authorization{
					Type:             "session",
					SessionDuration:  10,
					PBKDF2Iterations: 1,
					PBKDF2KeyLenght:  1,
				},
			},
			fail: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.c.Validate()
			if tc.fail {
				require.NotNil(t, err, "expected to see an error, but got nil")
				require.Contains(t, err.Error(), tc.expected, "expected to see a different error")
			} else {
				require.NoError(t, err, "expected to get no error, but got: %v", err)
			}
		})
	}
}

func Test_Database_Validate(t *testing.T) {
	tt := []struct {
		name     string
		d        Database
		fail     bool
		expected string
	}{
		{
			name:     "Empty type",
			d:        Database{},
			fail:     true,
			expected: "database type musy be non-empty string",
		},
		{
			name: "Unknown type",
			d: Database{
				Type: "test",
			},
			fail:     true,
			expected: "unsupported database type:",
		},
		{
			name: "Valid database (in-memory)",
			d: Database{
				Type: "in-memory",
			},
			fail: false,
		},
		{
			name: "Host is empty (not in-memory)",
			d: Database{
				Type: "postgres",
			},
			fail:     true,
			expected: "host must be non-empty string",
		},
		{
			name: "Port is 0 (not in-memory)",
			d: Database{
				Type: "postgres",
				Host: "test",
			},
			fail:     true,
			expected: "port must be not 0",
		},
		{
			name: "Username is empty (not in-memory)",
			d: Database{
				Type: "postgres",
				Host: "test",
				Port: 8443,
			},
			fail:     true,
			expected: "username must be non-empty string",
		},
		{
			name: "Password is empty (not in-memory)",
			d: Database{
				Type:     "postgres",
				Host:     "test",
				Port:     8443,
				Username: "test",
			},
			fail:     true,
			expected: "password must be non-empty string",
		},
		{
			name: "Name is empty (not in-memory)",
			d: Database{
				Type:     "postgres",
				Host:     "test",
				Port:     8443,
				Username: "test",
				Password: "test",
			},
			fail:     true,
			expected: "database name must be non-empty string",
		},
		{
			name: "Valid database (not in-memory)",
			d: Database{
				Type:     "postgres",
				Host:     "test",
				Port:     8443,
				Username: "test",
				Password: "test",
				Name:     "test",
			},
			fail: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.d.Validate()
			if tc.fail {
				require.NotNil(t, err, "expected to see an error, but got nil")
				require.Contains(t, err.Error(), tc.expected, "expected to see a different error")
			} else {
				require.NoError(t, err, "expected to get no error, but got: %v", err)
			}
		})
	}
}

func Test_AuthorizationValidate(t *testing.T) {
	tt := []struct {
		name     string
		a        Authorization
		fail     bool
		expected string
	}{
		{
			name:     "Empty struct",
			a:        Authorization{},
			fail:     true,
			expected: "authorization type must be provided",
		},
		{
			name: "Unknown type",
			a: Authorization{
				Type: "TESTING",
			},
			fail:     true,
			expected: "unknown authorization type:",
		},
		{
			name: "Invalid session duration (session)",
			a: Authorization{
				Type:            "session",
				SessionDuration: -1,
			},
			fail:     true,
			expected: "session duration in 'session' type must be greater than 0",
		},
		{
			name: "Invalid PBKDF2 iterations",
			a: Authorization{
				Type:            "session",
				SessionDuration: 10,
			},
			fail:     true,
			expected: "PBKDF2 hashing iterations count must be greater than 0",
		},
		{
			name: "Invalid PBKDF2 key lenght",
			a: Authorization{
				Type:             "session",
				SessionDuration:  10,
				PBKDF2Iterations: 1,
			},
			fail:     true,
			expected: "PBKDF2 key lenght must be greater than 0",
		},
		{
			name: "Valid authorization configuration (session)",
			a: Authorization{
				Type:             "session",
				SessionDuration:  10,
				PBKDF2Iterations: 1,
				PBKDF2KeyLenght:  1,
			},
			fail: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.a.Validate()
			if tc.fail {
				require.NotNil(t, err, "expected to see an error, but got nil")
				require.Contains(t, err.Error(), tc.expected, "expected to see a different error")
			} else {
				require.NoError(t, err, "expected to get no error, but got: %v", err)
			}
		})
	}
}

const (
	testFolder        = "../../test/"
	testConfigsFolder = testFolder + "configs/"
	badConfFolder     = testConfigsFolder + "bad/"
	goodConfFolder    = testConfigsFolder + "good/"
)

func Test_ReadConf(t *testing.T) {
	tt := []struct {
		name     string
		path     string
		fail     bool
		expected string
	}{
		{
			name:     "Empty path",
			path:     "",
			fail:     true,
			expected: "configuration path must be non-empty string",
		},
		{
			name:     "Wrong path",
			path:     "test",
			fail:     true,
			expected: "failed to read file:",
		},
		{
			name:     "Invalid conf structure",
			path:     badConfFolder + "config-1.yaml",
			fail:     true,
			expected: "failed to unmarshal configuration yaml:",
		},
		{
			name: "Valid config structure",
			path: goodConfFolder + "config.yaml",
			fail: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ReadConf(tc.path)
			if tc.fail {
				require.NotNil(t, err, "expected to see an error, but got nil")
				require.Contains(t, err.Error(), tc.expected, "expected to see a different error")
			} else {
				require.NoError(t, err, "expected to get no error, but got: %v", err)
			}
		})
	}
}
