package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := Response{
		Result: "",
		Error:  message,
	}

	json.NewEncoder(w).Encode(errorResponse)
}

func sendResponse(w http.ResponseWriter, statusCode int, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Result: responseData,
		Error:  "",
	}

	json.NewEncoder(w).Encode(response)
}

func sendResponseRaw(w http.ResponseWriter, statusCode int, responseData json.RawMessage) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Result: responseData,
		Error:  "",
	}

	json.NewEncoder(w).Encode(response)
}
