package apis

import (
	"database/sql"
	"os"
)

type Api struct {
	DB *sql.DB
}

func New() (*Api, error) {

	// To connect to the database
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	// Create DSN
	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + name

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return &Api{nil}, err
	}

	return &Api{db}, nil
}
