package db

import (
	"example/bootcamp_ex1/entities"

	"github.com/google/uuid"
)

type Storage interface {
	Get(id uuid.UUID) (entities.User, error)
	GetAll() ([]entities.User, error)
	Create(entities.User) (uuid.UUID, error)
	Update(id uuid.UUID, user entities.User) (entities.User, error)
	Delete(id uuid.UUID) (uuid.UUID, error)
}
