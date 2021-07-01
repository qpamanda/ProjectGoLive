package database

import "strings"

// Author: Ahmad Bahrudin
// RequestStatus struct for storing request status information
type RequestStatus struct {
	StatusCode      int    // To store request status' id
	Status          string // To store request status' name
	CreatedBy       string // To store request status' created by
	Created_dt      string // To store request status' created date/time
	LastModifiedBy  string // To store request status' last modified by
	LastModified_dt string // To store request status' last modified date/time
}

var (
	ReqS = &RequestStatus{0, "", "", "", "", ""}
)

// Author: Ahmad Bahrudin
// GetNextID function that get next request status' id from the database
func (reqS RequestStatus) GetNextID() (int, error) {
	query := "SELECT StatusCode " +
		"FROM RequestStatus"

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	}

	for results.Next() {
		err := results.Scan(&reqS.StatusCode)

		if err != nil {
			panic("error getting results from sql select")
		}
	}
	return (reqS.StatusCode + 1), nil
}

// Author: Ahmad Bahrudin
// Insert function that adds new request status to the database
func (*RequestStatus) Insert(statusCode int, status string, userName string) error {
	stmt, err := DB.Prepare("INSERT INTO RequestStatus " +
		"(StatusCode, Status, CreatedBy, Created_dt, LastModifiedBy, LastModified_dt) " +
		"VALUES (?, ?, ?, DATE_ADD(NOW(), INTERVAL 8 HOUR), ?, DATE_ADD(NOW(), INTERVAL 8 HOUR))")

	if err != nil {
		panic("error preparing sql insert")
	}

	_, err = stmt.Exec(statusCode, status, userName, userName)
	if err != nil {
		panic("error executing sql insert")
	}

	return nil
}

// Author: Ahmad Bahrudin
// Update function that update exist request status from the database
func (*RequestStatus) Update(statusCode int, status string, userName string) error {
	stmt, err := DB.Prepare("UPDATE RequestStatus " +
		"SET Status=?, LastModifiedBy=?, LastModified_dt=DATE_ADD(NOW(), INTERVAL 8 HOUR) " +
		"WHERE StatusCode=?")

	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(status, userName, statusCode)
	if err != nil {
		panic("error executing sql update")
	}

	return nil
}

// Author: Ahmad Bahrudin
// GetAll function that get all request status from the database
func (reqS RequestStatus) GetAll() (map[int]RequestStatus, error) {
	query := "SELECT * " +
		"FROM RequestStatus "

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	}

	m := make(map[int]RequestStatus)
	for results.Next() {
		err := results.Scan(&reqS.StatusCode, &reqS.Status, &reqS.CreatedBy, &reqS.Created_dt, &reqS.LastModifiedBy, &reqS.LastModified_dt)

		if err != nil {
			panic("error getting results from sql select")
		}
		reqS.Created_dt = strings.Replace(reqS.Created_dt, "T", " ", -1)
		reqS.Created_dt = strings.Replace(reqS.Created_dt, "Z", " ", -1)
		reqS.LastModified_dt = strings.Replace(reqS.LastModified_dt, "T", " ", -1)
		reqS.LastModified_dt = strings.Replace(reqS.LastModified_dt, "Z", " ", -1)

		m[reqS.StatusCode] = reqS
	}
	return m, nil
}

// Author: Ahmad Bahrudin
// Delete function that delete request status from the database
// IMPORTANT: ONLY USE FOR TESTING
func (reqS RequestStatus) Delete(categoryID int) error {
	stmt, err := DB.Prepare("DELETE FROM RequestStatus " +
		"WHERE StatusCode=?")

	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(categoryID)
	if err != nil {
		panic("error executing sql update")
	}

	return nil
}
