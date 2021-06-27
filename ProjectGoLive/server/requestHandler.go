package server

import (
	"ProjectGoLive/database"
	"net/http"
	"strconv"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/sirupsen/logrus"
)

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

type viewRecipient struct {
	RecipientID int
	Name        string
}

type viewRequest struct {
	RequestID     int
	Category      string
	RecipientName string
	Description   string
	ToCompleteBy  string
}

// Author: Tan Jun Jie
// addrequest is a handler func to create a new request.
// It creates a new request and adds in to the database.
func addrequest(res http.ResponseWriter, req *http.Request) {

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	currentUser, _ := getUser(res, req)

	// Initialize request information
	// Decide whether to add admin entry in Representative table
	// or to add an arbitary non-zero value for admin repID in this method
	// to avoid failing isValidRequest check.
	repID := 0
	categoryID := 0
	recipientID := 0
	reqStatus := 0
	reqDesc := ""
	address := ""
	createdBy := ""
	createdDT := time.Now()
	lastModifiedBy := ""
	lastModifiedDT := time.Now()
	clientMsg := "" // To display message to the user on the client
	var toCompleteBy time.Time

	viewRecipientSlice := make([]viewRecipient, 0)

	repDetails := database.GetRepresentativeDetails(currentUser.UserName)

	// Only 1 key-value pair in repDetails
	for k, v := range repDetails {
		repID = k
		tmpName := v[0] + " " + v[1]
		createdBy = tmpName
		lastModifiedBy = tmpName
	}
	recipients := database.GetRecipientDetails(repID)

	// Parse recipients into viewRecipient format
	for k, v := range recipients {
		viewR := viewRecipient{k, v[0]}
		viewRecipientSlice = append(viewRecipientSlice, viewR)
	}

	// Process the form submission
	if req.Method == http.MethodPost {
		categoryID, _ = strconv.Atoi(req.FormValue("requestcategory"))
		reqDesc = sanitize.Accents(req.FormValue("description"))

		// Convert UTC timestamp to time.Time object set to GMT +8.
		timezoneSuffix := ":00+08:00"
		tmpTime := req.FormValue("tocompletebyDT") + timezoneSuffix
		toCompleteBy, _ = time.Parse(time.RFC3339, tmpTime)

		recipientID, _ = strconv.Atoi(req.FormValue("recipientid"))
		address = sanitize.Accents(req.FormValue("address"))

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

		//fmt.Fprintln(res, isValidRequest(request))

		if isValidRequest(request) {

			err := database.AddRequest(
				repID,
				categoryID,
				recipientID,
				reqStatus,
				reqDesc,
				toCompleteBy,
				address,
				createdBy,
				createdDT,
				lastModifiedBy,
				lastModifiedDT)

			if err != nil {
				clientMsg = "Could not add request"

				log.WithFields(logrus.Fields{
					"repID":     repID,
					"createdBy": createdBy,
					"createdDT": createdDT,
				}).Warn(clientMsg)
			}

			clientMsg = "Request successfully added"

			log.WithFields(logrus.Fields{
				"repID":      repID,
				"createdBy":  createdBy,
				"createdDT":  createdDT,
				"reqDetails": request.RequestDetails,
			}).Info(clientMsg)

		} else {
			clientMsg = "No request added"

			log.WithFields(logrus.Fields{
				"repID":     repID,
				"createdBy": createdBy,
				"createdDT": createdDT,
			}).Error(clientMsg)
		}
	}

	data := struct {
		RecipientSlice []viewRecipient
		ClientMsg      string
	}{
		viewRecipientSlice,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "addrequest.gohtml", data)
}

// Author: Tan Jun Jie
// deleterequest is a handler func to delete an existing request.
// A representative can only choose from the requests he made,
// whereas an admin can choose any request to delete.
func deleterequest(res http.ResponseWriter, req *http.Request) {

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	currentUser, _ := getUser(res, req)

	repID := 0
	clientMsg := "" // To display message to the user on the client
	var isAdmin bool

	repDetails := database.GetRepresentativeDetails(currentUser.UserName)

	// Only 1 key-value pair in repDetails
	for k := range repDetails {
		repID = k
	}

	if currentUser.UserName == "admin" {
		isAdmin = true
		// Decide whether to add admin entry in Representative table
		// or to add an arbitary non-zero value for admin repID
		repID = 5221
	} else {
		isAdmin = false
	}

	viewRequestSlice := make([]viewRequest, 0)

	requests := database.GetRequest(repID, isAdmin)

	// Parse recipients into viewRecipient format
	// Consider writing a function to sort viewRequestSlice
	for k, v := range requests {
		tmpTime := v.ToCompleteBy.Format("Mon, 02 Jan 2006, 15:04")
		viewR := viewRequest{k, convertCategoryID(v.CategoryID), v.RecipientName, v.Description, tmpTime}
		viewRequestSlice = append(viewRequestSlice, viewR)
	}

	//fmt.Fprintln(res, viewRequestSlice)

	// Process the form submission
	if req.Method == http.MethodPost {
		for _, v := range viewRequestSlice {
			selectedRequest := req.FormValue(strconv.Itoa(v.RequestID))

			if selectedRequest != "" {
				if err := database.DeleteRequest(v.RequestID); err == nil {

					clientMsg = "Request successfully deleted"

					log.WithFields(logrus.Fields{
						"repID":       repID,
						"requestID":   v.RequestID,
						"recipient":   v.RecipientName,
						"description": v.Description,
						"deleteDT":    time.Now(),
					}).Info(clientMsg)

				} else {
					clientMsg = "No request deleted"

					log.WithFields(logrus.Fields{
						"repID": repID,
						"time":  time.Now(),
					}).Info(clientMsg)
				}
			}
		}

	}

	data := struct {
		RequestSlice    []viewRequest
		CntCurrentItems int
		ClientMsg       string
	}{
		viewRequestSlice,
		len(viewRequestSlice),
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "deleterequest.gohtml", data)
}

// Author: Tan Jun Jie
// isValidRequest checks that a request has non-empty fields
func isValidRequest(req newRequest) bool {

	// skip RequestCategoryId, RecipientId check
	// as there is a default option in the dropdown of the form
	if req.RepresentativeId == 0 ||
		req.RequestDetails.RequestDescription == "" ||
		req.RequestDetails.FulfilAt == "" ||
		req.CreatedBy == "" ||
		req.LastModifiedBy == "" {
		return false
	}
	return true
}

// Author: Tan Jun Jie
// convertCategoryID returns the string description of a category id
func convertCategoryID(id int) string {
	switch id {
	case 1:
		return "Monetary Donation"
	case 2:
		return "Item Donation"
	case 3:
		return "Errands"
	default:
		return "Invalid Category"
	}
}
