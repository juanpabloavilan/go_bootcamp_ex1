package db

import (
	"fmt"
)

type StorageError struct {
	Code        string
	Description string
}

func (s StorageError) Error() string {
	return fmt.Sprintf("Storage Error -%q -%q\n", s.Code, s.Description)
}

// Implementing Is interface to use error.Is in custom error
func (s StorageError) Is(err error) bool {
	other, ok := err.(StorageError)
	if !ok {
		return false
	}
	return s.Code == other.Code
}

// Error codes
const (
	RedisConnectionFailed = "RedisConnectionFailed "

	ErrEntityNotFound = "EntityNotFound"

	ErrUnmarshalingEntity = "UnmarshalingEntityError"

	ErrMarshalingEntity = "MarshalingEntityError"

	ErrGettingRecords = "ErrGettingRecords"
)
