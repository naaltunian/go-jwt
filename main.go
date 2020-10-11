package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/naaltunian/go-jwt/controllers"
	"github.com/naaltunian/go-jwt/driver"
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
