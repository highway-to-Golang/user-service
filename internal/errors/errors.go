package errors

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserIDEmpty    = errors.New("user ID cannot be empty")
	ErrInvalidRole    = errors.New("invalid role")
	ErrFailedToBuild  = errors.New("failed to build user")
	ErrFailedToSave   = errors.New("failed to save user")
	ErrFailedToGet    = errors.New("failed to get user")
	ErrFailedToRemove = errors.New("failed to remove user")
)

func ErrorWithID(err error, id string) error {
	return errors.New(err.Error() + " with ID " + id)
}

func ErrorWithRole(err error, role string) error {
	return errors.New(err.Error() + ": " + role)
}
