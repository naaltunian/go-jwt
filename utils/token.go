package utils

import (
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/naaltunian/go-jwt/models"
)

func GenerateToken(user models.User) (string, error) {
	var err error
	secret := "secret"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "course",
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}
