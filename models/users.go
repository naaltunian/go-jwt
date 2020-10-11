package models

import (
	"errors"

	"github.com/naaltunian/go-jwt/driver"
)

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

func (u User) SaveUser() error {
	stmt := "insert into users (email, password) values ($1, $2) RETURNING id;"

	err := driver.DB.QueryRow(stmt, u.Email, u.Password).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u User) QueryUser() (User, error) {
	var user User
	stmt := "select * from users where email = $1;"

	row := driver.DB.QueryRow(stmt, u.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}
