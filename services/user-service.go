package services

import (
	"example/bootcamp_ex1/db"
	"example/bootcamp_ex1/entities"

	"log/slog"

	"github.com/google/uuid"
)

type UserService struct {
	storage db.Storage[entities.User]
}

func NewUserService(storage db.Storage[entities.User]) *UserService {
	userService := new(UserService)
	userService.storage = storage
	return userService
}

func (u *UserService) Get(id uuid.UUID) (entities.User, error) {
	//Log action
	slog.Info("Getting a user by id", "id", id.String())
	user, err := u.storage.Get(id)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (u *UserService) GetAll() ([]entities.User, error) {
	//Log action
	slog.Info("Logging all users")
	//Return slice of users
	return u.storage.GetAll()
}

func (u *UserService) Create(userReq entities.UserRequest) (uuid.UUID, error) {
	id := uuid.New()
	newUser := entities.User{
		Id:       id,
		Name:     userReq.Name,
		LastName: userReq.LastName,
		Email:    userReq.Email,
		Address:  userReq.Address,
		Active:   userReq.Active,
	}
	//Log action
	slog.Info("Creating user ", "newUser", newUser)

	return u.storage.Create(newUser)

}

func (u *UserService) Update(id uuid.UUID, userReq entities.UserRequest) (entities.User, error) {
	newUser := entities.User{
		Id:       id,
		Name:     userReq.Name,
		LastName: userReq.LastName,
		Email:    userReq.Email,
		Address:  userReq.Address,
		Active:   userReq.Active,
	}

	//Log action
	slog.Info("Update user", "newUser", newUser)
	return u.storage.Update(id, newUser)
}

func (u *UserService) Delete(id uuid.UUID) (uuid.UUID, error) {
	slog.Info("Deleting user", "id", id.String())
	return u.storage.Delete(id)
}
