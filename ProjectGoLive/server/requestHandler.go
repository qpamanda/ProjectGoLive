package server

import (
	"ProjectGoLive/database"
	"net/http"
	"strconv"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/sirupsen/logrus"
)


var (
	// selecteditrequest passes these to editrequest
	rid         int
	hasSelected bool
)

// representative ID of admin user
const adminID = 5000

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
					}).Error(clientMsg)
				}
			}
		}
		if selectedCnt == 0 {
			clientMsg = "No option selected. No requests deleted."

			log.WithFields(logrus.Fields{
				"repID": repID,
				"time":  time.Now(),
			}).Info(clientMsg)
		}

		for _, id := range idSlice {

			info := requests[id]

			tmpTime := info.ToCompleteBy.Format("Mon, 02 Jan 2006, 15:04")
			// do not display address
			viewR := viewRequest{id, convertCategoryID(info.CategoryID), info.RecipientName, info.Description, tmpTime, ""}
			viewDeletedRequestSlice = append(viewDeletedRequestSlice, viewR)
		}

	}

	data := struct {
		RequestSlice        []viewRequest
		DeletedRequestSlice []viewRequest
		CntCurrentItems     int
		ClientMsg           string
		FormSubmitted       bool
	}{
		viewRequestSlice,
		viewDeletedRequestSlice,
		len(viewRequestSlice),
		clientMsg,
		submitted,
	}

	tpl.ExecuteTemplate(res, "deleterequest.gohtml", data)
}

// Author: Tan Jun Jie
// viewrequest is a handler func to view existing requests.
func viewrequest(res http.ResponseWriter, req *http.Request) {

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
		repID = adminID
	} else {
		isAdmin = false
	}

	viewRequestSlice := make([]viewRequest, 0)

	requests := database.GetRequestByRep(repID, isAdmin)

	for requestid, v := range requests {
		tmpTime := v.ToCompleteBy.Format("Mon, 02 Jan 2006, 15:04")
		viewR := viewRequest{requestid, convertCategoryID(v.CategoryID), v.RecipientName, v.Description, tmpTime, v.FulfillAt}
		viewRequestSlice = append(viewRequestSlice, viewR)
	}

	if len(viewRequestSlice) == 0 {
		clientMsg = "No requests have been made."
	}

	data := struct {
		RequestSlice []viewRequest
		ClientMsg    string
	}{
		viewRequestSlice,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "viewrequest.gohtml", data)

}

