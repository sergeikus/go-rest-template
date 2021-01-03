package storage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Connect(t *testing.T) {
	tt := []struct {
		name     string
		ims      InMemoryStorage
		fail     bool
		expected string
	}{
		{
			name: "Valid connect",
			ims:  InMemoryStorage{},
			fail: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ims.Connect()
			if tc.fail {
				require.NotNil(t, err, "expected to see an error, but got nil")
				require.Contains(t, err.Error(), tc.expected, "expected to see a different error")
			} else {
				require.NoError(t, err, "expected to get no error, but got: %v", err)
			}
		})
	}
}

func Test_Store(t *testing.T) {
	tt := []struct {
		name     string
		ims      InMemoryStorage
		key      string
		data     string
		fail     bool
		expected string
	}{
		{
			name: "Invalid key",
			ims: InMemoryStorage{
				data:  make(map[string]string),
				mutex: sync.Mutex{},
			},
			key:      "",
			data:     "",
			fail:     true,
			expected: "key must be non-empty string",
		},
		{
			name: "Existing key",
			ims: InMemoryStorage{
				data: map[string]string{
					"test": "test",
				},
				mutex: sync.Mutex{},
			},
			key:      "test",
			data:     "",
			fail:     true,
			expected: "key with 'test' ID already exists",
		},
		{
			name: "Valid addition",
			ims: InMemoryStorage{
				data:  make(map[string]string),
				mutex: sync.Mutex{},
			},
			key:  "test",
			data: "test",
			fail: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ims.Store(tc.key, tc.data)
			if tc.fail {
				require.NotNil(t, err, "expected to see an error, but got nil")
				require.Contains(t, err.Error(), tc.expected, "expected to see a different error")
			} else {
				require.NoError(t, err, "expected to get no error, but got: %v", err)
			}
		})
	}
}
