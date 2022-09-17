package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var (
	DBCon *sql.DB
)

func Init() {
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	var dbEndpoint = fmt.Sprintf("%s:%s", dbHost, dbPort)
	var dbSourceName = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbEndpoint, dbName)

	DBCon, err = sql.Open("mysql", dbSourceName)
	if err != nil {
		fmt.Println("Unable to open db connection,", err)
		os.Exit(1)
	}
}