// Author: Tan Jun Jie
// selecteditrequest is a handler func to select an existing request to edit.
func selecteditrequest(res http.ResponseWriter, req *http.Request) {

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// resets hasSelected so that users are unable to navigate to /editrequest
	// without going through /selecteditrequest
	hasSelected = false
	// resets rid because user has not selected the request to edit
	rid = 0

	currentUser, _ := getUser(res, req)

	repID := 0
	clientMsg := "" // To display message to the user on the client
	var isAdmin bool

	repDetails := database.GetRepresentativeDetails(currentUser.UserName)

	var createdBy string
	// Only 1 key-value pair in repDetails
	for k := range repDetails {
		repID = k
		createdBy = repDetails[repID][0] + " " + repDetails[repID][1]
	}

	if currentUser.UserName == "admin" {
		isAdmin = true
		// Decide whether to add admin entry in Representative table
		// or to add an arbitary non-zero value for admin repID
		repID = adminID
	} else {
		isAdmin = false
	}

	viewRequestSlice := make([]viewRequest, 0)

	requests := database.GetRequestByRep(repID, isAdmin)

	for requestid, v := range requests {
		tmpTime := v.ToCompleteBy.Format("Mon, 02 Jan 2006, 15:04")
		viewR := viewRequest{requestid, convertCategoryID(v.CategoryID), v.RecipientName, v.Description, tmpTime, v.FulfillAt}
		viewRequestSlice = append(viewRequestSlice, viewR)
	}

	if req.Method == http.MethodPost {
		var requestID int

		selectedRequest := req.FormValue("selection")
		if selectedRequest != "" {
			requestID, _ = strconv.Atoi(selectedRequest)
		}

		if requestID != 0 {
			rid = requestID
			hasSelected = true
			http.Redirect(res, req, "/editrequest", http.StatusSeeOther)
		} else {
			clientMsg = "No selection made. Please try again."

			log.WithFields(logrus.Fields{
				"repID": repID,
				//name of representative
				"createdBy": createdBy,
				"createdDT": time.Now(),
			}).Warn(clientMsg)
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


	tpl.ExecuteTemplate(res, "selecteditrequest.gohtml", data)

}

// Author: Tan Jun Jie
// editrequest is a handler func to edit an existing request.
func editrequest(res http.ResponseWriter, req *http.Request) {

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// redirects users if they have not selected a request to edit
	if !hasSelected {
		http.Redirect(res, req, "/selecteditrequest", http.StatusSeeOther)
	}

	currentUser, _ := getUser(res, req)
	repID := 0
	reqID := rid
	clientMsg := "" // To display message to the user on the client

	repDetails := database.GetRepresentativeDetails(currentUser.UserName)

	var createdBy string
	var submitted bool
	// Only 1 key-value pair in repDetails
	for k := range repDetails {
		repID = k
		createdBy = repDetails[repID][0] + " " + repDetails[repID][1]
	}

	if currentUser.UserName == "admin" {
		repID = adminID
	}

	viewRequestSlice := make([]viewRequest, 0)

	// get request details
	r := database.GetRequest(reqID)

	tmpTime := r.ToCompleteBy.Format("Mon, 02 Jan 2006, 15:04")
	viewR := viewRequest{reqID, convertCategoryID(r.CategoryID), r.RecipientName, r.Description, tmpTime, r.FulfillAt}
	viewRequestSlice = append(viewRequestSlice, viewR)

	if req.Method == http.MethodPost {

		categoryID, _ := strconv.Atoi(req.FormValue("requestcategory"))
		reqDesc := sanitize.Accents(req.FormValue("description"))

		// Convert UTC timestamp to time.Time object set to GMT +8.
		timezoneSuffix := ":00+08:00"
		newTime := req.FormValue("tocompletebyDT")
		tmpTime := newTime + timezoneSuffix
		toCompleteBy, _ := time.Parse(time.RFC3339, tmpTime)

		address := sanitize.Accents(req.FormValue("address"))

		// replace empty fields with existing fields
		if reqDesc == "" {
			reqDesc = viewRequestSlice[0].Description
		}
		if address == "" {
			address = viewRequestSlice[0].FulfillAt
		}
		if newTime == "" {
			toCompleteBy = r.ToCompleteBy
		}

		// disallow user to change toCompleteBy to before current datetime
		if toCompleteBy.Before(time.Now()) {
			clientMsg = "Time indicated has already passed. No request added."

			log.WithFields(logrus.Fields{
				"repID":     repID,
				"createdBy": createdBy,
				"requestID": reqID,
			}).Warn(clientMsg)
		} else {

			if err := database.EditRequest(reqID, categoryID, reqDesc, toCompleteBy, address); err != nil {
				clientMsg = "Could not edit request: "

				log.WithFields(logrus.Fields{
					"repID":     repID,
					"createdBy": createdBy,
					"requestID": reqID,
				}).Error(clientMsg + err.Error())
			} else {
				clientMsg = "Successfully edited request."

				rid = 0
				hasSelected = false
				submitted = true

				tmpTime := toCompleteBy.Format("Mon, 02 Jan 2006, 15:04")
				viewR := viewRequest{reqID, convertCategoryID(categoryID), r.RecipientName, reqDesc, tmpTime, address}
				viewRequestSlice = []viewRequest{viewR}

				log.WithFields(logrus.Fields{
					"repID":     repID,
					"createdBy": createdBy,
					"requestID": reqID,
				}).Info(clientMsg)
			}
		}

	}
	data := struct {
		RequestSlice    []viewRequest
		CntCurrentItems int
		ClientMsg       string
		FormSubmitted   bool
	}{
		viewRequestSlice,
		len(viewRequestSlice),
		clientMsg,
		submitted,
	}

	tpl.ExecuteTemplate(res, "editrequest.gohtml", data)
	
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
