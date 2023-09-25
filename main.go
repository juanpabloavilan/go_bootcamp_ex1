package main

import (
	"example/bootcamp_ex1/entities"
	"example/bootcamp_ex1/services"
	"fmt"
)

func main() {
	//Creating a user service
	userService := services.NewUserService()

	userList := userService.GetAll()
	fmt.Printf("\nPrinting all users \n %+v \n\n", userList)

	//Getting one user ID
	userId := userList[0].Id

	user, _ := userService.Get(userId)
	fmt.Printf("\nGet a user by id: %q \n  %+v \n\n", userId, user)

	user, _ = userService.Update(userId, entities.User{
		Name:     "Michaer",
		LastName: "Neuer",
		Email:    "neuer@gmail.com",
		Active:   false,
		Address: entities.Address{
			City:          "Munich",
			Country:       "Alemania",
			AddressString: "Munich stadium",
		},
	})
	fmt.Printf("\nUpdate a user \n %+v \n\n", user)

	fmt.Printf("\n Printing all users\n %+v \n\n", userService.GetAll())

	userService.Delete(userId)
	fmt.Printf("\nDelete a user %q \n\n", userId)

	fmt.Printf("Getting all users\n %+v \n\n", userService.GetAll())

}
