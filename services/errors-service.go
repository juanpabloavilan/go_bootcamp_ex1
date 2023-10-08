package services

import "fmt"

type ServiceError struct {
	Code        string
	Description error
}

func (s ServiceError) Error() string {
	return fmt.Sprintf("Service Error -%q -%q \n", s.Code, s.Description)
}

// Implementing Is interface to use error.Is in custom error
func (s ServiceError) Is(err error) bool {
	other, ok := err.(ServiceError)
	if !ok {
		return false
	}
	return s.Code == other.Code
}
