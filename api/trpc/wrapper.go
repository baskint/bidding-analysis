package trpc

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// HandlerFunc is the signature for business logic handlers
// They receive a context, user UUID, and decoded request object
// and return a result and error
type HandlerFunc func(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error)

// WithAuth wraps a handler function with common middleware:
// - User authentication and validation
// - Request timeout (60 seconds)
// - Request body decoding
// - Error handling
// - Response writing
//
// Usage:
//
//	http.HandleFunc("/trpc/endpoint", h.WithAuth(h.myHandler, &MyRequest{}))
func (h *Handler) WithAuth(fn HandlerFunc, requestType interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get and validate user from context
		userID := GetUserIDFromContext(r.Context())
		if userID == "" {
			h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			log.Printf("Invalid user ID format: %v", err)
			h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// 2. Setup context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// 3. Decode request body (if POST/PUT and requestType provided)
		var req interface{} = requestType
		if requestType != nil && r.Body != nil && (r.Method == "POST" || r.Method == "PUT") {
			// Create a new instance of the request type
			reqValue := reflect.New(reflect.TypeOf(requestType).Elem()).Interface()

			if err := json.NewDecoder(r.Body).Decode(reqValue); err != nil {
				log.Printf("Failed to decode request body: %v", err)
				h.writeErrorResponse(w, "Invalid request format", http.StatusBadRequest)
				return
			}
			req = reqValue
		}

		// 4. Execute the actual handler function
		result, err := fn(ctx, userUUID, req)
		if err != nil {
			log.Printf("Handler error: %v", err)

			// Handle specific error types
			if errors.Is(err, context.DeadlineExceeded) {
				h.writeErrorResponse(w, "Request timeout. Try a smaller date range.", http.StatusGatewayTimeout)
				return
			}

			// Default error response
			h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 5. Write successful response
		h.writeTRPCResponse(w, result)
	}
}

// WithAuthNoBody is a variant for GET requests that don't need request body parsing
//
// Usage:
//
//	http.HandleFunc("/trpc/endpoint", h.WithAuthNoBody(h.myHandler))
func (h *Handler) WithAuthNoBody(fn HandlerFunc) http.HandlerFunc {
	return h.WithAuth(fn, nil)
}

// WithAuthQuery wraps a handler that reads parameters from query string
// The request object will be nil, and the handler should parse query params itself
//
// Usage:
//
//	http.HandleFunc("/trpc/endpoint", h.WithAuthQuery(h.myHandler))
func (h *Handler) WithAuthQuery(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get and validate user
		userID := GetUserIDFromContext(r.Context())
		if userID == "" {
			h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			log.Printf("Invalid user ID format: %v", err)
			h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Setup timeout
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		// Execute handler with query string access via context
		// Handler can use r.URL.Query() by accessing it from the original request
		result, err := fn(ctx, userUUID, r)
		if err != nil {
			log.Printf("Handler error: %v", err)

			if errors.Is(err, context.DeadlineExceeded) {
				h.writeErrorResponse(w, "Request timeout", http.StatusGatewayTimeout)
				return
			}

			h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.writeTRPCResponse(w, result)
	}
}
