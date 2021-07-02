package server

import (
	"ProjectGoLive/database"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Author: Ahmad Bahrudin.
// aaMemTypeAdd is a handler func that adds new member type.
func aaMemTypeAdd(res http.ResponseWriter, req *http.Request) {
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
			log.Error("Empty Field in creating member type")

			tpl.ExecuteTemplate(res, "aaMemTypeAdd.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		// call database package and get next member type id from the database
		nextID, err := database.MemT.GetNextID()

		if err != nil {
			log.Error(err)
			msg = "Internal server error at database"

			tpl.ExecuteTemplate(res, "aaMemTypeAdd.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		log.WithFields(logrus.Fields{"userName": myUser.UserName}).Infof("[%s] get next member type id: [%s]", myUser.UserName, nextID)

		// call database package and add the member type into the database
		err = database.MemT.Insert(nextID, name, myUser.UserName)

		if err != nil {
			log.Error(err)
			msg = "Internal server error at database"

			tpl.ExecuteTemplate(res, "aaMemTypeAdd.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		log.WithFields(logrus.Fields{"userName": myUser.UserName}).Infof("[%s] created a new member type: [%s]", myUser.UserName, name)
		msg = "You have successfully created a new member type"

		tpl.ExecuteTemplate(res, "aaMemTypeAdd.gohtml", AdminHtml{
			User:      myUser,
			ClientMsg: msg})
		return
	}

	tpl.ExecuteTemplate(res, "aaMemTypeAdd.gohtml", AdminHtml{
		User: myUser})
}

// Author: Ahmad Bahrudin.
// aaMemTypeUpdate is a handler func that update existing member type.
func aaMemTypeUpdate(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	myUser, _ := getUser(res, req)
	memT, _ := database.MemT.GetAll()
	msg := ""

	if req.Method == http.MethodPost {
		Selection = req.FormValue("selection")
		if Selection == "" {
			msg = "No selection made. Please try again."
			log.Error("No selection made. Please try again.")

			tpl.ExecuteTemplate(res, "aaMemTypeUpdate.gohtml", AdminHtml{
				User:          myUser,
				ClientMsg:     msg,
				MapMemberType: memT})
			return
		}

		http.Redirect(res, req, "/aaMemTypeUpdate2", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(res, "aaMemTypeUpdate.gohtml", AdminHtml{
		User:          myUser,
		MapMemberType: memT})
}

// Author: Ahmad Bahrudin.
// aaMemTypeUpdate2 is a handler func that update existing member type.
func aaMemTypeUpdate2(res http.ResponseWriter, req *http.Request) {
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
			log.Error("Empty Field in creating member type")

			tpl.ExecuteTemplate(res, "aaMemTypeUpdate2.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		// call database package and update existing member type into the database
		memTID, _ := strconv.Atoi(Selection)
		err := database.MemT.Update(memTID, name, myUser.UserName)

		if err != nil {
			log.Error(err)
			msg = "Internal server error at database"

			tpl.ExecuteTemplate(res, "aaMemTypeUpdate2.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		log.WithFields(logrus.Fields{"userName": myUser.UserName}).Infof("[%s] updated existing member type: [%s]", myUser.UserName, name)

		http.Redirect(res, req, "/aaMemTypeView", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(res, "aaMemTypeUpdate2.gohtml", AdminHtml{
		User: myUser})
}

// Author: Ahmad Bahrudin.
// aaMemTypeView is a handler func that display the member type.
func aaMemTypeView(res http.ResponseWriter, req *http.Request) {
	myUser, _ := getUser(res, req)

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	memT, _ := database.MemT.GetAll()

	tpl.ExecuteTemplate(res, "aaMemTypeView.gohtml", AdminHtml{
		User:          myUser,
		MapMemberType: memT})
}
