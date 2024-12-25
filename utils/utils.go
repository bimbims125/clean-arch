package utils

import (
	"encoding/json"
	"net/http"
)

// ResponseData represent the response data struct
type ResponseData struct {
	Data interface{} `json:"data"`
}

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ResponseSuccess represent the response success struct
type ResponseSuccess struct {
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, ResponseError{Message: message})
}

func RespondWithSuccess(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, ResponseSuccess{Message: message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
