/*
Author: Huang Yanping
Last Updated: 26-Jun-2021
*/
package server

import (
	"ProjectGoLive/authenticate"
	"ProjectGoLive/database"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Recipient struct {
	RecipientID    int
	RepID          int
	Name           string
	Category       bool
	Profile        string
	ContactNo      string
	CreatedDT      time.Time
	LastModifiedDT time.Time
}

func manageRecipient(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	clientMsg := ""
	clientMsg2 := ""

	// get the user information using session cookie
	myUser, _ := getUser(w, r)
	authenticate.IsAdmin = database.IsAdmin(myUser.UserName)

	// Get all recipients from the database under this representative in a form of slice
	recipients, err := database.GetMyRecipient(myUser.RepID)

	if err != nil {
		clientMsg = "Internal server error at database"
	}

	// check if there is any recipient in the database under this RepID
	if len(recipients) == 0 {
		clientMsg = "You currently have no recipients"
	} else {
		// sort the slice from database in alphabetical order according to the recipient name
		sort.SliceStable(recipients, func(i, j int) bool { return recipients[i].Name < recipients[j].Name })
	}

	data := struct {
		User       authenticate.User
		IsAdmin    bool
		Recipients []database.Recipient
		ClientMsg  string
		ClientMsg2 string
	}{
		myUser,
		authenticate.IsAdmin,
		recipients,
		clientMsg,
		clientMsg2,
	}
	tpl.ExecuteTemplate(w, "manageRecipient.gohtml", data)
}

// Add recipient to the database
func addRecipient(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	clientMsg := ""

	myUser, _ := getUser(w, r)
	authenticate.IsAdmin = database.IsAdmin(myUser.UserName)

	if r.Method == http.MethodPost {
		// Process form
		name := r.FormValue("name")
		category := r.FormValue("category")
		profile := r.FormValue("profile")
		contactNo := r.FormValue("contact")

		// check if there is any empty field
		if name == "" || profile == "" || contactNo == "" {
			clientMsg = "Field cannot empty"
			data := struct {
				User      authenticate.User
				IsAdmin   bool
				ClientMsg string
			}{
				myUser,
				authenticate.IsAdmin,
				clientMsg,
			}
			tpl.ExecuteTemplate(w, "addRecipient.gohtml", data)
			return
		}

		var categoryBool bool

		if category == "Individual" {
			categoryBool = true
		}

		err := database.AddRecipient(myUser.RepID, name, categoryBool, profile, contactNo)

		if err != nil {
			fmt.Println(err)
			clientMsg = "Internal server error at database"
		} else {
			clientMsg = "You have successfully created a new recipient"
		}
	}
	data := struct {
		User      authenticate.User
		IsAdmin   bool
		ClientMsg string
	}{
		myUser,
		authenticate.IsAdmin,
		clientMsg,
	}
	tpl.ExecuteTemplate(w, "addRecipient.gohtml", data)
}
func getRecipient(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""
	myUser, _ := getUser(w, r)
	authenticate.IsAdmin = database.IsAdmin(myUser.UserName)

	// get the user input from URL query
	inputs := r.URL.Query()["recipientID"]

	// get the recipientID from the URL query
	recipientID := inputs[0]

	// convert the recipientID from string to int64
	recipientIDInt, err := strconv.ParseInt(recipientID, 10, 0)

	// checking for error when converting the string to int64
	if err != nil {
		clientMsg = "Internal server error"
	}

	// query the database and get the recipient information
	recipient, err := database.GetRecipient(myUser.RepID, recipientIDInt)

	// checking for error at the database query
	if err != nil {
		clientMsg = "Internal server error at database"
	}

	data := struct {
		User      authenticate.User
		IsAdmin   bool
		Recipient database.Recipient
		ClientMsg string
	}{
		myUser,
		authenticate.IsAdmin,
		recipient,
		clientMsg,
	}

	tpl.ExecuteTemplate(w, "getRecipient.gohtml", data)
}

func updateRecipient(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""
	myUser, _ := getUser(w, r)
	authenticate.IsAdmin = database.IsAdmin(myUser.UserName)
	// get the user input from URL query
	inputs := r.URL.Query()["recipientID"]

	// get the recipientID from the URL query
	recipientID := inputs[0]

	// convert the recipientID from string to int64
	recipientIDInt, err := strconv.ParseInt(recipientID, 10, 0)

	// checking for error when converting the string to int64
	if err != nil {
		clientMsg = "Internal server error"
	}

	if r.Method == http.MethodPost {
		// Process form
		name := r.FormValue("name")
		category := r.FormValue("category")
		profile := r.FormValue("profile")
		contactNo := r.FormValue("contact")

		// check if there is any empty field
		if name == "" || profile == "" || contactNo == "" {
			clientMsg = "Field cannot empty"
			recipient, _ := database.GetRecipient(myUser.RepID, recipientIDInt)
			data := struct {
				User      authenticate.User
				IsAdmin   bool
				Recipient database.Recipient
				ClientMsg string
			}{
				myUser,
				authenticate.IsAdmin,
				recipient,
				clientMsg,
			}
			tpl.ExecuteTemplate(w, "updateRecipient.gohtml", data)
			return
		}

		var categoryBool bool

		if category == "Individual" {
			categoryBool = true
		}

		err := database.UpdateRecipient(myUser.RepID, recipientIDInt, name, categoryBool, profile, contactNo)

		if err != nil {
			clientMsg = "Internal server error at database"
		} else {
			clientMsg = "You have successfully updated " + name + "."
			recipient, err := database.GetRecipient(myUser.RepID, recipientIDInt)
			if err != nil {
				clientMsg = "Internal server error at database"
			}
			data := struct {
				User      authenticate.User
				IsAdmin   bool
				Recipient database.Recipient
				ClientMsg string
			}{
				myUser,
				authenticate.IsAdmin,
				recipient,
				clientMsg,
			}
			tpl.ExecuteTemplate(w, "getRecipient.gohtml", data)
			return
		}
	}

	// query the database and get the recipient information
	recipient, err := database.GetRecipient(myUser.RepID, recipientIDInt)

	// checking for error at the database query
	if err != nil {
		clientMsg = "Internal server error at database"
	}

	data := struct {
		User      authenticate.User
		IsAdmin   bool
		Recipient database.Recipient
		ClientMsg string
	}{
		myUser,
		authenticate.IsAdmin,
		recipient,
		clientMsg,
	}
	tpl.ExecuteTemplate(w, "updateRecipient.gohtml", data)
}

func deleteRecipient(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	clientMsg := ""
	clientMsg2 := ""
	myUser, _ := getUser(w, r)
	authenticate.IsAdmin = database.IsAdmin(myUser.UserName)

	// get the user input from URL query
	inputs := r.URL.Query()["recipientID"]

	// get the recipientID from the URL query
	recipientID := inputs[0]

	// convert the recipientID from string to int64
	recipientIDInt, err := strconv.ParseInt(recipientID, 10, 0)

	// checking for error when converting the string to int64
	if err != nil {
		clientMsg2 = "Internal server error"
	}

	// query the database and get the recipient information
	err2 := database.DeleteRecipient(myUser.RepID, recipientIDInt)

	// checking for error at the database query
	if err2 != nil {
		clientMsg2 = "Internal server error at database"
	} else {
		clientMsg2 = "You have successfully deleted a recipient"
	}
	recipients, err3 := database.GetMyRecipient(myUser.RepID)

	if err3 != nil {
		clientMsg2 = "Internal server error at database"
	}

	if len(recipients) == 0 {
		clientMsg = "You currently have no recipients"
	} else {
		// sort the slice from database in alphabetical order according to the recipient name
		sort.SliceStable(recipients, func(i, j int) bool { return recipients[i].Name < recipients[j].Name })
	}

	data := struct {
		User       authenticate.User
		IsAdmin    bool
		Recipients []database.Recipient
		ClientMsg  string
		ClientMsg2 string
	}{
		myUser,
		authenticate.IsAdmin,
		recipients,
		clientMsg,
		clientMsg2,
	}
	tpl.ExecuteTemplate(w, "manageRecipient.gohtml", data)
}
