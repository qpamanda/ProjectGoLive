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

// Author: Amanda Soh.
// index is a handler func that display the home page of the application.
// On start, it will default as the login page first. Once user login,
// the page will change to show the main menu for the users.
// If user is an admin, it will display the admin menu as well.
func index(res http.ResponseWriter, req *http.Request) {

	clientMsg := "" // To display client-side message to user

	// Process the form submission
	if req.Method == http.MethodPost {
		// Reset MapUsers and MapSessions for new login
		resetSession()

		username := req.FormValue("username")
		password := req.FormValue("password")

		if username == "" || password == "" {
			clientMsg = "ERROR: username and/or password cannot be blank"
			log.Error("username and/or password cannot be blank")
		} else {
			if !database.UserNameExist(username) {
				clientMsg = "ERROR: username and/or password do not match"
				log.Error("username and/or password do not match")
			} else {
				hashpassword, err := database.GetHashPassword(username)
				if err != nil {
					clientMsg = "ERROR: username and/or password do not match"
					log.Error("username and/or password do not match")
				} else {
					// Matching of password entered
					err := bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(password))
					if err != nil {
						clientMsg = "ERROR: " + "username and/or password do not match"
						log.Error("username and/or password do not match")
					} else {
						// Call createCookie func to set the cookie
						sessionToken := createCookie(res, req)

						authenticate.MapSessions[sessionToken.Value] = username

						myUser, err := database.GetUser(username) // Get user from database
						if err != nil {
							clientMsg = "ERROR: " + "retrieving user information"
							log.Error("error retrieving user information")
						}

						// Set user to map users
						authenticate.MapUsers[username] = myUser

						// Update LastLogin datetime in database
						database.UpdateLoginDate(username)

						log.WithFields(logrus.Fields{
							"userName": username,
						}).Infof("[%s] user login successfully", username)

						// Redirect to the main index page
						http.Redirect(res, req, "/", http.StatusSeeOther)
						return
						//}
					}
				}
			}
		}
	}

	myUser, validSession := getUser(res, req)

	data := struct {
		User         authenticate.User
		ValidSession bool
		ClientMsg    string
	}{
		myUser,
		validSession,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "index.gohtml", data)
}

// Author: Amanda Soh.
// createCookie func creates sets the struct for a cookie
func createCookie(res http.ResponseWriter, req *http.Request) *http.Cookie {
	domain := req.Host                      // the domain can be localhost:5221 or //127.0.0.1:5221
	host, _, _ := net.SplitHostPort(domain) // get the host either localhost or 127.0.0.1

	// Add new session token cookie
	id, _ := uuid.NewV4()
	// Set an expiry time of 1200 seconds (20mins) for the cookie
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    id.String(),
		Expires:  time.Now().Add(1200 * time.Second),
		HttpOnly: true,
		Path:     "/",
		Domain:   host, // set cookie with the host
		Secure:   true,
	}
	// Set the session token as a cookie on the client
	http.SetCookie(res, cookie)
	return cookie
}

// Author: Amanda Soh.
// resetSession delete map data prior to create for new session
func resetSession() {
	for k1 := range authenticate.MapUsers {
		// Delete the map user from the server
		delete(authenticate.MapUsers, k1)
	}

	for k2 := range authenticate.MapSessions {
		// Delete the session token from the server
		delete(authenticate.MapSessions, k2)
	}
}
