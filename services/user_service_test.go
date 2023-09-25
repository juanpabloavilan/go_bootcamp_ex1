package services_test

import (
	"example/bootcamp_ex1/entities"
	"example/bootcamp_ex1/services"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// Hacer un package diferente para el test para simular un client
// go test ./... -v

func TestUserService(t *testing.T) {
	t.Run("user service creation returns userservice type", func(t *testing.T) {
		userService := services.NewUserService()
		if reflect.ValueOf(userService).IsZero() {
			t.Errorf("UserService cannot be zeroed ")
		}
	})

	t.Run("user service get all", func(t *testing.T) {
		t.Run("listing all five default users as an array", func(t *testing.T) {
			userService := services.NewUserService()
			userList := userService.GetAll()

			if len(userList) != 5 {
				t.Errorf("expected lenght of 5 but got %d", len(userList))
			}

		})

	})

	t.Run("user service create", func(t *testing.T) {
		t.Run("creating a new user should return uuid ", func(t *testing.T) {
			userService := services.NewUserService()
			newUserId := userService.Create(entities.User{
				Name:     "Neymar",
				LastName: "Junior",
				Email:    "neymar@gmail.com",
				Active:   false,
				Address: entities.Address{
					City:          "Sao Paulo",
					Country:       "Brazil",
					AddressString: "Estadio Macarena",
				},
			})
			got := reflect.TypeOf(newUserId)
			want := reflect.TypeOf(uuid.New())
			if got != want {
				t.Errorf("expected %v return value, but got %v", want, got)
			}

		})

	})

	t.Run("user service get", func(t *testing.T) {
		userService := services.NewUserService()
		userList := userService.GetAll()

		t.Run("getting a user that does exists", func(t *testing.T) {
			id := userList[0].Id
			user, err := userService.Get(id)
			if err != nil {
				t.Errorf("wasn't expecting an error but got %q", err.Error())
			}

			if id != user.Id {
				t.Errorf("expected user id to be %q, but got %q", id, user.Id)
			}
		})

		t.Run("getting a user that don't exist should return an error", func(t *testing.T) {
			id := uuid.New()
			_, err := userService.Get(id)

			if err != services.ErrUserNotFound {
				t.Errorf("wasn't expecting error %q but got %q", services.ErrUserNotFound, err)
			}
		})
	})

	t.Run("user service update", func(t *testing.T) {

		t.Run("updating a user that does exists", func(t *testing.T) {
			userService := services.NewUserService()
			userList := userService.GetAll()
			id := userList[0].Id
			newInfo := entities.User{
				Id:       id,
				Name:     "Michaer",
				LastName: "Neuer",
				Email:    "neuer@gmail.com",
				Active:   false,
				Address: entities.Address{
					City:          "Munich",
					Country:       "Alemania",
					AddressString: "Munich stadium",
				},
			}

			_, err := userService.Update(id, newInfo)

			got, _ := userService.Get(id)

			if err != nil {
				t.Errorf("wasn't expecting an error but got %q", err.Error())
			}

			if !reflect.DeepEqual(got, newInfo) {
				t.Errorf("expected %+v but got %v", got, newInfo)
			}

		})

		t.Run("updating a user that don't exist should return an error", func(t *testing.T) {
			userService := services.NewUserService()
			id := uuid.New()
			newInfo := entities.User{
				Id:       id,
				Name:     "Michaer",
				LastName: "Neuer",
				Email:    "neuer@gmail.com",
				Active:   false,
				Address: entities.Address{
					City:          "Munich",
					Country:       "Alemania",
					AddressString: "Munich stadium",
				},
			}

			_, err := userService.Update(id, newInfo)

			if err != services.ErrUserNotFound {
				t.Errorf("wasn expecting an error %q but got %q", services.ErrUserNotFound.Error(), err.Error())
			}

		})
	})

	t.Run("user service delete", func(t *testing.T) {
		t.Run("deleting a user that does exists", func(t *testing.T) {
			userService := services.NewUserService()
			userList := userService.GetAll()
			id := userList[0].Id

			err := userService.Delete(id)

			if err != nil {
				t.Errorf("wasn't expecting an error but got %v", err)
			}

			_, errGet := userService.Get(id)

			expectedErr := services.ErrUserNotFound
			if errGet != expectedErr {
				t.Errorf("expected %v but got %v", expectedErr, errGet)
			}

		})

		t.Run("deleting a user that don't exist", func(t *testing.T) {
			userService := services.NewUserService()
			id := uuid.New()

			err := userService.Delete(id)

			expectedErr := services.ErrUserNotFound
			if err != expectedErr {
				t.Errorf("expected %v but got %v", expectedErr, err)
			}
		})
	})

}
