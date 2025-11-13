package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/highway-to-Golang/user-service/config"
	"github.com/highway-to-Golang/user-service/internal/client"
)

func main() {
	cfg, err := config.NewClientConfig()
	if err != nil {
		slog.Error("Error loading client config", "error", err)
		os.Exit(1)
	}

	cl, err := client.NewClient(*cfg)
	if err != nil {
		slog.Error("Error creating client", "error", err)
		os.Exit(1)
	}
	defer cl.Close()

	ctx := context.Background()

	createReq, err := client.LoadCreateUserRequest()
	if err != nil {
		slog.Error("Error loading create user request", "error", err)
		os.Exit(1)
	}

	user, err := cl.CreateUser(ctx, createReq, "test-idempotency-key-1")
	if err != nil {
		slog.Error("Error creating user", "error", err)
		os.Exit(1)
	}

	fmt.Println("Created user:")
	printJSON(user)

	userID := user.ID
	user, err = cl.GetUser(ctx, userID)
	if err != nil {
		slog.Error("Error getting user", "error", err)
		os.Exit(1)
	}

	fmt.Println("\nGot user:")
	printJSON(user)

	updateReq, err := client.LoadUpdateUserRequest()
	if err != nil {
		slog.Error("Error loading update user request", "error", err)
		os.Exit(1)
	}

	user, err = cl.UpdateUser(ctx, userID, updateReq)
	if err != nil {
		slog.Error("Error updating user", "error", err)
		os.Exit(1)
	}

	fmt.Println("\nUpdated user:")
	printJSON(user)

	users, err := cl.GetAllUsers(ctx)
	if err != nil {
		slog.Error("Error getting all users", "error", err)
		os.Exit(1)
	}

	fmt.Println("\nAll users:")
	printJSON(users)

	err = cl.DeleteUser(ctx, userID)
	if err != nil {
		slog.Error("Error deleting user", "error", err)
		os.Exit(1)
	}

	fmt.Println("\nUser deleted successfully")
}

func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		slog.Error("Error marshaling JSON", "error", err)
		return
	}
	fmt.Println(string(data))
}
