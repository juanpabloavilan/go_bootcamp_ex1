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
	slog.Info("Getting a user by id", id)
	return u.storage.Get(id)
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
	slog.Info("Creating user: ", newUser)

	id, err := u.storage.Create(newUser)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
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
	slog.Info("Update user: ", newUser)
	return u.storage.Update(id, newUser)
}

func (u *UserService) Delete(id uuid.UUID) (uuid.UUID, error) {
	slog.Info("Deleting user: ", id)
	return u.storage.Delete(id)
}

//métodos create, get, get all, update y delete. Este struct debe ser privado y debe contar con un método constructor.
