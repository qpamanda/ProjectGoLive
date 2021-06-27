// Package database implements the connection to the database server at the designated port.
// It performs the DB operations as invoked by the server.
package database

import (
	"database/sql"
	"errors"
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

// Author: Tan Jun Jie
// AddRequest inserts a new request into the database.
func AddRequest(repID, categoryID, recipientID, reqStatus int, reqDesc string, toCompleteBy time.Time, address string, createdBy string, createdDT time.Time, lastModifiedBy string, lastModifiedDT time.Time) (e error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
			switch x := err.(type) {
			case string:
				e = errors.New(x)
			case error:
				e = x
			default:
				e = errors.New("unknown panic")
			}
		}
	}()

	request := newRequest{
		repID,
		categoryID,
		recipientID,
		reqStatus,
		requestDetails{reqDesc, toCompleteBy, address},
		createdBy,
		createdDT,
		lastModifiedBy,
		lastModifiedDT,
	}

	query := `
	INSERT INTO Requests 
		(RepID_FK,
		CategoryID,
		RecipientID_FK,
		RequestStatusCode,
		RequestDescription,
		ToCompleteBy,
		FulfillAt,
		CreatedBy,
		CreatedDT,
		LastModifiedBy,
		LastModifiedDT)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)`

	stmt, err := DB.Prepare(query)
	if err != nil {
		panic("error preparing sql insert: " + err.Error())
	}

	_, err = stmt.Exec(unpackRequest(request))

	if err != nil {
		panic("error executing sql insert: " + err.Error())
	}

	return nil
}

// Author: Tan Jun Jie
// GetRequest gets all requests tied to a representative
// or all requests if one is logged in as an admin.
func GetRequest(repID int, isAdmin bool) map[int]viewRequest {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var requestID int
	var categoryID int
	var toCompleteBy time.Time
	var name string
	var desc string

	// Instantiate requests
	var requests = make(map[int]viewRequest)
	if isAdmin {

		query := `
				SELECT Requests.RequestID, 
				Requests.CategoryID,
				Recipients.Name,
				Requests.RequestDescription,
				Requests.ToCompleteBy
				FROM Requests
				INNER JOIN Recipients
				ON Requests.RecipientID_FK = Recipients.RecipientID;
				`

		results, err := DB.Query(query)

		if err != nil {
			panic("error executing sql select: " + err.Error())
		} else {
			for results.Next() {
				err := results.Scan(&requestID, &categoryID, &name, &desc, &toCompleteBy)
				if err != nil {
					panic("error getting results from sql select: " + err.Error())
				}
				requests[requestID] = viewRequest{categoryID, name, desc, toCompleteBy}
			}
		}
	} else {
		//query := "SELECT RequestID, RequestID, CategoryID, Recipient, RequestDescription FROM Requests WHERE RepID_FK=?"
		query := `
				SELECT Requests.RequestID, 
				Requests.CategoryID,
				Recipients.Name,
				Requests.RequestDescription,
				Requests.ToCompleteBy
				FROM Requests
				INNER JOIN Recipients
				ON Requests.RecipientID_FK = Recipients.RecipientID
				WHERE Requests.RepID_FK=?;
				`
		results, err := DB.Query(query, repID)

		if err != nil {
			panic("error executing sql select: " + err.Error())
		} else {
			for results.Next() {
				err := results.Scan(&requestID, &categoryID, &name, &desc, &toCompleteBy)
				if err != nil {
					panic("error getting results from sql select: " + err.Error())
				}
				requests[requestID] = viewRequest{categoryID, name, desc, toCompleteBy}
			}
		}
	}
	return requests
}

// Author: Tan Jun Jie
// DeleteRequest deletes the request belonging to reqID.
func DeleteRequest(reqID int) (e error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
			switch x := err.(type) {
			case string:
				e = errors.New(x)
			case error:
				e = x
			default:
				e = errors.New("unknown panic")
			}
		}
	}()

	stmt, err := DB.Prepare("DELETE FROM Requests WHERE RequestID=?")
	if err != nil {
		panic("error preparing sql update: " + err.Error())
	}

	_, err = stmt.Exec(reqID)
	if err != nil {
		panic("error executing sql update: " + err.Error())
	}
	return nil
}

func unpackRequest(request newRequest) (repID, categoryID, recipientID, reqStatus int, reqDesc string, toCompleteBy time.Time, address string, createdBy string, createdDT time.Time, lastModifiedBy string, lastModifiedDT time.Time) {
	return request.RepresentativeId,
		request.RequestCategoryId,
		request.RecipientId,
		request.RequestStatus,
		request.RequestDetails.RequestDescription,
		request.RequestDetails.ToCompleteBy,
		request.RequestDetails.FulfilAt,
		request.CreatedBy,
		request.CreatedDT,
		request.LastModifiedBy,
		request.LastModifiedDT
}
