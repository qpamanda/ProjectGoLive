package server

import (
	"ProjectGoLive/database"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Author: Ahmad Bahrudin.
// aaReqSAdd is a handler func that add request status.
func aaReqSAdd(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	myUser, _ := getUser(res, req)
	msg := ""

	if req.Method == http.MethodPost {
		// Process form
		name := req.FormValue("name")

		// check if there is any empty field , return a clientMsg if there is any empty field
		if name == "" {
			msg = "Field cannot empty"
			log.Error("Empty field in creating request status")

			tpl.ExecuteTemplate(res, "aaReqSAdd.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		// call database package and get next request status id from the database
		nextID, err := database.ReqS.GetNextID()

		if err != nil {
			log.Error(err)
			msg = "Internal server error at database"

			tpl.ExecuteTemplate(res, "aaReqSAdd.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		log.WithFields(logrus.Fields{"userName": myUser.UserName}).Infof("[%s] get next request status id: [%s]", myUser.UserName, nextID)

		// call database package and add request status into the database
		err = database.ReqS.Insert(nextID, name, myUser.UserName)

		if err != nil {
			log.Error(err)
			msg = "Internal server error at database"

			tpl.ExecuteTemplate(res, "aaReqSAdd.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		log.WithFields(logrus.Fields{"userName": myUser.UserName}).Infof("[%s] created a new request status: [%s]", myUser.UserName, name)
		msg = "You have successfully created a new request status"

		tpl.ExecuteTemplate(res, "aaReqSAdd.gohtml", AdminHtml{
			User:      myUser,
			ClientMsg: msg})
		return
	}

	tpl.ExecuteTemplate(res, "aaReqSAdd.gohtml", AdminHtml{
		User: myUser})
}

// Author: Ahmad Bahrudin.
// aaReqSUpdate is a handler func that update existing request status.
func aaReqSUpdate(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	myUser, _ := getUser(res, req)
	reqS, _ := database.ReqS.GetAll()
	msg := ""

	if req.Method == http.MethodPost {
		Selection = req.FormValue("selection")
		if Selection == "" {
			msg = "No selection made. Please try again."
			log.Error("No selection made. Please try again.")

			tpl.ExecuteTemplate(res, "aaReqSUpdate.gohtml", AdminHtml{
				User:             myUser,
				ClientMsg:        msg,
				MapRequestStatus: reqS})
			return
		}

		http.Redirect(res, req, "/aaReqSUpdate2", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(res, "aaReqSUpdate.gohtml", AdminHtml{
		User:             myUser,
		MapRequestStatus: reqS})
}

// Author: Ahmad Bahrudin.
// aaReqSUpdate2 is a handler func that update existing request status.
func aaReqSUpdate2(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	myUser, _ := getUser(res, req)
	msg := ""

	if req.Method == http.MethodPost {
		// Process form
		name := req.FormValue("name")

		// check if there is any empty field , return a clientMsg if there is any empty field
		if name == "" {
			msg = "Field cannot empty"
			log.Error("Empty Field in updating request status")

			tpl.ExecuteTemplate(res, "aaReqSUpdate2.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		// call database package and update existing request status into the database
		reqSID, _ := strconv.Atoi(Selection)
		err := database.ReqS.Update(reqSID, name, myUser.UserName)

		if err != nil {
			log.Error(err)
			msg = "Internal server error at database"

			tpl.ExecuteTemplate(res, "aaReqSUpdate2.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		log.WithFields(logrus.Fields{"userName": myUser.UserName}).Infof("[%s] updated existing request status: [%s]", myUser.UserName, name)

		http.Redirect(res, req, "/aaReqSView", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(res, "aaReqSUpdate2.gohtml", AdminHtml{
		User: myUser})
}

// Author: Ahmad Bahrudin.
// aaReqSView is a handler func that display the request status.
func aaReqSView(res http.ResponseWriter, req *http.Request) {
	myUser, _ := getUser(res, req)

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	reqS, _ := database.ReqS.GetAll()

	tpl.ExecuteTemplate(res, "aaReqSView.gohtml", AdminHtml{
		User:             myUser,
		MapRequestStatus: reqS})
}
