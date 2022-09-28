package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "127.0.0.1"
	database = "userCrud"
	user     = "root"
	password = "b4rc3l0s?"
)

var Client *sql.DB

func Connect() (*sql.DB, error) {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	connection, err := sql.Open("mysql", connectionString)
	Client = connection
	return connection, err

}
