package entities

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	Name     string
	LastName string
	Email    string
	Active   bool
	Address  Address
}

// func (u User) NewUser(name string, lastname string, email string, active bool, address Address) User {
// 	return User{
// 		Name:     name,
// 		LastName: lastname,
// 		Email:    email,
// 		Active:   active,
// 		Address:  address,
// 	}
// }
