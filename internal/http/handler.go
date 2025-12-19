package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/highway-to-Golang/user-service/internal/domain"
	apperrors "github.com/highway-to-Golang/user-service/internal/errors"
	"github.com/highway-to-Golang/user-service/internal/usecase"
)

type UserHandler struct {
	uc *usecase.UseCase
}

func NewUserHandler(uc *usecase.UseCase) *UserHandler {
	return &UserHandler{
		uc: uc,
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeErrorJSON(w http.ResponseWriter, status int, message string) {
	errorResponse := map[string]interface{}{
		"error": message,
	}
	writeJSON(w, status, errorResponse)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := r.Header.Get("Idempotency-Key")

	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request body", "error", err)
		writeErrorJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.uc.CreateUser(r.Context(), idempotencyKey, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrRequestAlreadyInProgress) {
			writeErrorJSON(w, http.StatusUnprocessableEntity, "Request already in progress")
			return
		}
		slog.Error("failed to create user", "error", err)
		writeErrorJSON(w, http.StatusUnprocessableEntity, "Failed to create user")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := h.uc.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeErrorJSON(w, http.StatusNotFound, "User not found")
			return
		}
		slog.Error("failed to get user", "error", err, "user_id", id)
		writeErrorJSON(w, http.StatusUnprocessableEntity, "Failed to get user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.GetAllUsers(r.Context())
	if err != nil {
		slog.Error("failed to get users", "error", err)
		writeErrorJSON(w, http.StatusUnprocessableEntity, "Failed to get users")
		return
	}

	response := map[string]interface{}{
		"users": users,
		"total": len(users),
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request body", "error", err)
		writeErrorJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.uc.UpdateUser(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeErrorJSON(w, http.StatusNotFound, "User not found")
			return
		}
		slog.Error("failed to update user", "error", err, "user_id", id)
		writeErrorJSON(w, http.StatusUnprocessableEntity, "Failed to update user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.uc.DeleteUser(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeErrorJSON(w, http.StatusNotFound, "User not found")
			return
		}
		slog.Error("failed to delete user", "error", err, "user_id", id)
		writeErrorJSON(w, http.StatusUnprocessableEntity, "Failed to delete user")
		return
	}

	response := map[string]interface{}{
		"message": "User deleted successfully",
	}

	writeJSON(w, http.StatusOK, response)
}
