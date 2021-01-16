package storage

import (
	"sync"
	"testing"

	"github.com/sergeikus/go-rest-template/pkg/types"
	"github.com/stretchr/testify/require"
)

func Test_Connect(t *testing.T) {
	tt := []struct {
		name     string
		ims      *InMemoryStorage
		fail     bool
		expected string
	}{
		{
			name: "Valid connect",
			ims:  &InMemoryStorage{},
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
		ims      *InMemoryStorage
		data     string
		fail     bool
		expected string
	}{
		{
			name: "Invalid data",
			ims: &InMemoryStorage{
				data:  make(map[int]types.Data),
				mutex: sync.Mutex{},
				index: 1,
			},
			data:     "",
			fail:     true,
			expected: "data must be non-empty string",
		},
		{
			name: "Valid addition",
			ims: &InMemoryStorage{
				data:  make(map[int]types.Data),
				mutex: sync.Mutex{},
				index: 1,
			},
			data: "test",
			fail: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			id, err := tc.ims.Store(tc.data)
			if tc.fail {
				require.NotNil(t, err, "expected to see an error, but got nil")
				require.Contains(t, err.Error(), tc.expected, "expected to see a different error")
			} else {
				require.NoError(t, err, "expected to get no error, but got: %v", err)
				require.NotEqual(t, 0, id, "expected to see a different ID")
			}
		})
	}
}
