package http

import (
	"net/http"
)

func NewRouter(userHandler *UserHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/users", userHandler.GetAllUsers)
	mux.HandleFunc("POST /api/users", userHandler.CreateUser)
	mux.HandleFunc("GET /api/users/{id}", userHandler.GetUser)
	mux.HandleFunc("PUT /api/users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", userHandler.DeleteUser)

	return mux
}
