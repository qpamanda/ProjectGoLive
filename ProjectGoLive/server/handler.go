package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

/*
// index is the handler function to display the home page of the server.
func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", "")
}
*/

// index is a handler func that display the home page of the application.
// On start, it will default as the login page first. Once user login,
// the page will change to show the main menu for the users.
// If user is an admin, it will display the admin menu as well.
func index(res http.ResponseWriter, req *http.Request) {
	bFirst = true
	clientMsg := "" // To display client-side message to user

	// Process the form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")

		if username == "" || password == "" {
			clientMsg = "ERROR: username and/or password cannot be blank"
			log.Error("username and/or password cannot be blank")
		} else {
			// Check if user exist with username
			myUser, ok := mapUsers[username]

			if !ok {
				clientMsg = "ERROR: username and/or password do not match"
				log.Error("username and/or password do not match")
			} else {
				// Matching of password entered
				err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
				if err != nil {
					clientMsg = "ERROR: " + "username and/or password do not match"
					log.Error("username and/or password do not match")
				} else {
					sessionToken, err := req.Cookie("sessionToken")
					if err != nil {
						clientMsg = "ERROR: " + "session cookie not found"
						log.Error("session cookie not found")
					} else {
						http.SetCookie(res, sessionToken)
						// Set user to session token cookie
						mapSessions[sessionToken.Value] = username

						updateLoginDate(myUser)

						log.WithFields(logrus.Fields{
							"userName": username,
						}).Infof("[%s] user login successfully", username)
					}
				}
			}
		}
	}

	myUser := getUser(res, req)

	data := struct {
		User      user
		ClientMsg string
	}{
		myUser,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "index.gohtml", data)
}
