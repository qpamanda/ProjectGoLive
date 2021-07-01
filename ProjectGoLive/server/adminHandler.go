package server

import (
	"net/http"
)

type adminHtml struct {
	User string // To store user name
}

// Author: Ahmad Bahrudin
// categorytable is a handler func that display the category table.
func categorytable(res http.ResponseWriter, req *http.Request) {
	myUser, _ := getUser(res, req)

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "categorytable.gohtml", adminHtml{
		User: myUser.UserName})
}

// Author: Ahmad Bahrudin
// membertypetable is a handler func that display the member type table.
func membertypetable(res http.ResponseWriter, req *http.Request) {
	myUser, _ := getUser(res, req)

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "membertypetable.gohtml", adminHtml{
		User: myUser.UserName})
}

// Author: Ahmad Bahrudin
// requeststatustable is a handler func that display the request status table.
func requeststatustable(res http.ResponseWriter, req *http.Request) {
	myUser, _ := getUser(res, req)

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "requeststatustable.gohtml", adminHtml{
		User: myUser.UserName})
}
