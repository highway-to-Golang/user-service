package client

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/highway-to-Golang/user-service/internal/domain"
)

//go:embed testdata/create_user.json
var createUserData []byte

//go:embed testdata/update_user.json
var updateUserData []byte

func LoadCreateUserRequest() (domain.CreateUserRequest, error) {
	var req domain.CreateUserRequest
	if err := json.Unmarshal(createUserData, &req); err != nil {
		return domain.CreateUserRequest{}, fmt.Errorf("failed to unmarshal create user data: %w", err)
	}
	return req, nil
}

func LoadUpdateUserRequest() (domain.UpdateUserRequest, error) {
	var req domain.UpdateUserRequest
	if err := json.Unmarshal(updateUserData, &req); err != nil {
		return domain.UpdateUserRequest{}, fmt.Errorf("failed to unmarshal update user data: %w", err)
	}
	return req, nil
}
