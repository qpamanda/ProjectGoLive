package server

import (
	"net/http"
	"strconv"
	"time"
)

// addrequest is a handler func to create a new request.
// It creates a new request and adds in to the database.
func addrequest(res http.ResponseWriter, req *http.Request) {

	/*if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	*/

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

	// Process the form submission
	if req.Method == http.MethodPost {
		categoryID, _ = strconv.Atoi(req.FormValue("requestcategory"))
		reqDesc = req.FormValue("description")

		// Convert UTC timestamp to time.Time object set to GMT +8.
		timezoneSuffix := ":00+08:00"
		tmpTime := req.FormValue("tocompletebyDT") + timezoneSuffix
		toCompleteBy, _ = time.Parse(time.RFC3339, tmpTime)

		address = req.FormValue("address")
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
		details   newRequest
		ClientMsg string
	}{
		request,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "addrequest.gohtml", data)
}
