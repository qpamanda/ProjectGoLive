package server

import (
	"ProjectGoLive/database"
	"net/http"
	"strconv"
	"time"

	"github.com/kennygrant/sanitize"
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
	reqID := 0
	repID := 0
	categoryID := 0
	recipientID := 0
	reqStatus := 0
	reqDesc := ""
	toCompleteBy := time.Now()
	address := ""
	createdBy := ""
	createdDT := time.Now()
	lastModifiedBy := ""
	lastModifiedDT := time.Now()
	clientMsg := "" // To display message to the user on the client

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

		address = sanitize.Accents(req.FormValue("address"))

		//fmt.Fprintln(res, viewRecipientSlice)
	}

	request := newRequest{
		reqID,
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

	data := struct {
		details        newRequest
		RecipientSlice []viewRecipient
		ClientMsg      string
	}{
		request,
		viewRecipientSlice,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "addrequest.gohtml", data)
}
