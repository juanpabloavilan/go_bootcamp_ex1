package main

import (
	"example/bootcamp_ex1/db"
	"example/bootcamp_ex1/handlers"
	"example/bootcamp_ex1/services"
	"log"
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
	var storage db.Storage

	switch os.Getenv(ENV_STORAGE) {

	case STORAGE_MEMORY:
		storage = db.NewMemoryStorage()
	case STORAGE_REDIS:
		storage = db.NewRedisStorage()
	default:
		slog.Error(ErrNotValidStorage, os.Getenv(ENV_STORAGE))

	}

	userService := services.NewUserService(storage)

	r := mux.NewRouter()
	// Declaring user subrouter
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/", handlers.GetAllUsers(userService)).Methods("GET")
	userRouter.HandleFunc("/{id}", handlers.GetUserById(userService)).Methods("GET")
	userRouter.HandleFunc("/", handlers.CreateUser(userService)).Methods("POST")
	userRouter.HandleFunc("/{id}", handlers.UpdateUser(userService)).Methods("PUT")
	userRouter.HandleFunc("/{id}", handlers.DeleteUser(userService)).Methods("DELETE")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
