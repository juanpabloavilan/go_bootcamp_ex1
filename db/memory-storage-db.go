package db

import (
	"example/bootcamp_ex1/entities"
	"fmt"

	"github.com/google/uuid"
)

type memoryStorage[T entities.StorageObject] struct {
	entities map[uuid.UUID]T
}

func NewMemoryStorage[T entities.StorageObject]() *memoryStorage[T] {
	return &memoryStorage[T]{
		entities: make(map[uuid.UUID]T),
	}
}

func (m *memoryStorage[T]) Create(thing T) (uuid.UUID, error) {
	id := thing.GetId()
	m.entities[id] = thing
	return id, nil
}

func (m *memoryStorage[T]) Get(key uuid.UUID) (T, error) {
	var zeroValue T
	value, ok := m.entities[key]
	//If tning doesn't exist we return a nil value and a error
	if !ok {
		return zeroValue, StorageError{
			Code:        ErrEntityNotFound,
			Description: fmt.Sprintf("cannot find entity with id: %q", key),
		}
	}
	return value, nil
}

func (m *memoryStorage[T]) GetAll() ([]T, error) {
	things := make([]T, 0, len(m.entities))
	for _, t := range m.entities {
		things = append(things, t)
	}
	return things, nil
}

func (m *memoryStorage[T]) Update(key uuid.UUID, thing T) (T, error) {
	var zeroValue T
	// If not exists return error
	_, err := m.Get(key)
	if err != nil {
		return zeroValue, err
	}
	m.entities[key] = thing

	return m.entities[key], nil
}
func (m *memoryStorage[T]) Delete(key uuid.UUID) (uuid.UUID, error) {
	// If not exists return error
	_, err := m.Get(key)
	if err != nil {
		return uuid.Nil, err
	}
	// delete
	delete(m.entities, key)
	return key, nil
}
