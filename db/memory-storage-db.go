package db

import (
	"errors"
	"example/bootcamp_ex1/entities"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("cannot find a user with this id")
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
	value, ok := m.entities[key]
	//If user doesn't exist we return a nil value and a error
	if !ok {
		var zeroValue T
		return zeroValue, ErrUserNotFound
	}
	return value, nil
}

func (u *memoryStorage[T]) GetAll() ([]T, error) {
	userList := make([]T, 0, len(u.entities))
	for _, user := range u.entities {
		userList = append(userList, user)
	}
	return userList, nil
}

func (u *memoryStorage[T]) Update(key uuid.UUID, newUser T) (T, error) {
	// If not exists return error
	_, err := u.Get(key)
	if err != nil {
		var zeroValue T
		return zeroValue, ErrUserNotFound
	}
	u.entities[key] = newUser

	return u.entities[key], nil
}
func (u *memoryStorage[T]) Delete(key uuid.UUID) (uuid.UUID, error) {
	// If not exists return error
	_, err := u.Get(key)
	if err != nil {
		return uuid.Nil, err
	}
	// delete
	delete(u.entities, key)
	return key, nil
}
