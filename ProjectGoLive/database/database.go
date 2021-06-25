// Package database implements the connection to the database server at the designated port.
// It performs the DB operations as invoked by the server.
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Config struct to maintain DB configuration properties
type Config struct {
	ServerName string
	User       string
	Password   string
	DB         string
}

// Connector variable used for DB operations
var DB *sql.DB

// Connect creates the database connection
func Connect(connectionString string) error {
	var err error
	DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	} else {
		fmt.Println("Database opened")
	}
	return nil
}

// GetConnectionString formats the database connection string and returns it.
func GetConnectionString(config Config) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.User, config.Password, config.ServerName, config.DB)
	return connectionString
}
