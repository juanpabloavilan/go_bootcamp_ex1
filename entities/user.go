package entities

import (
	"github.com/google/uuid"
)

type StorageObject interface {
	GetId() uuid.UUID
}

type User struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	LastName string    `json:"lastname"`
	Email    string    `json:"email"`
	Active   bool      `json:"active"`
	Address  Address   `json:"address"`
}

func (u User) GetId() uuid.UUID {
	return u.Id
}

type UserRequest struct {
	Name     string  `json:"name" validate:"required"`
	LastName string  `json:"lastname" validate:"required"`
	Email    string  `json:"email" validate:"required"`
	Active   bool    `json:"active"`
	Address  Address `json:"address" validate:"required"`
}

type Address struct {
	City          string `json:"city" validate:"required"`
	Country       string `json:"country" validate:"required"`
	AddressString string `json:"address_string" validate:"required"`
}
