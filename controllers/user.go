package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/naaltunian/go-jwt/driver"
	"github.com/naaltunian/go-jwt/models"
	"github.com/naaltunian/go-jwt/utils"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}

	err = user.ValidateUser()
	if err != nil {
		log.Println(err)
		utils.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Println("error hashing password:", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = string(hashedPassword)

	err = driver.SaveUser(user)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]string{"response": "user created"}
	utils.ResponseWithJSON(w, 201, response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	json.NewDecoder(r.Body).Decode(&user)

	err := user.ValidateUser()
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userFromDB, err := driver.QueryUser(user.Email)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword := userFromDB.Password
	// TODO: move compare password to another function
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		err = errors.New("invalid credentials")
		utils.ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]string{"token": token}
	utils.ResponseWithJSON(w, http.StatusOK, response)
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			err := errors.New("No auth token supplied")
			utils.ResponseWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		authToken := strings.Split(authHeader, " ")[1]
		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err := errors.New("Error parsing token")
				return nil, err
			}
			return []byte("secret"), nil
		})
		if err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			err := errors.New("Invalid token")
			utils.ResponseWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	})
}
