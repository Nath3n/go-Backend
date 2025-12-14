package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

//Defining a user struct

type User struct {
	ID    int
	Name  string
	Email string
}

//funtion to greet a user

func greetUser(u User) {
	fmt.Printf("Hello %s! your email is %s\n", u.Name, u.Email)
}

// function to update a user's information
func updateUser(users []User, id int, newName string, newEmail string) bool {
	userPtr := getUserByID(users, id)
	if userPtr != nil {
		userPtr.Name = newName
		userPtr.Email = newEmail
		return true
	}
	return false
}

// function to add a new user
func addUser(users []User, newUser User) []User {
	return append(users, newUser)
}

// function to get a user by ID
func getUserByID(users []User, id int) *User {
	for i := range users {
		if users[i].ID == id {
			return &users[i]
		}
	}
	return nil
}

// function to delete a user by ID
func deleteUser(users []User, id int) []User {
	for i, u := range users {
		if u.ID == id {
			return append(users[:i], users[i+1:]...)
		}

	}
	return users
}

// get all users handler
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// get single user handler
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userPtr := getUserByID(users, id)
	if userPtr == nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(userPtr)
}

var users []User
var mu sync.Mutex

// Main function to demonstrate the user management functionalities
func main() {

	//using imported packages to avoid unused import errors
	_ = http.MethodGet
	_, _ = json.Marshal([]int{}) // uses json
	user := User{
		ID:    1,
		Name:  "Nathan",
		Email: "nathan@email.com",
	}
	users = []User{
		{ID: 1, Name: "Joshua", Email: "Joshua@email.com"},
		{ID: 2, Name: "Duke", Email: "Duke@email.com"},
	}

	//Adding a new user
	newUser := User{ID: 3, Name: "Alice", Email: "alice@wonderland.com"}
	users = addUser(users, newUser)

	//Greeting all users from the range
	for _, u := range users {
		greetUser(u)
	}
	//Getting a user by ID and greeting them
	userPtr := getUserByID(users, 2)

	//checking if userPtr is not nil before dereferencing
	if userPtr != nil {
		greetUser(*userPtr)
	} else {
		fmt.Println("User not found")
	}

	//updating a user and greeting all the users again
	updated := updateUser(users, 2, "Duke Updated", "duke@newgmail.com")

	//checking if the user was updated successfully
	if updated {
		fmt.Println("User has been updated successfully.")
	} else {
		fmt.Println("User not found ")
	}

	//greeting all users after the update
	for _, u := range users {
		greetUser(u)
	}

	//deleting a user
	users = deleteUser(users, 1)
	fmt.Println("After deletion:")

	//greeting all users again after the deletion
	for _, u := range users {
		greetUser(u)
		fmt.Println(user)
		fmt.Println(user.Name)

		greetUser(user)
	}
	http.HandleFunc("/users", getUsersHandler)
	http.HandleFunc("/user", getUserHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
