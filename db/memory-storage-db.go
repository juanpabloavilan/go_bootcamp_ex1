package db

import (
	"errors"
	"example/bootcamp_ex1/entities"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("cannot find a user with this id")
)

type userMap map[uuid.UUID]entities.User

type memoryStorage struct {
	users userMap
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{users: userMap{}}
}

func (u *memoryStorage) Create(user entities.User) (uuid.UUID, error) {
	id := user.Id
	u.users[id] = user
	return id, nil
}

func (u *memoryStorage) Get(key uuid.UUID) (entities.User, error) {
	value, ok := u.users[key]
	//If user doesn't exist we return a nil value and a error
	if !ok {
		return entities.User{}, ErrUserNotFound
	}
	return value, nil
}

func (u *memoryStorage) GetAll() ([]entities.User, error) {
	userList := make([]entities.User, 0, len(u.users))
	for _, user := range u.users {
		userList = append(userList, user)
	}
	return userList, nil
}

func (u *memoryStorage) Update(key uuid.UUID, newUser entities.User) (entities.User, error) {
	// If not exists return error
	_, err := u.Get(key)
	if err != nil {
		return entities.User{}, err
	}
	u.users[key] = newUser

	return u.users[key], nil
}
func (u *memoryStorage) Delete(key uuid.UUID) (uuid.UUID, error) {
	// If not exists return error
	_, err := u.Get(key)
	if err != nil {
		return uuid.Nil, err
	}
	// delete
	delete(u.users, key)
	return key, nil
}
