package database

import "errors"

// Author: Ahmad Bahrudin
// requestStatus struct for storing request status information
type requestStatus struct {
	StatusCode      int    // To store requestStatus's id
	Status          string // To store requestStatus's name
	CreatedBy       string // To store requestStatus's created by
	Created_dt      string // To store requestStatus's created date/time
	LastModifiedBy  string // To store requestStatus's last modified by
	LastModified_dt string // To store requestStatus's last modified date/time
}

var (
	ReqS = &requestStatus{0, "", "", "", "", ""}
)

func (*requestStatus) Insert(statusCode int, status string, userName string) error {
	stmt, err := DB.Prepare("INSERT INTO RequestStatus " +
		"(StatusCode, Status, CreatedBy, Created_dt, LastModifiedBy, LastModified_dt) " +
		"VALUES (?, ?, ?, now(), ?, now())")

	if err != nil {
		panic("error preparing sql insert")
	}

	_, err = stmt.Exec(statusCode, status, userName, userName)
	if err != nil {
		panic("error executing sql insert")
	}

	return nil
}

func (*requestStatus) Update(statusCode int, status string, userName string) error {
	stmt, err := DB.Prepare("UPDATE RequestStatus " +
		"SET StatusCode=?, Status=?, LastModifiedBy=?, LastModified_dt=now() " +
		"WHERE StatusCode=?")

	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(statusCode, status, userName, statusCode)
	if err != nil {
		panic("error executing sql update")
	}

	return nil
}

func (reqS requestStatus) Get(statusCode int) (requestStatus, error) {
	query := "SELECT * " +
		"FROM RequestStatus " +
		"WHERE StatusCode=?"

	results, err := DB.Query(query, statusCode)
	if err != nil {
		panic("error executing sql select")
	}

	if results.Next() {
		err := results.Scan(&reqS.StatusCode, &reqS.Status, &reqS.CreatedBy, &reqS.Created_dt, &reqS.LastModifiedBy, &reqS.LastModified_dt)

		if err != nil {
			panic("error getting results from sql select")
		}
	}
	return reqS, errors.New("member type not found")
}

func (reqS requestStatus) GetAll() (map[int]requestStatus, error) {
	query := "SELECT * " +
		"FROM RequestStatus "

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	}

	m := make(map[int]requestStatus)
	for results.Next() {
		err := results.Scan(&reqS.StatusCode, &reqS.Status, &reqS.CreatedBy, &reqS.Created_dt, &reqS.LastModifiedBy, &reqS.LastModified_dt)

		if err != nil {
			panic("error getting results from sql select")
		}

		m[reqS.StatusCode] = reqS
	}
	return m, nil
}

func (reqS requestStatus) Delete(categoryID int) error {
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
