package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Validate(t *testing.T) {
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
			name: "Database type is not provided (tls disabled)",
			c: Conf{
				TLS:  false,
				Port: 8080,
			},
			fail:     true,
			expected: "database type must be provided",
		},
		{
			name: "Valid conf (tls disabled)",
			c: Conf{
				TLS:          false,
				Port:         8080,
				DatabaseType: "test",
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
				TLS:          false,
				TLSKeyPath:   "test",
				TLSCertPath:  "test",
				Port:         8080,
				DatabaseType: "test",
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
