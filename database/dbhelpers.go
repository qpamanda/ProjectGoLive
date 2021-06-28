package database

import (
	"ProjectGoLive/recipients"
	"fmt"
	"time"
)

// GetRequestsByStatus implements the sql operations to retrieve all requests by statuscode
// Author: Amanda
func GetRequestsByStatus(currstatus int, helperid int) ([]recipients.Request, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	// Instantiate request
	var (
		requestid       int
		repid           int
		catid           int
		recid           int
		statuscode      int
		reqdesc         string
		tocompleteby    time.Time
		fulfilat        string
		username        string
		firstname       string
		lastname        string
		email           string
		contactno       string
		organisation    string
		category        string
		recname         string
		reccategory     int
		recprofile      string
		reccontactno    string
		status          string
		reccategorydesc string
	)
	requests := make([]recipients.Request, 0)

	query := "SELECT req.RequestID, req.RepID_FK, req.CategoryID, req.RecipientID_FK, " +
		"req.RequestStatusCode, req.RequestDescription, req.ToCompleteBy, req.FulfillAt, " +
		"rep.UserName, rep.FirstName, rep.LastName, rep.Email, rep.ContactNo, rep.Organisation, " +
		"cat.Category, rec.Name, rec.Category AS RecCategory, rec.Profile, rec.ContactNo, " +
		"rq.Status, rc.CategoryDesc " +
		"FROM Requests req " +
		"INNER JOIN Representatives rep ON rep.RepID =  req.RepID_FK " +
		"INNER JOIN Category cat ON cat.CategoryID = req.CategoryID " +
		"INNER JOIN Recipients rec ON rec.RecipientID = req.RecipientID_FK " +
		"INNER JOIN RequestStatus rq ON rq.StatusCode = req.RequestStatusCode " +
		"INNER JOIN RecipientCategory rc ON rc.Category = rec.Category " +
		"WHERE req.RequestStatusCode = ? " +
		"AND req.RepID_FK <> ? "

	results, err := DB.Query(query, currstatus, helperid)
	if err != nil {
		panic("error executing sql select in GetRequestsByStatus")
	} else {
		for results.Next() {
			err := results.Scan(&requestid, &repid, &catid, &recid, &statuscode, &reqdesc,
				&tocompleteby, &fulfilat, &username, &firstname, &lastname, &email,
				&contactno, &organisation, &category, &recname, &reccategory, &recprofile,
				&reccontactno, &status, &reccategorydesc)

			if err != nil {
				panic("error getting results from sql select")
			}

			request := recipients.Request{
				RequestID:       requestid,
				RepID:           repid,
				CatID:           catid,
				RecID:           recid,
				StatusCode:      statuscode,
				ReqDesc:         reqdesc,
				ToCompleteBy:    tocompleteby,
				FulfilAt:        fulfilat,
				UserName:        username,
				FirstName:       firstname,
				LastName:        lastname,
				Email:           email,
				ContactNo:       contactno,
				Organisation:    organisation,
				Category:        category,
				RecName:         recname,
				RecCategory:     reccategory,
				RecProfile:      recprofile,
				RecContactNo:    reccontactno,
				Status:          status,
				RecCategoryDesc: reccategorydesc,
			}
			requests = append(requests, request)
		}
		return requests, nil
	}
}

