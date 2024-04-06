package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func OpenCon() (*sql.DB, error) {
	var err error
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	dbname := os.Getenv("dbname")
	host := os.Getenv("dbhost")
	port := os.Getenv("dbport")
	user := os.Getenv("user")
	password := os.Getenv("password")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable application_name=notesapp",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	//defer db.Close()

	return db, err
}

func CloseCon() {
	db.Close()

}
