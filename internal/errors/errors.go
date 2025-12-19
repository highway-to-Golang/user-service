package errors

import "errors"

var (
	ErrNotFound                 = errors.New("not found")
	ErrFailedToBuild            = errors.New("failed to build")
	ErrInvalidInput             = errors.New("invalid input")
	ErrRequestAlreadyInProgress = errors.New("request already in progress")
)
