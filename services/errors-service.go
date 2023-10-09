package services

import "fmt"

type ServiceError struct {
	Code        string
	Description error
}

func (s ServiceError) Error() string {
	return fmt.Sprintf("Service Error -%q -%q \n", s.Code, s.Description)
}