// GetRequestsToHandle implements the sql operations to retrieve all requests that helper has selected to fulfil.
// Author: Amanda
func GetRequestsToHandle(currstatus int, helperid int) ([]recipients.Request, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	// Instantiate request
	var (
		requestid       int
		repid           int
		catid           int
		recid           int
		statuscode      int
		reqdesc         string
		tocompleteby    time.Time
		fulfilat        string
		username        string
		firstname       string
		lastname        string
		email           string
		contactno       string
		organisation    string
		category        string
		recname         string
		reccategory     int
		recprofile      string
		reccontactno    string
		status          string
		reccategorydesc string
	)
	requests := make([]recipients.Request, 0)

	query := "SELECT req.RequestID, req.RepID_FK, req.CategoryID, req.RecipientID_FK, " +
		"req.RequestStatusCode, req.RequestDescription, req.ToCompleteBy, req.FulfillAt, " +
		"rep.UserName, rep.FirstName, rep.LastName, rep.Email, rep.ContactNo, rep.Organisation, " +
		"cat.Category, rec.Name, rec.Category AS RecCategory, rec.Profile, rec.ContactNo, " +
		"rq.Status, rc.CategoryDesc " +
		"FROM Requests req " +
		"INNER JOIN Representatives rep ON rep.RepID =  req.RepID_FK " +
		"INNER JOIN Category cat ON cat.CategoryID = req.CategoryID " +
		"INNER JOIN Recipients rec ON rec.RecipientID = req.RecipientID_FK " +
		"INNER JOIN RequestStatus rq ON rq.StatusCode = req.RequestStatusCode " +
		"INNER JOIN RecipientCategory rc ON rc.Category = rec.Category " +
		"WHERE req.RequestStatusCode = ? " +
		"AND req.RequestID IN " +
		"(SELECT RequestID FROM Helpers WHERE RepID = ?)"

	results, err := DB.Query(query, currstatus, helperid)
	if err != nil {
		panic("error executing sql select in GetRequestsToHandle")
	} else {
		for results.Next() {
			err := results.Scan(&requestid, &repid, &catid, &recid, &statuscode, &reqdesc,
				&tocompleteby, &fulfilat, &username, &firstname, &lastname, &email,
				&contactno, &organisation, &category, &recname, &reccategory, &recprofile,
				&reccontactno, &status, &reccategorydesc)

			if err != nil {
				panic("error getting results from sql select")
			}

			request := recipients.Request{
				RequestID:       requestid,
				RepID:           repid,
				CatID:           catid,
				RecID:           recid,
				StatusCode:      statuscode,
				ReqDesc:         reqdesc,
				ToCompleteBy:    tocompleteby,
				FulfilAt:        fulfilat,
				UserName:        username,
				FirstName:       firstname,
				LastName:        lastname,
				Email:           email,
				ContactNo:       contactno,
				Organisation:    organisation,
				Category:        category,
				RecName:         recname,
				RecCategory:     reccategory,
				RecProfile:      recprofile,
				RecContactNo:    reccontactno,
				Status:          status,
				RecCategoryDesc: reccategorydesc,
			}
			requests = append(requests, request)
		}
		return requests, nil
	}
}

// UpdateRequestStatus implements the sql operations to update request status from the database.
// Author: Amanda
func UpdateRequestStatus(requestid int, username string, statuscode int) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("UPDATE Requests SET RequestStatusCode=?, " +
		"LastModifiedBy=?, LastModifiedDT=? " +
		"WHERE RequestID=?")
	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(statuscode, username, time.Now(), requestid)
	if err != nil {
		panic("error executing sql update")
	}
	return nil
}

// AddRequestHelper implements the sql operations to insert a requests that helpers selected to fulfil into the database.
// Author: Amanda
func AddRequestHelper(repid int, requestid int, username string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("INSERT INTO Helpers (RepID, RequestID, " +
		"CreatedBy, Created_dt, LastModifiedBy, LastModified_DT) VALUES " +
		"(?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic("error preparing sql insert")
	}

	_, err = stmt.Exec(repid, requestid, username, time.Now(), username, time.Now())
	if err != nil {
		panic("error executing sql insert")
	}
	return nil
}

// DeleteRequestHelper implements the sql operations to delete requests that helper has selected to fulfil from the database.
// Author: Amanda
func DeleteRequestHelper(repid int, requestid int) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	// Do not delete admin record
	stmt, err := DB.Prepare("DELETE FROM Helpers " +
		"WHERE RepID = ? AND RequestID = ?")
	if err != nil {
		panic("error preparing sql delete")
	}

	_, err = stmt.Exec(repid, requestid)
	if err != nil {
		panic("error executing sql delete")
	}
	return nil
}
