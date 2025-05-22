package utils

import (
	"Rodrigo-Guinazu-Arias-TpFinal/models/Sales"
	"Rodrigo-Guinazu-Arias-TpFinal/models/users"
	"fmt"
)

func InitSystem(saleService *Sales.SaleService, userService *users.UserService) ([]*users.User, []*Sales.Sale) {

	var userList []*users.User
	var saleList []*Sales.Sale
	// Create users
	user1 := &users.User{
		Name:     "John Doe",
		Address:  "123 Main St",
		NickName: "johndoe",
	}
	user2 := &users.User{
		Name:     "Jane Smith",
		Address:  "456 Oak Ave",
		NickName: "janesmith",
	}
	user3 := &users.User{
		Name:     "Bob Johnson",
		Address:  "789 Pine Rd",
		NickName: "bobjohnson",
	}

	// Store users and get their IDs
	err := userService.Create(user1)
	if err != nil {
		return nil, nil
	}

	err = userService.Create(user2)
	if err != nil {
		return nil, nil
	}

	err = userService.Create(user3)
	if err != nil {
		return nil, nil
	}

	// Add users to the list
	userList = append(userList, user1, user2, user3)

	// Create sales for each user
	fmt.Println("\nSales created:")
	fmt.Println("-------------")

	fmt.Printf("User: %s, Address: %s, Nickname: %s, UserID: %s\n", user1.Name, user1.Address, user1.NickName, user1.ID)

	fmt.Printf("User: %s, Address: %s, Nickname: %s, UserID: %s\n", user2.Name, user2.Address, user2.NickName, user2.ID)

	fmt.Printf("User: %s, Address: %s, Nickname: %s, UserID: %s\n", user3.Name, user3.Address, user3.NickName, user3.ID)
	// For user1
	sale1, err := saleService.Create(user1.ID, 100.50)
	if err != nil {
		return nil, nil
	}
	fmt.Printf("Sale ID: %s, User: %s, Amount: %.2f, Status: %s\n", sale1.Id, user1.Name, sale1.Amount, sale1.Status)

	sale2, err := saleService.Create(user1.ID, 200.75)
	if err != nil {
		return nil, nil
	}
	fmt.Printf("Sale ID: %s, User: %s, Amount: %.2f, Status: %s\n", sale2.Id, user1.Name, sale2.Amount, sale2.Status)

	// For user2
	sale3, err := saleService.Create(user2.ID, 150.25)
	if err != nil {
		return nil, nil
	}
	fmt.Printf("Sale ID: %s, User: %s, Amount: %.2f, Status: %s\n", sale3.Id, user2.Name, sale3.Amount, sale3.Status)

	sale4, err := saleService.Create(user2.ID, 300.00)
	if err != nil {
		return nil, nil
	}
	fmt.Printf("Sale ID: %s, User: %s, Amount: %.2f, Status: %s\n", sale4.Id, user2.Name, sale4.Amount, sale4.Status)

	// For user3
	sale5, err := saleService.Create(user3.ID, 75.99)
	if err != nil {
		return nil, nil
	}
	fmt.Printf("Sale ID: %s, User: %s, Amount: %.2f, Status: %s\n", sale5.Id, user3.Name, sale5.Amount, sale5.Status)

	sale6, err := saleService.Create(user3.ID, 125.45)
	if err != nil {
		return nil, nil
	}
	fmt.Printf("Sale ID: %s, User: %s, Amount: %.2f, Status: %s\n", sale6.Id, user3.Name, sale6.Amount, sale6.Status)

	// Add sales to the list
	saleList = append(saleList, sale1, sale2, sale3, sale4, sale5, sale6)

	fmt.Println("-------------")

	return userList, saleList
}
