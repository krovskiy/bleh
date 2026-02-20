package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Username string
	Password []byte
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// implement db!
var users = []User{
	{ID: 1, Username: "cat", Password: mustHash("meow")},
}

func mustHash(password string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return hash
}

func findUser(username string) *User {
	for i := range users {
		if users[i].Username == username {
			return &users[i]
		}
	}
	return nil
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		var newLoginRequest LoginRequest
		err := json.NewDecoder(r.Body).Decode(&newLoginRequest)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
		user := findUser(newLoginRequest.Username)
		if user == nil {
			http.Error(w, "wrong", http.StatusUnauthorized)
			return
		}
		err = bcrypt.CompareHashAndPassword(user.Password, []byte(newLoginRequest.Password))
		if err != nil {
			http.Error(w, "wrong", http.StatusUnauthorized)
			return
		}
	}
}
