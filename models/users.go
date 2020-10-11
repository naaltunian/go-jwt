package models

import "errors"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u User) ValidateUser() error {

	// TODO: add more email validation
	if u.Email == "" {
		err := errors.New("Invalid email")
		return err
	}

	// TODO: add more password validation
	if u.Password == "" {
		err := errors.New("Invalid password")
		return err
	}
	return nil
}
