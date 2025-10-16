package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/highway-to-Golang/user-service/internal/domain"
	"github.com/highway-to-Golang/user-service/internal/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Name == "" {
		http.Error(w, "Email and name are required", http.StatusBadRequest)
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	user, err := domain.NewUser(req.Name, req.Email, req.Role)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	if err := h.userRepo.Create(r.Context(), user); err != nil {
		slog.Error("failed to save user", "error", err)
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	slog.Info("user created successfully", "user_id", user.ID, "email", user.Email)

	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	slog.Info("getting user", "id", id)

	user, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to get user", "error", err, "user_id", id)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	slog.Info("getting all users")

	users, err := h.userRepo.GetAll(r.Context())
	if err != nil {
		slog.Error("failed to get users", "error", err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingUser, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to get user for update", "error", err, "user_id", id)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	if req.Email != nil {
		existingUser.Email = *req.Email
	}
	if req.Name != nil {
		existingUser.Name = *req.Name
	}
	if req.Role != "" {
		existingUser.Role = req.Role
	}

	slog.Info("updating user", "id", id, "email", existingUser.Email, "name", existingUser.Name)

	if err := h.userRepo.Update(r.Context(), id, existingUser); err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to update user", "error", err, "user_id", id)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	updatedUser, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("failed to get updated user", "error", err, "user_id", id)
		http.Error(w, "Failed to get updated user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	slog.Info("deleting user", "id", id)

	if err := h.userRepo.Delete(r.Context(), id); err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to delete user", "error", err, "user_id", id)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User deleted successfully",
	}

	writeJSON(w, http.StatusOK, response)
}
