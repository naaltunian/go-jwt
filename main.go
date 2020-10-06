package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

func GenerateToken(user User) (string, error) {
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

func validateUser(user User) error {

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

func signup(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}

	err = validateUser(user)
	if err != nil {
		log.Println(err)
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Println("error hashing password:", err)
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = string(hashedPassword)

	err = saveUser(user)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]string{"response": "user created"}
	responseWithJSON(w, 201, response)
}

func login(w http.ResponseWriter, r *http.Request) {
	var user User

	json.NewDecoder(r.Body).Decode(&user)

	err := validateUser(user)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userFromDB, err := queryUser(user.Email)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword := userFromDB.Password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		err = errors.New("invalid credentials")
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := GenerateToken(user)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]string{"token": token}
	responseWithJSON(w, http.StatusOK, response)
}

func queryUser(email string) (User, error) {
	var user User
	stmt := "select * from users where email = $1;"
	// password := user.Password

	row := db.QueryRow(stmt, email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("protected")
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			err := errors.New("No auth token supplied")
			responseWithError(w, http.StatusBadRequest, err.Error())
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
			responseWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			err := errors.New("Invalid token")
			responseWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	})
}

func saveUser(user User) error {
	stmt := "insert into users (email, password) values ($1, $2) RETURNING id;"

	err := db.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return err
	}

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
