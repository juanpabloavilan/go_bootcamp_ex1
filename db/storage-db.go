package db

import (
	"example/bootcamp_ex1/entities"

	"github.com/google/uuid"
)

type Storage[T entities.StorageObject] interface {
	Get(id uuid.UUID) (T, error)
	GetAll() ([]T, error)
	Create(thing T) (uuid.UUID, error)
	Update(id uuid.UUID, thing T) (T, error)
	Delete(id uuid.UUID) (uuid.UUID, error)
}
