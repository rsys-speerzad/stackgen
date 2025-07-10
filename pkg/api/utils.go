package api

import (
	"encoding/json"
	"net/http"
)

func ResponseWriter(w http.ResponseWriter, data interface{}, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if jsonMsg, err := json.Marshal(data); err == nil {
			w.Write(jsonMsg)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
