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

// Author: Tan Jun Jie
// GetRepresentativeDetails queries the database for a logged-in representative's ID and name.
func GetRepresentativeDetails(username string) map[int][]string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var repID int
	var firstName string
	var lastName string

	// Instantiate representative details
	var details = make(map[int][]string)

	query := "SELECT RepID, FirstName, LastName FROM Representatives WHERE UserName=?"

	results, err := DB.Query(query, username)
	if err != nil {
		panic("error executing sql select: " + err.Error())
	} else {
		if results.Next() {
			err := results.Scan(&repID, &firstName, &lastName)
			if err != nil {
				panic("error getting results from sql select")
			}
			details[repID] = []string{firstName, lastName}
		}
		return details
	}
}

// Author: Tan Jun Jie
// GetRecipientDetails queries the database for the ID and names of recipients that a representative is in charge of.
func GetRecipientDetails(RepresentativeId int) map[int][]string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var recipientID int
	var name string

	// Instantiate recipients' details
	var details = make(map[int][]string)

	query := "SELECT RecipientID, Name FROM Recipients WHERE RepID_FK=?"

	results, err := DB.Query(query, RepresentativeId)
	if err != nil {
		panic("error executing sql select: " + err.Error())
	} else {
		for results.Next() {
			err := results.Scan(&recipientID, &name)
			if err != nil {
				panic("error getting results from sql select")
			}
			details[recipientID] = []string{name}
		}
		return details
	}
}

/*
// AddCourse implements the sql operations to insert a new course as invoked by the REST API.
func AddXXX(courseID string, courseTitle string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("INSERT INTO Courses (CourseID, CourseTitle, Created_DT, LastModified_DT) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic("error preparing sql insert")
	}

	_, err = stmt.Exec(courseID, courseTitle, time.Now(), time.Now())
	if err != nil {
		panic("error executing sql insert")
	}
}

// UpdateCourse implements the sql operations to update a course as invoked by the REST API.
func UpdateXXX(courseID string, courseTitle string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("UPDATE Courses SET CourseTitle=?, LastModified_DT=? WHERE CourseID=?")
	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(courseTitle, time.Now(), courseID)
	if err != nil {
		panic("error executing sql update")
	}
}

// DeleteCourse implements the sql operations to delete a course as invoked by the REST API.
func DeleteXXX(courseID string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("DELETE FROM Courses WHERE CourseID=?")
	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(courseID)
	if err != nil {
		panic("error executing sql update")
	}
}
*/
