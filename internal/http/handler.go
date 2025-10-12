package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	slog.Info("creating user", "email", req["email"], "name", req["name"])

	response := map[string]interface{}{
		"id":      1,
		"email":   req["email"],
		"name":    req["name"],
		"message": "User created successfully",
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	slog.Info("getting user", "id", id)

	response := map[string]interface{}{
		"id":    id,
		"email": "user@example.com",
		"name":  "John Doe",
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	slog.Info("getting all users")

	response := map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": 1, "email": "user1@example.com", "name": "John Doe"},
			{"id": 2, "email": "user2@example.com", "name": "Jane Smith"},
		},
		"total": 2,
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	slog.Info("updating user", "id", id, "email", req["email"], "name", req["name"])

	response := map[string]interface{}{
		"id":      id,
		"email":   req["email"],
		"name":    req["name"],
		"message": "User updated successfully",
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	slog.Info("deleting user", "id", id)

	response := map[string]interface{}{
		"message": "User deleted successfully",
	}

	writeJSON(w, http.StatusOK, response)
}
