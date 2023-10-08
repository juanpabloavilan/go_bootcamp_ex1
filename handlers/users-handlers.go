package handlers

import (
	"encoding/json"
	"errors"
	"example/bootcamp_ex1/db"
	"example/bootcamp_ex1/entities"
	"example/bootcamp_ex1/services"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// TODO: Put error signatures to all handlers
type RootHandler func(http.ResponseWriter, *http.Request) error

func (h RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		//log error
		slog.Error(err.Error())

		switch err := err.(type) {
		case db.StorageError:
			if errors.Is(err, db.StorageError{Code: db.ErrEntityNotFound}) {
				sendError(w, r, HTTPError{
					Code:        ErrNotFound,
					Status:      http.StatusNotFound,
					Description: "cannot find entity with this id",
				})
			}

			if errors.Is(err, db.StorageError{Code: db.ErrUnmarshalingEntity}) {
				sendError(w, r, HTTPError{
					Code:        ErrBadRequest,
					Status:      http.StatusBadRequest,
					Description: "invalid body",
				})
			}
			if errors.Is(err, db.StorageError{Code: db.ErrMarshalingEntity}) {
				sendError(w, r, HTTPError{
					Code:        ErrBadRequest,
					Status:      http.StatusBadRequest,
					Description: "invalid body",
				})
			}
			if errors.Is(err, db.StorageError{Code: db.ErrGettingRecords}) {
				sendError(w, r, HTTPError{
					Code:        ErrInternalServerError,
					Status:      http.StatusInternalServerError,
					Description: "cannot get records from db",
				})
			}

		case services.ServiceError:
		case HTTPError:
			sendError(w, r, err)
		default:
			sendError(w, r, HTTPError{
				Code:        ErrInternalServerError,
				Status:      http.StatusInternalServerError,
				Description: "internal server error",
			})
		}
	}
}

func GetUserById(userService *services.UserService) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		params := mux.Vars(r)
		idParam := params["id"]
		id, err := uuid.Parse(idParam)

		if err != nil {
			return HTTPError{
				Status:      http.StatusBadRequest,
				Code:        ErrBadRequest,
				Description: "invalid uuid",
			}
		}

		user, err := userService.Get(id)
		if err != nil {
			return err
		}

		userPayload, err := json.Marshal(user)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(userPayload)
		return nil
	}
}

func GetAllUsers(userService *services.UserService) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		users, err := userService.GetAll()
		if err != nil {
			return err
		}
		usersPayload, err := json.Marshal(users)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(usersPayload)
		return nil
	}
}

func CreateUser(userService *services.UserService) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		var newUser entities.UserRequest
		err := json.NewDecoder(r.Body).Decode(&newUser)

		if err != nil {
			return HTTPError{
				Status:      http.StatusBadRequest,
				Code:        ErrBadRequest,
				Description: "invalid user body",
			}
		}

		validate := validator.New()
		err = validate.Struct(newUser)
		if err != nil {
			return HTTPError{
				Code:        ErrBadRequest,
				Status:      http.StatusBadRequest,
				Description: err.Error(),
			}
		}

		id, err := userService.Create(newUser)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id": id.String(),
		})
		return nil

	}
}

func UpdateUser(userService *services.UserService) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			return HTTPError{
				Status:      http.StatusBadRequest,
				Code:        ErrBadRequest,
				Description: "invalid uuid",
			}
		}

		var newUser entities.UserRequest
		err = json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			return HTTPError{
				Status:      http.StatusBadRequest,
				Code:        ErrBadRequest,
				Description: "invalid user body",
			}
		}

		validate := validator.New()
		err = validate.Struct(newUser)
		if err != nil {
			return HTTPError{
				Status:      http.StatusBadRequest,
				Code:        ErrBadRequest,
				Description: err.Error(),
			}
		}

		user, err := userService.Update(id, newUser)

		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		return nil
	}
}

func DeleteUser(userService *services.UserService) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		params := mux.Vars(r)
		id, err := uuid.Parse(params["id"])
		if err != nil {
			return HTTPError{
				Status:      http.StatusBadRequest,
				Code:        ErrBadRequest,
				Description: "invalid uuid",
			}
		}

		id, err = userService.Delete(id)

		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id": id.String(),
		})
		return nil
	}
}

func sendError(w http.ResponseWriter, r *http.Request, err HTTPError) {
	errorPayload, _ := json.Marshal(err)
	w.WriteHeader(err.Status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(errorPayload)
}
