package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const getTag = "GetData"

// GetData queries some key in database
func (api *API) GetData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tm := time.Now()
		keyString := r.URL.Query().Get("key")
		if len(keyString) == 0 {
			fail(w, getTag, fmt.Errorf("key must be provided"), http.StatusInternalServerError)
			return
		}

		key, err := strconv.Atoi(keyString)
		if err != nil {
			fail(w, getTag, fmt.Errorf("key must be an integer: %v", err), http.StatusInternalServerError)
			return
		}
		d, err := api.DB.GetKey(key)
		if err != nil {
			fail(w, getTag, fmt.Errorf("failed to get data for '%d' key: %v", key, err), http.StatusInternalServerError)
			return
		}

		writeReponseObject(w, d, getTag, fmt.Sprintf("Successfully got key with '%d' ID, action took: %v", key, time.Since(tm)))
	}
}

const getAllTag = "GetAllData"

// GetAllData queries all data from main table
func (api *API) GetAllData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tm := time.Now()

		data, err := api.DB.GetAll()
		if err != nil {
			fail(w, getAllTag, fmt.Errorf("failed to get all data from 'data_table': %v", err), http.StatusInternalServerError)
			return
		}
		writeReponseObject(w, data, getAllTag, fmt.Sprintf("Successfully got all data from 'data_table', got %d rows, action took: %v", len(data), time.Since(tm)))
	}
}

const storeTag = "Store"

// Store performs key addition to the database
func (api *API) Store(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var dar DataAdditionRequest
		if err := decoder.Decode(&dar); err != nil {
			fail(w, storeTag, fmt.Errorf("error while decoding body request: %v", err), http.StatusInternalServerError)
			return
		}

		if err := dar.Validate(); err != nil {
			fail(w, storeTag, fmt.Errorf("validation of data addition request failed: %v", err), http.StatusInternalServerError)
			return
		}

		if _, err := api.DB.Store(dar.Data); err != nil {
			fail(w, storeTag, fmt.Errorf("failed to store data: %v", err), http.StatusInternalServerError)
			return
		}

		writeResponseString(w, MsgStatusOK, getTag, "")
	}
}

func fail(w http.ResponseWriter, tag string, err error, code int) {
	log.Printf("[%s] Error: %v", tag, err)
	http.Error(w, err.Error(), code)
}

func writeResponseString(w http.ResponseWriter, msg, logTag, logMsg string) {
	if _, err := w.Write([]byte(msg)); err != nil {
		fail(w, logTag, fmt.Errorf("failed to write response: %v", err), http.StatusInternalServerError)
		return
	}
	if len(logMsg) != 0 {
		log.Printf("[%s] %s", logTag, logMsg)
	}
}

func writeReponseObject(w http.ResponseWriter, obj interface{}, logTag, logMsg string) {
	if err := json.NewEncoder(w).Encode(&obj); err != nil {
		fail(w, logTag, fmt.Errorf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
	if len(logMsg) != 0 {
		log.Printf("[%s] %s", logTag, logMsg)
	}
}
