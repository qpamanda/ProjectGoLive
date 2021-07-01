package server

import (
	"ProjectGoLive/authenticate"
	"ProjectGoLive/database"
	"ProjectGoLive/recipients"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Request struct {
	RepID int
}

// selectrequest is a handler func to select request(s) to fulfil
// Redirects to index page if user has not login.
// Author: Amanda
func selectrequest(res http.ResponseWriter, req *http.Request) {
	myUser, _ := getUser(res, req)

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := "" // To display message to the user on the client

	// Get all requests with status - Pending and Requester is not current user himself
	requests, err := database.GetRequestsByStatus(0, myUser.RepID)

	if err != nil {
		clientMsg = "No requests to select"
	}

	// Process the form submission
	if req.Method == http.MethodPost {
		for _, val := range requests {
			// Get requestid from form checkbox using requests.RequestID as its value
			// Since vRequestID is an int, use strconv.Itoa to convert it to a string
			requestid := req.FormValue(strconv.Itoa(val.RequestID))

			// If requestid is not an empty string, the request is selected i.e. checkbox is checked in the form
			// Otherwise, no action required since the request is not selected
			if requestid != "" {
				// Convert RequestID to int
				reqid, _ := strconv.Atoi(requestid)

				// Update request status to 'Being Handled'
				err := database.UpdateRequestStatus(reqid, myUser.UserName, 1)

				if err != nil {
					clientMsg = "ERROR: " + "error updating request status"
					log.WithFields(logrus.Fields{
						"requestid": requestid,
					}).Error(err.Error())
					break
				} else {

					err := database.AddRequestHelper(myUser.RepID, reqid, myUser.UserName)

					if err != nil {
						clientMsg = "ERROR: " + "error adding requests to helper"
						log.WithFields(logrus.Fields{
							"requestid": requestid,
						}).Error(err.Error())
						break
					}
				}
			}
		}
		// Redirect to the view request to fulfil page
		http.Redirect(res, req, "/fulfilrequest", http.StatusSeeOther)
		return
	}

	data := struct {
		User        authenticate.User
		Requests    []recipients.Request
		CntRequests int
		ClientMsg   string
	}{
		myUser,
		requests,
		len(requests),
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "selectrequest.gohtml", data)
}

// fulfilrequest is a handler func to view selected requests to fulfil
// Redirects to index page if user has not login.
// Author: Amanda
func fulfilrequest(res http.ResponseWriter, req *http.Request) {
	myUser, _ := getUser(res, req)

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := "" // To display message to the user on the client

	// Get all requests with status - Being Handled
	requests, err := database.GetRequestsToHandle(1, myUser.RepID)

	if err != nil {
		clientMsg = "No requests selected"
	}

	// Process the form submission
	if req.Method == http.MethodPost {
		buttonClick := req.FormValue("buttonClick")

		bCompleteRequest := false
		if buttonClick == "Requests Completed" {
			bCompleteRequest = true
		}

		for _, val := range requests {
			requestid := req.FormValue(strconv.Itoa(val.RequestID))

			// If requestid is not an empty string, the request is selected i.e. checkbox is checked in the form
			// Otherwise, no action required since the request is not selected
			if requestid != "" {
				// Convert RequestID to int
				reqid, _ := strconv.Atoi(requestid)

				if !bCompleteRequest {
					// Update request status to 'Pending'
					err := database.UpdateRequestStatus(reqid, myUser.UserName, 0)

					if err != nil {
						clientMsg = "ERROR: " + "error updating request status"
						log.WithFields(logrus.Fields{
							"requestid": requestid,
						}).Error(err.Error())
						break
					} else {

						err := database.DeleteRequestHelper(myUser.RepID, reqid)

						if err != nil {
							clientMsg = "ERROR: " + "error deleting requests for helper"
							log.WithFields(logrus.Fields{
								"requestid": requestid,
							}).Error(err.Error())
							break
						}
					}
				} else {
					// Update request status to 'Complete'
					err := database.UpdateRequestStatus(reqid, myUser.UserName, 2)

					if err != nil {
						clientMsg = "ERROR: " + "error updating request status"
						log.WithFields(logrus.Fields{
							"requestid": requestid,
						}).Error(err.Error())
						break
					}
				}
			}
		}

		if !bCompleteRequest {
			// Redirect to the select request(s) to fulfil page
			http.Redirect(res, req, "/selectrequest", http.StatusSeeOther)
			return
		} else {
			// Redirect to the request(s) completed page
			http.Redirect(res, req, "/requestcompleted", http.StatusSeeOther)
			return
		}
	}

	data := struct {
		User        authenticate.User
		Requests    []recipients.Request
		CntRequests int
		ClientMsg   string
	}{
		myUser,
		requests,
		len(requests),
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "fulfilrequest.gohtml", data)
}

// requestcompleted is a handler func to view all completed requests of current user (the helper)
// Redirects to index page if user has not login.
// Author: Amanda
func requestcompleted(res http.ResponseWriter, req *http.Request) {
	myUser, _ := getUser(res, req)

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := "" // To display message to the user on the client

	// Get all requests with status - Completed
	requests, err := database.GetRequestsToHandle(2, myUser.RepID)

	if err != nil {
		clientMsg = "No requests completed"
	}

	data := struct {
		User        authenticate.User
		Requests    []recipients.Request
		CntRequests int
		ClientMsg   string
	}{
		myUser,
		requests,
		len(requests),
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "requestcompleted.gohtml", data)
}
