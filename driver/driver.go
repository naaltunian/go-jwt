package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     string
	port     string
	user     string
	password string
	dbname   string
)

var DB *sql.DB

func init() {
	hostVar, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		log.Fatal("Requires a DB host environment variable")
	}
	host = hostVar

	portVar, ok := os.LookupEnv("POSTGRES_DB_PORT")
	if !ok {
		log.Fatal("Requires a DB port environment variable")
	}
	port = portVar

	userVar, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		log.Fatal("Requires a DB user environment variable")
	}
	user = userVar

	dbNameVar, ok := os.LookupEnv("POSTGRES_DB_NAME")
	if !ok {
		log.Fatal("Requires a DB name environment variable")
	}
	dbname = dbNameVar

	passwordVar, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatal("Requires a DB host environment variable")
	}
	password = passwordVar
}

func ConnectToDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
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
	DB = db
}
