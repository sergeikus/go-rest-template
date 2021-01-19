package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sergeikus/go-rest-template/pkg/storage"
	"github.com/sergeikus/go-rest-template/pkg/types"
	"github.com/stretchr/testify/require"
)

func marshal(obj interface{}, t *testing.T) []byte {
	b, err := json.Marshal(&obj)
	require.NoError(t, err, "failed to marshal an object: %v", err)
	return b
}

func Test_GetKey(t *testing.T) {
	tt := []struct {
		name         string
		key          string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Invalid query key (empty)",
			key:          "",
			expectedCode: 500,
			expectedBody: "key must be provided\n",
		},
		{
			name:         "Invalid query key (not a number)",
			key:          "test",
			expectedCode: 500,
			expectedBody: "key must be an integer:",
		},
		{
			name:         "Valid query",
			key:          "1",
			expectedCode: 200,
			expectedBody: string(marshal(types.Data{ID: 1, String: "test"}, t)) + "\n",
		},
	}

	api := API{
		DB: &storage.InMemoryStorage{},
	}

	err := api.DB.Connect()
	require.NoError(t, err, "expected to see no errors, but got: %v", err)

	_, err = api.DB.Store("test")
	require.NoError(t, err, "expected Store() to succeed")

	hnd := http.HandlerFunc(api.GetData)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/data/get?key=%s", tc.key), nil)
			hnd.ServeHTTP(rec, req)
			require.Equal(t, tc.expectedCode, rec.Code)
			require.Contains(t, string(rec.Body.Bytes()), tc.expectedBody)
		})
	}
}

func Test_Store(t *testing.T) {
	tt := []struct {
		name         string
		request      DataAdditionRequest
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Empty request",
			request:      DataAdditionRequest{},
			expectedCode: 500,
			expectedBody: "validation of data addition request failed: data to be added must be non-empty string\n",
		},
		{
			name: "Valid data addition request",
			request: DataAdditionRequest{
				Data: "test",
			},
			expectedCode: 200,
			expectedBody: MsgStatusOK,
		},
	}

	api := API{
		DB: &storage.InMemoryStorage{},
	}
	err := api.DB.Connect()
	require.NoError(t, err, "expected to see no errors, but got: %v", err)

	hnd := http.HandlerFunc(api.Store)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tc.request)
			require.NoError(t, err, "failed to marshal request: %", err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/store", strings.NewReader(string(bodyBytes)))
			hnd.ServeHTTP(rec, req)
			require.Equal(t, tc.expectedCode, rec.Code)
			require.Equal(t, tc.expectedBody, string(rec.Body.Bytes()))
		})
	}
}
