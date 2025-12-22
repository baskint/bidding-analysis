package trpc

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// generateToken creates a JWT token for the given user
func (h *Handler) generateToken(userID, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

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

// Helper function to parse date range with defaults
func parseDateRange(startDateStr, endDateStr string) (time.Time, time.Time) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30) // Default to last 30 days

	// Handle Z or lack of explicit time by using ParseInLocation and setting time to 00:00:00
	location, _ := time.LoadLocation("UTC")

	if startDateStr != "" {
		if parsed, err := time.ParseInLocation("2006-01-02", startDateStr, location); err == nil {
			startDate = parsed
		} else if parsed, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = parsed.In(location)
		}
	}

	if endDateStr != "" {
		if parsed, err := time.ParseInLocation("2006-01-02", endDateStr, location); err == nil {
			endDate = parsed
		} else if parsed, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = parsed.In(location)
		}
	}

	// Ensure endDate includes full day (up to 23:59:59.999...)
	// Only add if the date parsing was ambiguous (e.g., "YYYY-MM-DD")
	// If RFC3339 was used (like in the curl), it already has time info, but we apply the shift here defensively.
	endDate = endDate.Add(24 * time.Hour).Add(-1 * time.Second)

	return startDate, endDate
}
