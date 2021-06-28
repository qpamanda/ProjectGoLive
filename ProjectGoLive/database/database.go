// Package database implements the connection to the database server at the designated port.
// It performs the DB operations as invoked by the server.
package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config struct to maintain DB configuration properties
type Config struct {
	ServerName string
	User       string
	Password   string
	DB         string
}

// req struct for storing request information
type newRequest struct {
	RepresentativeId int // id of the coordinator/representative
	/*
		RequestCategoryId
		1 (monetary donation)
		2 (item donation)
		3 (errands)
	*/
	RequestCategoryId int
	RecipientId       int // id of recipient who receives the aid
	/*
		RequestStatus
		0 (pending/waiting to be matched to a helper)
		1 (being handled)
		2 (completed)
	*/
	RequestStatus  int
	RequestDetails requestDetails
	CreatedBy      string
	CreatedDT      time.Time
	LastModifiedBy string
	LastModifiedDT time.Time
}

//requestDetails struct for storing request detail information
type requestDetails struct {
	RequestDescription string
	ToCompleteBy       time.Time
	FulfilAt           string
}

type viewRequest struct {
	CategoryID    int
	RecipientName string
	Description   string
	ToCompleteBy  time.Time
	FulfillAt     string
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
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", config.User, config.Password, config.ServerName, config.DB)
	return connectionString
}
