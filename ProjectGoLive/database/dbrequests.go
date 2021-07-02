package database

import (
	"errors"
	"fmt"
	"time"
)

// Author: Tan Jun Jie.
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
		results.Close()
		return details
	}
}

// Author: Tan Jun Jie.
// GetRecipientDetails queries the database for the ID and names of recipients that a representative is in charge of.
func GetRecipientDetails(RepresentativeId int, isAdmin bool) map[int][]string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var recipientID int
	var name string

	// Instantiate recipients' details
	var details = make(map[int][]string)

	if !isAdmin {

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
			results.Close()
		}
	} else {
		query := "SELECT RecipientID, Name FROM Recipients"

		results, err := DB.Query(query)
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
			results.Close()
		}
	}
	return details
}

// Author: Tan Jun Jie.
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
				e = errors.New("Unknown panic")
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

	stmt.Close()
	return nil
}

// Author: Tan Jun Jie.
// GetRequestByRep gets all requests tied to a representative
// or all requests if one is logged in as an admin.
func GetRequestByRep(repID int, isAdmin bool) map[int]viewRequest {
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
	var addr string
	var status int

	// Instantiate requests
	var requests = make(map[int]viewRequest)
	if isAdmin {

		query := `
				SELECT Requests.RequestID, 
				Requests.CategoryID,
				Recipients.Name,
				Requests.RequestDescription,
				Requests.ToCompleteBy,
				Requests.FulfillAt,
				Requests.RequestStatusCode
				FROM Requests
				INNER JOIN Recipients
				ON Requests.RecipientID_FK = Recipients.RecipientID;
				`

		results, err := DB.Query(query)

		if err != nil {
			panic("error executing sql select: " + err.Error())
		} else {
			for results.Next() {
				err := results.Scan(&requestID, &categoryID, &name, &desc, &toCompleteBy, &addr, &status)
				if err != nil {
					panic("error getting results from sql select: " + err.Error())
				}
				requests[requestID] = viewRequest{categoryID, name, desc, toCompleteBy, addr, status}
			}
			results.Close()
		}
	} else {
		//query := "SELECT RequestID, RequestID, CategoryID, Recipient, RequestDescription FROM Requests WHERE RepID_FK=?"
		query := `
				SELECT Requests.RequestID, 
				Requests.CategoryID,
				Recipients.Name,
				Requests.RequestDescription,
				Requests.ToCompleteBy,
				Requests.FulfillAt,
				Requests.RequestStatusCode
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
				err := results.Scan(&requestID, &categoryID, &name, &desc, &toCompleteBy, &addr, &status)
				if err != nil {
					panic("error getting results from sql select: " + err.Error())
				}
				requests[requestID] = viewRequest{categoryID, name, desc, toCompleteBy, addr, status}
			}
			results.Close()
		}
	}
	return requests
}

// Author: Tan Jun Jie.
// GetRequest gets the request details of an existing request.
func GetRequest(reqID int) viewRequest {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var categoryID int
	var toCompleteBy time.Time
	var name string
	var desc string
	var addr string
	var status int

	// Instantiate request
	var request = viewRequest{}
	query := `
				SELECT Requests.CategoryID,
				Recipients.Name,
				Requests.RequestDescription,
				Requests.ToCompleteBy,
				Requests.FulfillAt,
				Requests.RequestStatusCode
				FROM Requests
				INNER JOIN Recipients
				ON Requests.RecipientID_FK = Recipients.RecipientID
				WHERE Requests.RequestID=?;
				`
	results, err := DB.Query(query, reqID)

	if err != nil {
		panic("error executing sql select: " + err.Error())
	} else {
		for results.Next() {
			err := results.Scan(&categoryID, &name, &desc, &toCompleteBy, &addr, &status)
			if err != nil {
				panic("error getting results from sql select: " + err.Error())
			}
			request = viewRequest{categoryID, name, desc, toCompleteBy, addr, status}
		}
	}
	results.Close()
	return request
}

// Author: Tan Jun Jie.
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
				e = errors.New("Unknown panic")
			}
		}
	}()

	stmt, err := DB.Prepare("DELETE FROM Requests WHERE RequestID=?")
	if err != nil {
		panic("error preparing sql delete: " + err.Error())
	}

	_, err = stmt.Exec(reqID)
	if err != nil {
		panic("error executing sql delete: " + err.Error())
	}
	stmt.Close()
	return nil
}

// Author: Tan Jun Jie.
// EditRequest edits an existing request.
func EditRequest(requestid, categoryid int, reqDesc string, toCompleteBy time.Time, address string) (e error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
			switch x := err.(type) {
			case string:
				e = errors.New(x)
			case error:
				e = x
			default:
				e = errors.New("Unknown panic")
			}
		}
	}()

	stmt, err := DB.Prepare("UPDATE Requests SET CategoryID=?, RequestDescription=?, ToCompleteBy=?,FulfillAt=? WHERE RequestID=?")
	if err != nil {
		panic("error preparing sql update: " + err.Error())
	}

	_, err = stmt.Exec(categoryid, reqDesc, toCompleteBy, address, requestid)
	if err != nil {
		panic("error executing sql update: " + err.Error())
	}
	stmt.Close()
	return nil
}

// Author: Tan Jun Jie.
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
