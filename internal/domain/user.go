package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/highway-to-Golang/user-service/internal/errors"
)

var (
	ErrNotFound = errors.ErrNotFound
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type UpdateUserRequest struct {
	Email *string `json:"email,omitempty"`
	Name  *string `json:"name,omitempty"`
	Role  string  `json:"role,omitempty"`
}

func NewUser(name, email, role string) (User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return User{}, errors.ErrFailedToBuild
	}
	return User{
		ID:    id.String(),
		Name:  name,
		Email: email,
		Role:  role,
	}, nil
}
