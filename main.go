package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "go-jwt"
)

var db *sql.DB

func main() {

	db = connectToDB()
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/signup", signup).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/protected", TokenVerifyMiddleWare(protectedEndpoint)).Methods("GET")

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func connectToDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("Error connecting to db:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging DB", err)
	}

	log.Println("Connected to DB")
	return db
}

func validateUser(user User) error {

	// TODO: add email validation
	if user.Email == "" {
		err := errors.New("Invalid email")
		return err
	}

	// TODO: add password validation
	if user.Password == "" {
		err := errors.New("Invalid password")
		return err
	}
	return nil
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}

	log.Println(user)
	err = validateUser(user)
	if err != nil {
		log.Println(err)
		responseWithError(w, 400, err.Error())
		return
	}

	response := map[string]string{"response": "user created"}
	responseWithJSON(w, 201, response)
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

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling json", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func responseWithError(w http.ResponseWriter, code int, message string) {
	responseWithJSON(w, code, map[string]string{"error": message})
}
