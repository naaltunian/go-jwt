package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/naaltunian/go-jwt/controllers"
	"github.com/naaltunian/go-jwt/driver"
	"github.com/naaltunian/go-jwt/models"
)

func main() {

	driver.ConnectToDB()
	defer driver.DB.Close()

	router := mux.NewRouter()

	router.HandleFunc("/signup", controllers.Signup).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/protected", controllers.TokenVerifyMiddleWare(controllers.ProtectedEndpoint)).Methods("GET")

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func validateUser(user models.User) error {

	// TODO: add more email validation
	if user.Email == "" {
		err := errors.New("Invalid email")
		return err
	}

	// TODO: add more password validation
	if user.Password == "" {
		err := errors.New("Invalid password")
		return err
	}
	return nil
}
