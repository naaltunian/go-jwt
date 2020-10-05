package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/signup", signup).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/protected", TokenVerifyMiddleWare(protectedEndpoint)).Methods("GET")

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func signup(w http.ResponseWriter, r *http.Request) {
	log.Println("signup")
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("login")
}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("protected")
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return nil
}
