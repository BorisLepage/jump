package helper_api

import (
	"encoding/json"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func SendErrorResponse(w http.ResponseWriter, errorType, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := ErrorResponse{
		Error:   errorType,
		Message: message,
		Code:    statusCode,
	}

	json.NewEncoder(w).Encode(errorResponse)
}

// Validator interface for objects that can validate themselves
type Validator interface {
	Validate() error
}

// ReadAndValidate reads JSON from request body and validates the object
func ReadAndValidate(r *http.Request, v any) error {
	// Read and unmarshal JSON
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, v); err != nil {
		return err
	}

	// Validate if the object implements Validator
	if validator, ok := v.(Validator); ok {
		return validator.Validate()
	}

	return nil
}
