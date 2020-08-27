package main

import (
	"context"
	"github.com/volatiletech/authboss/v3"
	"log"
)

// MemStorer stores users in memory
type MemStorer struct {
	Users  map[string]User
	Tokens map[string][]string
}

// User struct for authboss
type User struct {
	ID int
	// Non-authboss related field
	Name string
	Role string
	// Auth
	Email    string
	Password string
}

func NewMemStorer() *MemStorer {
	return &MemStorer{
		Users: map[string]User{
			"ginkt95@gmail.com": {
				ID:                 1,
				Name:               "Vlad",
				Role: 				"admin",
				Password:           "$2a$10$XtW/BrS5HeYIuOCXYe8DFuInetDMdaarMUJEOg/VA/JAIDgw3l4aG", // pass = 1234
				Email:              "ginkt95@gmail.com",
			},
			"john95@gmail.com": {
				ID:                 2,
				Name:               "John",
				Role: 				"user",
				Password:           "$2a$10$XtW/BrS5HeYIuOCXYe8DFuInetDMdaarMUJEOg/VA/JAIDgw3l4aG", // pass = 1234
				Email:              "john95@gmail.com",
			},
		},
		Tokens: make(map[string][]string),
	}
}

// GetPID from user
func (u User) GetPID() string { return u.Email }

// GetPassword from user
func (u User) GetPassword() string { return u.Password }

// GetRole from user
func (u User) GetRole() string { return u.Role }

// PutPID into user
func (u *User) PutPID(pid string) { u.Email = pid }

// PutPassword into user
func (u *User) PutPassword(password string) { u.Password = password }

// PutRole into user
func (u *User) PutRole(role string) { u.Role = role }

// Load the user
func (m MemStorer) Load(_ context.Context, key string) (user authboss.User, err error) {
	u, ok := m.Users[key]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	log.Println("Loaded user:", u.Name, "| Role:", u.Role)
	return &u, nil
}

// Save the user
func (m MemStorer) Save(_ context.Context, user authboss.User) error {
	u := user.(*User)
	m.Users[u.Email] = *u

	log.Println("Saved user:", u.Name, "| Role:", u.Role)
	return nil
}
