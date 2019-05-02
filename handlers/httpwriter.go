package handlers

import (
	"encoding/json"
	"net/http"
)

func Write(w http.ResponseWriter, code int, response interface{}) {
	b, err := json.Marshal(response)
	if err != nil {
		Error(w, http.StatusInternalServerError, ErrMarshalling.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func Error(w http.ResponseWriter, code int, message string) {
	data := struct {
		Failure string `json:"failure"`
	}{
		message,
	}

	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}
