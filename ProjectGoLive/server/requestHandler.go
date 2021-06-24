package server

import (
	"ProjectGoLive/database"
	"net/http"
	"strconv"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/sirupsen/logrus"
)

// addrequest is a handler func to create a new request.
// It creates a new request and adds in to the database.
func addrequest(res http.ResponseWriter, req *http.Request) {

	/*if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	*/
	currentUser := getUser(res, req)

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
				}).Info(clientMsg)
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
			}).Info(clientMsg)
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
