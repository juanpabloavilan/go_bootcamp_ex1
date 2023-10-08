package main

import (
	"example/bootcamp_ex1/db"
	"example/bootcamp_ex1/entities"
	"example/bootcamp_ex1/handlers"
	"example/bootcamp_ex1/services"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const (
	ENV_STAGE      = "STAGE"
	ENV_STORAGE    = "STORAGE"
	STORAGE_REDIS  = "REDIS"
	STORAGE_MEMORY = "MEMORY"

	ErrNotValidStorage = "storage is not valid"
)

func main() {
	// Loading env variables
	godotenv.Load()
	slog.Info("ENVIRONMENT: ", os.Getenv(ENV_STAGE), os.Getenv(ENV_STORAGE))

	// Selecting storage from .env
	var storage db.Storage[entities.User]

	switch os.Getenv(ENV_STORAGE) {

	case STORAGE_MEMORY:
		storage = db.NewMemoryStorage[entities.User]()
	case STORAGE_REDIS:
		storage = db.NewRedisStorage[entities.User]()
	default:
		slog.Error(ErrNotValidStorage, os.Getenv(ENV_STORAGE))

	}

	userService := services.NewUserService(storage)

	r := mux.NewRouter()
	// Declaring user subrouter
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Handle("/", handlers.RootHandler(handlers.GetAllUsers(userService))).Methods("GET")
	userRouter.Handle("/{id}", handlers.RootHandler(handlers.GetUserById(userService))).Methods("GET")
	userRouter.Handle("/", handlers.RootHandler(handlers.CreateUser(userService))).Methods("POST")
	userRouter.Handle("/{id}", handlers.RootHandler(handlers.UpdateUser(userService))).Methods("PUT")
	userRouter.Handle("/{id}", handlers.RootHandler(handlers.DeleteUser(userService))).Methods("DELETE")

	// Bind to a port and pass our router in
	slog.Error(http.ListenAndServe(":8000", r).Error())
}
