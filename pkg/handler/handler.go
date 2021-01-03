package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const getTag = "GetKey"

// GetKey queries some key in database
func (api *API) GetKey(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		key := r.URL.Query().Get("key")
		if len(key) == 0 {
			fail(w, getTag, fmt.Errorf("key must be provided"), 500)
			return
		}

		if _, err := w.Write([]byte("{\"status\": \"OK\"}")); err != nil {
			fail(w, getTag, fmt.Errorf("failed to write response: %v", err), 500)
			return
		}
	}
}

const addTag = "AddKey"

// AddKey performs key addition to the database
func (api *API) AddKey(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var kar KeyAdditionRequest
		if err := decoder.Decode(&kar); err != nil {
			fail(w, addTag, fmt.Errorf("error while decoding body request: %v", err), 500)
			return
		}

		if err := kar.Validate(); err != nil {
			fail(w, addTag, fmt.Errorf("validation of key addition request failed: %v", err), 500)
			return
		}

		if err := api.DB.Store(kar.Key, kar.Data); err != nil {
			fail(w, addTag, fmt.Errorf("failed to store key with '%s' ID: %v", kar.Key, err), 500)
			return
		}

		if _, err := w.Write([]byte("{\"status\": \"OK\"}")); err != nil {
			fail(w, getTag, fmt.Errorf("failed to write response: %v", err), 500)
			return
		}
	}
}

func fail(w http.ResponseWriter, tag string, err error, code int) {
	log.Printf("[%s] %v", tag, err)
	http.Error(w, err.Error(), code)
}
