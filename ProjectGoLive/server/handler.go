package server

import (
	"ProjectGoLive/authenticate"
	"ProjectGoLive/database"
	"net"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// index is a handler func that display the home page of the application.
// On start, it will default as the login page first. Once user login,
// the page will change to show the main menu for the users.
// If user is an admin, it will display the admin menu as well.
// Author: Amanda
func index(res http.ResponseWriter, req *http.Request) {

	clientMsg := "" // To display client-side message to user

	// Process the form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")

		if username == "" || password == "" {
			clientMsg = "ERROR: username and/or password cannot be blank"
			log.Error("username and/or password cannot be blank")
			return
		} else {
			if !database.UserNameExist(username) {
				clientMsg = "ERROR: username and/or password do not match"
				log.Error("username and/or password do not match")
				return
			} else {
				hashpassword, err := database.GetHashPassword(username)
				if err != nil {
					clientMsg = "ERROR: username and/or password do not match"
					log.Error("username and/or password do not match")
					return
				} else {
					// Matching of password entered
					err := bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(password))
					if err != nil {
						clientMsg = "ERROR: username and/or password do not match"
						log.Error("username and/or password do not match")
						return
					} else {
						// Create a new session token
						id, _ := uuid.NewV4()

						// Set an expiry time of 120 seconds for the cookie, the same as the cache
						sessionToken := &http.Cookie{
							Name:     "sessionToken",
							Value:    id.String(),
							Expires:  time.Now().Add(120 * time.Second),
							HttpOnly: true,
							Path:     "/",
							Domain:   "localhost",
							Secure:   true,
						}
						// Set the session token as a cookie on the client
						http.SetCookie(res, sessionToken)
						// Set user to session token cookie
						authenticate.MapSessions[sessionToken.Value] = username

						// Checks if user is an admin
						authenticate.IsAdmin = database.IsAdmin(username)

						// Update LastLogin datetime in database
						database.UpdateLoginDate(username)

						log.WithFields(logrus.Fields{
							"userName": username,
						}).Infof("[%s] user login successfully", username)
					}
				}
			}
		}
	}

	myUser, validSession := getUser(res, req)

	data := struct {
		User         authenticate.User
		ValidSession bool
		IsAdmin      bool
		ClientMsg    string
	}{
		myUser,
		validSession,
		authenticate.IsAdmin,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "index.gohtml", data)
}

// createCookie func creates sets the struct for a cookie
// Author: Amanda
func createCookie(res http.ResponseWriter, req *http.Request) *http.Cookie {
	domain := req.Host                      // the domain can be localhost:5221 or //127.0.0.1:5221
	host, _, _ := net.SplitHostPort(domain) // get the host either localhost or 127.0.0.1

	// Add new session token cookie
	id, _ := uuid.NewV4()
	// Set an expiry time of 120 seconds for the cookie
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    id.String(),
		Expires:  time.Now().Add(120 * time.Second),
		HttpOnly: true,
		Path:     "/",
		Domain:   host, // set cookie with the host
		Secure:   true,
	}
	// Set the session token as a cookie on the client
	http.SetCookie(res, cookie)
	return cookie
}
