package handler

import (
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

func fail(w http.ResponseWriter, tag string, err error, code int) {
	log.Printf("[%s] %v", tag, err)
	http.Error(w, err.Error(), code)
}
