package server

import (
	"ProjectGoLive/database"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Author: Ahmad Bahrudin
// aaCatAdd is a handler func that adds new category.
func aaCatAdd(res http.ResponseWriter, req *http.Request) {
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
			log.Error("Empty Field in creating category")

			tpl.ExecuteTemplate(res, "aaCatAdd.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		// call database package and get next category id from the database
		nextID, err := database.Cat.GetNextID()

		if err != nil {
			log.Error(err)
			msg = "Internal server error at database"

			tpl.ExecuteTemplate(res, "aaCatAdd.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		log.WithFields(logrus.Fields{"userName": myUser.UserName}).Infof("[%s] get next category id: [%s]", myUser.UserName, nextID)

		// call database package and add the category into the database
		err = database.Cat.Insert(nextID, name, myUser.UserName)

		if err != nil {
			log.Error(err)
			msg = "Internal server error at database"

			tpl.ExecuteTemplate(res, "aaCatAdd.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		log.WithFields(logrus.Fields{"userName": myUser.UserName}).Infof("[%s] created a new category: [%s]", myUser.UserName, name)
		msg = "You have successfully created a new category"

		tpl.ExecuteTemplate(res, "aaCatAdd.gohtml", AdminHtml{
			User:      myUser,
			ClientMsg: msg})
		return
	}

	tpl.ExecuteTemplate(res, "aaCatAdd.gohtml", AdminHtml{
		User: myUser})
}

// Author: Ahmad Bahrudin
// aaCatUpdate is a handler func that update existing category.
func aaCatUpdate(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	myUser, _ := getUser(res, req)
	cat, _ := database.Cat.GetAll()
	msg := ""

	if req.Method == http.MethodPost {
		Selection = req.FormValue("selection")
		if Selection == "" {
			msg = "No selection made. Please try again."
			log.Error("No selection made. Please try again.")

			tpl.ExecuteTemplate(res, "aaCatUpdate.gohtml", AdminHtml{
				User:        myUser,
				ClientMsg:   msg,
				MapCategory: cat})
			return
		}

		http.Redirect(res, req, "/aaCatUpdate2", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(res, "aaCatUpdate.gohtml", AdminHtml{
		User:        myUser,
		MapCategory: cat})
}

// Author: Ahmad Bahrudin
// aaCatUpdate2 is a handler func that update existing category.
func aaCatUpdate2(res http.ResponseWriter, req *http.Request) {
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
			log.Error("Empty Field in creating category")

			tpl.ExecuteTemplate(res, "aaCatUpdate2.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		// call database package and update existing category into the database
		catID, _ := strconv.Atoi(Selection)
		err := database.Cat.Update(catID, name, myUser.UserName)

		if err != nil {
			log.Error(err)
			msg = "Internal server error at database"

			tpl.ExecuteTemplate(res, "aaCatUpdate2.gohtml", AdminHtml{
				User:      myUser,
				ClientMsg: msg})
			return
		}

		log.WithFields(logrus.Fields{"userName": myUser.UserName}).Infof("[%s] updated existing category: [%s]", myUser.UserName, name)

		http.Redirect(res, req, "/aaCatView", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(res, "aaCatUpdate2.gohtml", AdminHtml{
		User: myUser})
}

// Author: Ahmad Bahrudin
// aaCatView is a handler func that display the category.
func aaCatView(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	myUser, _ := getUser(res, req)

	cat, _ := database.Cat.GetAll()

	tpl.ExecuteTemplate(res, "aaCatView.gohtml", AdminHtml{
		User:        myUser,
		MapCategory: cat})
}
