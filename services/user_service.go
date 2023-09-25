package services

import (
	"errors"
	"example/bootcamp_ex1/entities"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("cannot find a user with this id")
)

// Creating initial users using a helper
func createInitialUsers() userMap {
	usersMap := userMap{}
	initialUsers := []entities.User{
		{
			Name:     "Leo",
			LastName: "Messi",
			Email:    "leo.messi@gmail.com",
			Active:   true,
			Address: entities.Address{
				City:          "Miami",
				Country:       "USA",
				AddressString: "Avenida 112f # 88 - 15",
			},
		},
		{
			Name:     "Cristiano",
			LastName: "Ronaldo",
			Email:    "cr7.@gmail.com",
			Active:   true,
			Address: entities.Address{
				City:          "Abu Dhabi",
				Country:       "Emiratos Arabes",
				AddressString: "Cra 112f # 88 - 15",
			},
		},
		{
			Name:     "Kilian",
			LastName: "Mbappe",
			Email:    "mbappe.@gmail.com",
			Active:   true,
			Address: entities.Address{
				City:          "Paris",
				Country:       "Francia",
				AddressString: "Rue 112f # 88 - 15",
			},
		},
		{
			Name:     "Joao",
			LastName: "Felix",
			Email:    "felix.@gmail.com",
			Active:   true,
			Address: entities.Address{
				City:          "Barcelona",
				Country:       "España",
				AddressString: "Calle 112f # 88 - 15",
			},
		},
		{
			Name:     "Robert",
			LastName: "Lewandoski",
			Email:    "lewandoski.@gmail.com",
			Active:   true,
			Address: entities.Address{
				City:          "Barcelona",
				Country:       "España",
				AddressString: "Calle 112f # 88 - 18",
			},
		},
	}
	// Creating uuid for map key and also for user id property
	for _, u := range initialUsers {
		id := uuid.New()
		u.Id = id
		usersMap[id] = u
	}

	return usersMap
}

type userMap map[uuid.UUID]entities.User

type userService struct {
	users userMap
}

func NewUserService() *userService {
	userService := userService{users: createInitialUsers()}
	return &userService
}

func (u *userService) Create(user entities.User) uuid.UUID {
	id := uuid.New()
	user.Id = id
	u.users[id] = user
	return id
}

func (u *userService) Get(key uuid.UUID) (entities.User, error) {
	value, ok := u.users[key]
	//If user doesn't exist we return a nil value and a error
	if !ok {
		return entities.User{}, ErrUserNotFound
	}
	return value, nil
}

func (u *userService) GetAll() []entities.User {
	userList := make([]entities.User, 0, len(u.users))
	for _, user := range u.users {
		userList = append(userList, user)
	}
	return userList
}

func (u *userService) Update(key uuid.UUID, newUser entities.User) (entities.User, error) {
	// If not exists return error
	_, err := u.Get(key)
	if err != nil {
		return entities.User{}, err
	}
	// update
	newUser.Id = key
	u.users[key] = newUser

	return u.users[key], nil
}
func (u *userService) Delete(key uuid.UUID) error {
	// If not exists return error
	_, err := u.Get(key)
	if err != nil {
		return err
	}
	// delete
	delete(u.users, key)
	return nil
}

//métodos create, get, get all, update y delete. Este struct debe ser privado y debe contar con un método constructor.
