package services_test

import (
	"errors"
	"example/bootcamp_ex1/db"
	"example/bootcamp_ex1/entities"
	"example/bootcamp_ex1/services"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestUserServiceGetAll(t *testing.T) {
	tt := []struct {
		description   string
		db            db.Storage[entities.User]
		expected      []entities.User
		expectedError error
		assert        func(t *testing.T, output []entities.User, expectedOutput []entities.User, err error, expectedError error)
	}{
		{
			description: "list users succesfully",
			db: RedisStorageMock[entities.User]{
				GetAllEntities: func() ([]entities.User, error) {
					return usersMapListMock, nil
				},
			},
			expected: usersMapListMock,
			assert: func(t *testing.T, output, expectedOutput []entities.User, err, expectedError error) {
				if err != nil {
					t.Errorf("wasn't expecting an error but got %q", err)
				}
				if !reflect.DeepEqual(expectedOutput, output) {
					t.Errorf("wanted %v but got %v", expectedOutput, output)
				}
			},
		},
		{
			description: "returns err ErrGettingRecords from database",
			db: RedisStorageMock[entities.User]{
				GetAllEntities: func() ([]entities.User, error) {
					return nil, db.StorageError{
						Code: db.ErrGettingRecords,
					}
				},
			},
			expectedError: db.StorageError{
				Code: db.ErrGettingRecords,
			},
			assert: func(t *testing.T, output, expectedOutput []entities.User, err, expectedError error) {
				if !errors.Is(err, expectedError) {
					t.Errorf("error expected %q but got %q", expectedError, err)
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			userService := services.NewUserService(tc.db)
			result, err := userService.GetAll()
			tc.assert(t, result, tc.expected, err, tc.expectedError)
		})
	}

}

func TestUserServiceGet(t *testing.T) {
	tt := []struct {
		description   string
		db            db.Storage[entities.User]
		expected      entities.User
		expectedError error
		assert        func(t *testing.T, output, expectedOutput entities.User, err error, expectedError error)
	}{
		{
			description: "get a user succesfully",
			db: RedisStorageMock[entities.User]{
				GetEntity: func(id uuid.UUID) (entities.User, error) {
					return userGetMock, nil
				},
			},
			expected: userGetMock,
			assert: func(t *testing.T, output, expectedOutput entities.User, err, expectedError error) {
				if err != nil {
					t.Errorf("wasn't expecting an error but got %q", err)
				}
				if !reflect.DeepEqual(expectedOutput, output) {
					t.Errorf("wanted %v but got %v", expectedOutput, output)
				}
			},
		},
		{
			description: "should return error not found entity",
			db: RedisStorageMock[entities.User]{
				GetEntity: func(id uuid.UUID) (entities.User, error) {
					return entities.User{}, db.StorageError{
						Code: db.ErrEntityNotFound,
					}
				},
			},
			expectedError: db.StorageError{Code: db.ErrEntityNotFound},
			assert: func(t *testing.T, output, expectedOutput entities.User, err, expectedError error) {
				if !errors.Is(err, expectedError) {
					t.Errorf("error expected %q but got %q", expectedError, err)
				}
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			userService := services.NewUserService(tc.db)
			result, err := userService.Get(uuid.New())
			tc.assert(t, result, tc.expected, err, tc.expectedError)
		})
	}

}

func TestUserServiceCreate(t *testing.T) {
	tt := []struct {
		description   string
		db            db.Storage[entities.User]
		input         entities.UserRequest
		expected      uuid.UUID
		expectedError error
		assert        func(t *testing.T, input entities.UserRequest, output uuid.UUID, expectedOutput uuid.UUID, err error, expectedError error)
	}{
		{
			description: "should return new id and no error",
			db: RedisStorageMock[entities.User]{
				CreateEntity: func(thing entities.User) (uuid.UUID, error) {
					return thing.Id, nil
				},
			},
			input: userRequestMock,
			assert: func(t *testing.T, input entities.UserRequest, output uuid.UUID, expectedOutput uuid.UUID, err, expectedError error) {
				if err != nil {
					t.Errorf("wasn't expecting an error but got %q", err)
				}
				if reflect.TypeOf(output) != reflect.TypeOf(uuid.UUID{}) {
					t.Errorf("fields in userRequest input and user output are different")
				}
			},
		},
		{
			description: "should return error marshaling entity from db",
			db: RedisStorageMock[entities.User]{
				CreateEntity: func(thing entities.User) (uuid.UUID, error) {
					return uuid.Nil, db.StorageError{
						Code: db.ErrMarshalingEntity,
					}
				},
			},
			input: userRequestMock,
			expectedError: db.StorageError{
				Code: db.ErrMarshalingEntity,
			},
			assert: func(t *testing.T, input entities.UserRequest, output uuid.UUID, expectedOutput uuid.UUID, err, expectedError error) {
				if !errors.Is(err, expectedError) {
					t.Errorf("error expected %q but got %q", expectedError, err)
				}

			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			userService := services.NewUserService(tc.db)
			result, err := userService.Create(tc.input)
			tc.assert(t, tc.input, result, tc.expected, err, tc.expectedError)
		})
	}

}

func TestUserServiceUpdate(t *testing.T) {
	tt := []struct {
		description   string
		db            db.Storage[entities.User]
		inputId       uuid.UUID
		inputUser     entities.UserRequest
		expected      entities.User
		expectedError error
		assert        func(t *testing.T, inputId uuid.UUID, inputUser entities.UserRequest, output entities.User, expectedOutput entities.User, err error, expectedError error)
	}{
		{description: "should return new user and no error",
			db: RedisStorageMock[entities.User]{
				UpdateEntity: func(id uuid.UUID, thing entities.User) (entities.User, error) {
					return thing, nil
				},
			},
			inputId:   uuid.New(),
			inputUser: userRequestMock,
			expected:  userMockUpdate,
			assert: func(t *testing.T, inputId uuid.UUID, inputUser entities.UserRequest, output entities.User, expectedOutput entities.User, err error, expectedError error) {
				if err != nil {
					t.Errorf("wasn't expecting an error but got %q", err)
				}
				if !reflect.DeepEqual(output.Address, expectedOutput.Address) {
					t.Errorf("wanted %+v but got %+v", expectedOutput, output)
				}
				if output.Name != expectedOutput.Name {
					t.Errorf("wanted %+v but got %+v", expectedOutput.Name, output.Name)
				}
				if output.LastName != expectedOutput.LastName {
					t.Errorf("wanted %+v but got %+v", expectedOutput.LastName, output.LastName)
				}
				if output.Email != expectedOutput.Email {
					t.Errorf("wanted %+v but got %+v", expectedOutput.Email, output.Email)
				}
				if output.Active != expectedOutput.Active {
					t.Errorf("wanted %+v but got %+v", expectedOutput.Active, output.Active)
				}

			},
		},
		{
			description: "should return error marshaling entity from db",
			db: RedisStorageMock[entities.User]{
				UpdateEntity: func(id uuid.UUID, thing entities.User) (entities.User, error) {
					return entities.User{}, db.StorageError{
						Code: db.ErrMarshalingEntity,
					}
				},
			},
			inputId:   uuid.New(),
			inputUser: userRequestMock,
			expectedError: db.StorageError{
				Code: db.ErrMarshalingEntity,
			},
			assert: func(t *testing.T, inputId uuid.UUID, inputUser entities.UserRequest, output entities.User, expectedOutput entities.User, err error, expectedError error) {
				if !errors.Is(err, expectedError) {
					t.Errorf("error expected %q but got %q", expectedError, err)
				}

			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			userService := services.NewUserService(tc.db)
			result, err := userService.Update(tc.inputId, tc.inputUser)
			tc.assert(t, tc.inputId, tc.inputUser, result, tc.expected, err, tc.expectedError)
		})
	}
}

func TestUserServiceDelete(t *testing.T) {
	tt := []struct {
		description   string
		db            db.Storage[entities.User]
		input         uuid.UUID
		expected      uuid.UUID
		expectedError error
		assert        func(t *testing.T, output, expectedOutput uuid.UUID, err error, expectedError error)
	}{
		{
			description: "delete a user succesfully should return same id from db",
			input:       uuid.MustParse("e49e654d-8a3b-48f3-bdab-a58eb1c79e22"),
			db: RedisStorageMock[entities.User]{
				DeleteEntity: func(id uuid.UUID) (uuid.UUID, error) {
					return id, nil
				},
			},
			expected: uuid.MustParse("e49e654d-8a3b-48f3-bdab-a58eb1c79e22"),
			assert: func(t *testing.T, output, expectedOutput uuid.UUID, err, expectedError error) {
				if err != nil {
					t.Errorf("wasn't expecting an error but got %q", err)
				}
				if !reflect.DeepEqual(expectedOutput, output) {
					t.Errorf("wanted %v but got %v", expectedOutput, output)
				}
			},
		},
		{
			description: "should return error not found entity",
			db: RedisStorageMock[entities.User]{
				DeleteEntity: func(id uuid.UUID) (uuid.UUID, error) {
					return uuid.Nil, db.StorageError{
						Code: db.ErrEntityNotFound,
					}
				},
			},
			expectedError: db.StorageError{Code: db.ErrEntityNotFound},
			assert: func(t *testing.T, output, expectedOutput uuid.UUID, err, expectedError error) {
				if !errors.Is(err, expectedError) {
					t.Errorf("error expected %q but got %q", expectedError.Error(), err.Error())
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			userService := services.NewUserService(tc.db)
			result, err := userService.Delete(tc.input)
			tc.assert(t, result, tc.expected, err, tc.expectedError)
		})
	}

}

type funcGet[T entities.StorageObject] func(id uuid.UUID) (T, error)
type funcCreate[T entities.StorageObject] func(thing T) (uuid.UUID, error)
type funcGetAll[T entities.StorageObject] func() ([]T, error)
type funcUpdate[T entities.StorageObject] func(id uuid.UUID, thing T) (T, error)
type funcDelete[T entities.StorageObject] func(id uuid.UUID) (uuid.UUID, error)

type RedisStorageMock[T entities.StorageObject] struct {
	Entities       map[uuid.UUID]T
	GetEntity      funcGet[T]
	GetAllEntities funcGetAll[T]
	CreateEntity   funcCreate[T]
	UpdateEntity   funcUpdate[T]
	DeleteEntity   funcDelete[T]
}

func (rm RedisStorageMock[T]) Get(id uuid.UUID) (T, error) {
	return rm.GetEntity(id)
}

func (rm RedisStorageMock[T]) GetAll() ([]T, error) {
	return rm.GetAllEntities()
}

func (rm RedisStorageMock[T]) Create(thing T) (uuid.UUID, error) {
	return rm.CreateEntity(thing)
}

func (rm RedisStorageMock[T]) Update(id uuid.UUID, thing T) (T, error) {
	return rm.UpdateEntity(id, thing)
}

func (rm RedisStorageMock[T]) Delete(id uuid.UUID) (uuid.UUID, error) {
	return rm.DeleteEntity(id)
}

var usersMapListMock []entities.User = []entities.User{
	{
		Id:       uuid.MustParse("a1ff79d6-df55-469b-8386-99cf245dda77"),
		Name:     "Luke",
		LastName: "Skywalker",
		Email:    "luke@starwars.com",
		Active:   false,
		Address: entities.Address{
			Country:       "Tatooine",
			City:          "Somewhere in Tatooine",
			AddressString: "Some address in Tatooine",
		},
	},
}

var userGetMock entities.User = entities.User{
	Name:     "Juan Pablo",
	LastName: "Avilan",
	Email:    "juan.avilan@gmail.com",
	Active:   true,
	Address: entities.Address{
		Country:       "Tatooine",
		City:          "Somewhere in Tatooine",
		AddressString: "Some address in Tatooine",
	},
}

var userRequestMock entities.UserRequest = entities.UserRequest{
	Name:     "Anakin",
	LastName: "Skywalker",
	Email:    "anakin@starwars.com",
	Active:   false,
	Address: entities.Address{
		Country:       "Tatooine",
		City:          "Somewhere in Tatooine",
		AddressString: "Some address in Tatooine",
	},
}
var userMockUpdate entities.User = entities.User{
	Name:     "Anakin",
	LastName: "Skywalker",
	Email:    "anakin@starwars.com",
	Active:   false,
	Address: entities.Address{
		Country:       "Tatooine",
		City:          "Somewhere in Tatooine",
		AddressString: "Some address in Tatooine",
	},
}
