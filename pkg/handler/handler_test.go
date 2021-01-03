package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sergeikus/go-rest-template/pkg/storage"
	"github.com/stretchr/testify/require"
)

func Test_GetKey(t *testing.T) {
	tt := []struct {
		name         string
		key          string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Invalid query key",
			key:          "",
			expectedCode: 500,
			expectedBody: "key must be provided\n",
		},
		{
			name:         "Valid query",
			key:          "OK",
			expectedCode: 200,
			expectedBody: "{\"status\": \"OK\"}",
		},
	}

	api := API{
		DB: &storage.InMemoryStorage{},
	}
	err := api.DB.Connect()
	require.NoError(t, err, "expected to see no errors, but got: %v", err)

	hnd := http.HandlerFunc(api.GetKey)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/data/get?key=%s", tc.key), nil)
			hnd.ServeHTTP(rec, req)
			require.Equal(t, tc.expectedCode, rec.Code)
			require.Equal(t, tc.expectedBody, string(rec.Body.Bytes()))
		})
	}
}

func Test_Store(t *testing.T) {
	tt := []struct {
		name         string
		request      KeyAdditionRequest
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Empty request",
			request:      KeyAdditionRequest{},
			expectedCode: 500,
			expectedBody: "validation of key addition request failed: key must be non-empty string\n",
		},
		{
			name: "Empty request",
			request: KeyAdditionRequest{
				Key:  "test",
				Data: "test",
			},
			expectedCode: 200,
			expectedBody: "{\"status\": \"OK\"}",
		},
	}

	api := API{
		DB: &storage.InMemoryStorage{},
	}
	err := api.DB.Connect()
	require.NoError(t, err, "expected to see no errors, but got: %v", err)

	hnd := http.HandlerFunc(api.AddKey)
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
