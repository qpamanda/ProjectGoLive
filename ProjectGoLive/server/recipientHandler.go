package server

import (
	"ProjectGoLive/authenticate"
	"ProjectGoLive/database"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// struct for storing JSON file from numverify API
type Result struct {
	Valid                bool
	Number               string
	Local_format         string
	International_format string
	Country_prefix       string
	Country_code         string
	Country_name         string
	Location             string
	Carrier              string
	Line_type            string
}

// struct for storing recipient details
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

// Author: Huang Yanping.
// manageRecipient a handler func which will get all recipients under the user account
// from the database.
func manageRecipient(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// client msg to be displayed to the user indicating the result of the function
	clientMsg := ""
	clientMsg2 := ""

	// get the user information using session cookie
	myUser, _ := getUser(w, r)
	// check if the user is a admin
	authenticate.IsAdmin = database.IsAdmin(myUser.UserName)

	// Get all recipients from the database under this representative in a form of slice
	recipients, err := database.GetMyRecipient(myUser.RepID)

	if err != nil {
		clientMsg = "Internal server error at database"
		log.Error(err)
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

// Author: Huang Yanping.
// addRecipient is a handler func which add recipient to the database
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

		// check if there is any empty field , return a clientMsg if there is any empty field
		if name == "" || profile == "" || contactNo == "" {
			clientMsg = "Field cannot empty"
			log.Error("Empty Field in creating recipient")
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

		// validating the contact number
		validContact, err := validateContactNo(contactNo)

		// return clientMsg if there is any error in using the external API
		if err != nil {
			log.Error(err)
			log.Error("Problem with numverify API")
			clientMsg = "Server error"
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

		// return clientMsg if it is a invalid contact number
		if !validContact {
			log.Error("Invalid contact number at creating recipient")
			clientMsg = "Invalid contact number, please only enter the 8-digit numbers"
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

		} else { // executing adding of recipient once all field is validated

			var categoryBool bool

			// check for category input and converting it from string to bool
			if category == "Individual" {
				categoryBool = true
			}

			// call database package and add the recipient into the database
			err := database.AddRecipient(myUser.RepID, name, categoryBool, profile, contactNo)

			if err != nil {
				log.Error(err)
				clientMsg = "Internal server error at database"
			} else {
				log.WithFields(logrus.Fields{
					"userName": myUser.UserName,
				}).Infof("[%s] created a new recipient: [%s]", myUser.UserName, name)
				clientMsg = "You have successfully created a new recipient"
			}
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

// Author: Huang Yanping.
// getRecipient is a handler func that gets the recipient details
// from the database using the recipientID.
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
		log.Error(err)
		clientMsg = "Internal server error"
	}

	// query the database and get the recipient information
	recipient, err := database.GetRecipient(myUser.RepID, recipientIDInt)

	// checking for error at the database query
	if err != nil {
		log.Error(err)
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

// Author: Huang Yanping.
// updateRecipient is a handler func that update the recipient details
// in the database using the recipientID.
func updateRecipient(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""
	myUser, _ := getUser(w, r)
	authenticate.IsAdmin = database.IsAdmin(myUser.UserName)

	inputs := r.URL.Query()["recipientID"]

	recipientID := inputs[0]

	recipientIDInt, err := strconv.ParseInt(recipientID, 10, 0)

	if err != nil {
		log.Error(err)
		clientMsg = "Internal server error, update unsuccessful"
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		category := r.FormValue("category")
		profile := r.FormValue("profile")
		contactNo := r.FormValue("contact")

		if name == "" || profile == "" || contactNo == "" {
			log.Error("Empty Field in updating recipient")
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

		validContact, err := validateContactNo(contactNo)

		if err != nil {
			log.Error(err)
			log.Error("Problem with numverify API")
			clientMsg = "Server error"
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

		if !validContact {
			log.Error("Invalid contact number at updating recipient")
			clientMsg = "Invalid contact number, please only enter the 8-digit numbers"
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
		} else {

			var categoryBool bool

			if category == "Individual" {
				categoryBool = true
			}

			err = database.UpdateRecipient(myUser.RepID, recipientIDInt, name, categoryBool, profile, contactNo)

			if err != nil {
				log.Error(err)
				clientMsg = "Internal server error at database, update unsuccessful"
			} else {
				log.WithFields(logrus.Fields{
					"userName": myUser.UserName,
				}).Infof("[%s] updated recipientID: [%s]", myUser.UserName, recipientIDInt)
				clientMsg = "You have successfully updated " + name + "."
				recipient, err := database.GetRecipient(myUser.RepID, recipientIDInt)
				if err != nil {
					log.Error(err)
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
	}

	recipient, err := database.GetRecipient(myUser.RepID, recipientIDInt)

	if err != nil {
		log.Error(err)
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

// Author: Huang Yanping.
// deleteRecipient is a handler func that delete the recipient
// in the database using the recipientID.
func deleteRecipient(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	clientMsg := ""
	clientMsg2 := ""
	myUser, _ := getUser(w, r)
	authenticate.IsAdmin = database.IsAdmin(myUser.UserName)

	inputs := r.URL.Query()["recipientID"]

	recipientID := inputs[0]

	recipientIDInt, err := strconv.ParseInt(recipientID, 10, 0)

	if err != nil {
		log.Error(err)
		clientMsg2 = "Internal server error, delete unsuccessful."
	}

	err2 := database.DeleteRecipient(myUser.RepID, recipientIDInt)

	if err2 != nil {
		log.Error(err2)
		clientMsg2 = "Internal server error at database, delete unsuccessful."
	} else {
		log.WithFields(logrus.Fields{
			"userName": myUser.UserName,
		}).Infof("[%s] deleted recipientID: [%s]", myUser.UserName, recipientIDInt)
		clientMsg2 = "You have successfully deleted a recipient."
	}
	recipients, err3 := database.GetMyRecipient(myUser.RepID)

	if err3 != nil {
		log.Error(err3)
		clientMsg2 = "Internal server error at database"
	}

	if len(recipients) == 0 {
		clientMsg = "You currently have no recipients"
	} else {
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

// Author: Huang Yanping.
// validateContactNo is a func that accepts a string input
// and validate the string input whether it is a valid
// contact number from Singapore using a external API from numverify.
// return true if it is a valid contact number or
// false if it is invalid and any error using the API.
func validateContactNo(contactNo string) (bool, error) {
	// access key from the API website
	accessKey := "ce2db76b60065070752e50aab06b23a5"

	// url for the API
	url := "http://apilayer.net/api/validate?access_key=" + accessKey + "&number=" + contactNo + "&country_code=SG&format=1"

	// connecting to the API
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		// reading the result
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var result Result
			json.Unmarshal(body, &result)
			return result.Valid, nil
		} else {
			return false, err
		}
	} else {
		return false, err
	}
}
