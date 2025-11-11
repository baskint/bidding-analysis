// api/trpc/auth_handlers.go
package trpc

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Auth handlers
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from database
	user, err := h.userStore.GetByUsername(req.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		User: User{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if len(req.Username) < 3 || len(req.Password) < 6 {
		http.Error(w, "Username must be at least 3 characters and password at least 6 characters", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	if _, err := h.userStore.GetByUsername(req.Username); err == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create user
	user, err := h.userStore.Create(req.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		User: User{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)

	// Use the safe assertion (comma-ok idiom) to prevent the panic,
	// but don't immediately fail the tRPC call if it's missing, as
	// the auth middleware should have caught this first.

	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok || userID == "" {
		// Log the failure, but return the error in the format the frontend expects
		// (often by letting the rest of the request handler fail, or using the tRPC error structure).
		// For now, let's proceed with a known-good structure and assume the error
		// is due to premature error handling.
		log.Println("WARNING: UserID was missing or empty in getMe context.")
		// You can use a generic 401 or let the subsequent logic handle the error:
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// CRITICAL: REMOVE the second check for username that we added, as it may be line 128
	// and was causing confusion, and is not needed for the subsequent GetByID call.

	user, err := h.userStore.GetByID(userID)
	if err != nil {
		// This is where the application normally returns a 404/database error if the user ID is valid but not in the DB
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]User{
		"user": {
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
	}

	// Since the cURL works, this part is likely fine:
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
