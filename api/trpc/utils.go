package trpc

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// parseInput parses request input from either JSON body (POST) or query params (GET)
func parseInput(r *http.Request, input interface{}) error {
	if r.Method == "GET" {
		// Parse query parameters for GET requests
		return parseQueryParams(r, input)
	}

	// Parse JSON body for POST requests
	return json.NewDecoder(r.Body).Decode(input)
}

// parseQueryParams converts URL query parameters to struct
func parseQueryParams(r *http.Request, input interface{}) error {
	// This is a simplified implementation
	// In a real tRPC setup, you'd parse the query parameters properly
	query := r.URL.Query()

	// Convert to JSON and back for simplicity
	jsonData := make(map[string]interface{})
	for key, values := range query {
		if len(values) > 0 {
			// Try to parse as int, fallback to string
			if intVal, err := strconv.Atoi(values[0]); err == nil {
				jsonData[key] = intVal
			} else {
				jsonData[key] = values[0]
			}
		}
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, input)
}

// writeSuccess writes a successful tRPC response
func writeSuccess(w http.ResponseWriter, data interface{}) {
	response := TRPCResponse{
		Result: &TRPCResult{
			Data: data,
			Type: "data",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// writeError writes an error tRPC response
func writeError(w http.ResponseWriter, code int, message string, err error) {
	var errorData interface{}
	if err != nil {
		errorData = err.Error()
	}

	response := TRPCResponse{
		Error: &TRPCError{
			Code:    code,
			Message: message,
			Data:    errorData,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
