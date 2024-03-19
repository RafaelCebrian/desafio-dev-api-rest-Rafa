package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "testeapi"
	dbname   = "dbApi"
)

func ConnectDB() (*sql.DB, error) {
	dblInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dblInfo)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database %v", err)
	}
	return db, nil
}
