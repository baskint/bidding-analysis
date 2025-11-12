// api/trpc/ml_model_handlers.go
package trpc

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
)

// listMLModels retrieves all ML models for the current user
func (h *Handler) listMLModels(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Parse pagination parameters
	page := 1
	pageSize := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	response, err := h.mlModelStore.List(r.Context(), userUUID, page, pageSize)
	if err != nil {
		log.Printf("Failed to list ML models: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve ML models", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, response)
}

// getMLModel retrieves a specific ML model by ID
func (h *Handler) getMLModel(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID string `json:"id"`
	}

	// Try to get ID from query parameter first
	modelID := r.URL.Query().Get("id")
	if modelID == "" {
		// If not in query, try body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		modelID = req.ID
	}

	if modelID == "" {
		h.writeErrorResponse(w, "Model ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(modelID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid model ID format", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	model, err := h.mlModelStore.GetByID(r.Context(), id, userUUID)
	if err != nil {
		log.Printf("Failed to get ML model: %v", err)
		h.writeErrorResponse(w, "Model not found", http.StatusNotFound)
		return
	}

	h.writeTRPCResponse(w, model)
}

// createMLModel creates a new ML model
func (h *Handler) createMLModel(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	var input models.MLModelCreate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if input.Name == "" {
		h.writeErrorResponse(w, "Model name is required", http.StatusBadRequest)
		return
	}
	if input.Type == "" {
		h.writeErrorResponse(w, "Model type is required", http.StatusBadRequest)
		return
	}
	if input.Version == "" {
		h.writeErrorResponse(w, "Model version is required", http.StatusBadRequest)
		return
	}
	if input.Provider == "" {
		h.writeErrorResponse(w, "Model provider is required", http.StatusBadRequest)
		return
	}

	// Initialize config if nil
	if input.Config == nil {
		input.Config = make(map[string]interface{})
	}

	model, err := h.mlModelStore.Create(r.Context(), userUUID, &input)
	if err != nil {
		log.Printf("Failed to create ML model: %v", err)
		h.writeErrorResponse(w, "Failed to create ML model", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, model)
}

// updateMLModel updates an existing ML model
func (h *Handler) updateMLModel(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID string `json:"id"`
		models.MLModelUpdate
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		h.writeErrorResponse(w, "Model ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid model ID format", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	model, err := h.mlModelStore.Update(r.Context(), id, userUUID, &req.MLModelUpdate)
	if err != nil {
		log.Printf("Failed to update ML model: %v", err)
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, model)
}

// deleteMLModel soft deletes an ML model
func (h *Handler) deleteMLModel(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		h.writeErrorResponse(w, "Model ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid model ID format", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	if err := h.mlModelStore.Delete(r.Context(), id, userUUID); err != nil {
		log.Printf("Failed to delete ML model: %v", err)
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, map[string]interface{}{
		"success": true,
		"message": "Model deleted successfully",
	})
}

// setDefaultMLModel sets a model as the default for its type
func (h *Handler) setDefaultMLModel(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		h.writeErrorResponse(w, "Model ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid model ID format", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	if err := h.mlModelStore.SetAsDefault(r.Context(), id, userUUID); err != nil {
		log.Printf("Failed to set default ML model: %v", err)
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated model
	model, err := h.mlModelStore.GetByID(r.Context(), id, userUUID)
	if err != nil {
		h.writeErrorResponse(w, "Model updated but failed to retrieve", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, model)
}

// getDefaultMLModel retrieves the default model for a specific type
func (h *Handler) getDefaultMLModel(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		Type string `json:"type"`
	}

	// Try query parameter first
	modelType := r.URL.Query().Get("type")
	if modelType == "" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		modelType = req.Type
	}

	if modelType == "" {
		h.writeErrorResponse(w, "Model type is required", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	model, err := h.mlModelStore.GetDefaultForType(r.Context(), userUUID, modelType)
	if err != nil {
		log.Printf("Failed to get default ML model: %v", err)
		h.writeErrorResponse(w, "No default model found for this type", http.StatusNotFound)
		return
	}

	h.writeTRPCResponse(w, model)
}
