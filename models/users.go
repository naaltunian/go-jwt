package models

import (
	"errors"
	"regexp"
	"unicode"

	"github.com/naaltunian/go-jwt/driver"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u User) ValidateUser() error {

	err := validateEmail(u.Email)
	if err != nil {
		return err
	}

	err = validatePassword(u.Password)
	if err != nil {
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

func validateEmail(email string) error {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	isValid := emailRegex.MatchString(email)

	if !isValid {
		err := errors.New("invalid email format")
		return err
	}
	return nil
}

func validatePassword(password string) error {
	var (
		hasUpper       = false
		hasLower       = false
		hasNumber      = false
		hasSpecialChar = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}

	if len(password) < 8 || !hasUpper || !hasLower || !hasNumber || !hasSpecialChar {
		err := errors.New("Password must have at least 8 characters, at least one uppercase character, a number, and a special character")
		return err
	}
	return nil
}
