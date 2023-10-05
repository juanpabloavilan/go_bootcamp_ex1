package handlers

import (
	"encoding/json"
	"example/bootcamp_ex1/entities"
	"example/bootcamp_ex1/services"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetUserById(userService *services.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		idParam := params["id"]
		id, err := uuid.Parse(idParam)

		if err != nil {
			sendError(w, r, "Invalid id", http.StatusBadRequest, err.Error())
			return
		}

		user, err := userService.Get(id)

		if err != nil {
			sendError(w, r, "User not found with this id", http.StatusNotFound, err.Error())
			return
		}

		userPayload, err := json.Marshal(user)
		if err != nil {
			sendError(w, r, "There was an error", http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(userPayload)

	}
}

func GetAllUsers(userService *services.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := userService.GetAll()
		if err != nil {
			sendError(w, r, "There was an error", http.StatusInternalServerError, err.Error())
			return
		}
		usersPayload, err := json.Marshal(users)
		if err != nil {
			sendError(w, r, "There was an error", http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(usersPayload)

	}
}

func CreateUser(userService *services.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser entities.UserRequest
		err := json.NewDecoder(r.Body).Decode(&newUser)

		if err != nil {
			sendError(w, r, "Unvalid body", http.StatusBadRequest, err.Error())
			return
		}

		validate := validator.New()
		err = validate.Struct(newUser)
		if err != nil {
			sendError(w, r, "Unvalid body", http.StatusBadRequest, err.Error())
			return
		}

		id, err := userService.Create(newUser)
		if err != nil {
			sendError(w, r, "Unvalid body", http.StatusBadRequest, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id": id.String(),
		})

	}
}

func UpdateUser(userService *services.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			sendError(w, r, "Invalid id", http.StatusBadRequest, err.Error())
			return
		}

		var newUser entities.UserRequest
		err = json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			sendError(w, r, "Unvalid body", http.StatusBadRequest, err.Error())
			return
		}

		validate := validator.New()
		err = validate.Struct(newUser)
		if err != nil {
			sendError(w, r, "Unvalid body", http.StatusBadRequest, err.Error())
			return
		}

		user, err := userService.Update(id, newUser)

		if err != nil {
			sendError(w, r, "Error", http.StatusNotFound, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func DeleteUser(userService *services.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			sendError(w, r, "Invalid id", http.StatusBadRequest, err.Error())
			return
		}

		id, err = userService.Delete(id)

		if err != nil {
			sendError(w, r, "Error", http.StatusNotFound, err.Error())
			return

		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id": id.String(),
		})
	}
}

func sendError(w http.ResponseWriter, r *http.Request, msg string, statusCode int, errorDetails string) {
	slog.Error(errorDetails)
	err := struct {
		Code         int
		Message      string
		ErrorDetails string
	}{
		Code:    statusCode,
		Message: msg,
	}

	if os.Getenv("STAGE") == "development" {
		err.ErrorDetails = errorDetails
	}

	errorPayload, _ := json.Marshal(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(errorPayload)
}
