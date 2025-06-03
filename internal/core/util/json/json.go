package json

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func ResponseWithSuccess(w http.ResponseWriter, data interface{}) {
	response := Response{
		Data: data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func ResponseWithError(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := ErrorResponse{
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		http.Error(w, "Error encoding error response", http.StatusInternalServerError)
		return
	}
}

func NewDecoder(r *http.Request) *json.Decoder {
	decoder := json.NewDecoder(r.Body)
	return decoder
}

func JoinErrors(errors []string) string {
	if len(errors) == 0 {
		return ""
	}
	var result string
	for _, err := range errors {
		if result != "" {
			result += "; "
		}
		result += err
	}
	return result
}
