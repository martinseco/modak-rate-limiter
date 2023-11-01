package utils

import (
	"net/http"
)

func HandleResponse(w http.ResponseWriter, s int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
}
