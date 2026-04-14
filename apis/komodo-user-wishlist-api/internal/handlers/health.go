package handlers

import (
	"encoding/json"
	"net/http"
)

func HealthHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(map[string]string{"status": "OK"})
}
