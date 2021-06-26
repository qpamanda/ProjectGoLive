package server

import (
	"ProjectGoLive/authenticate"
	"ProjectGoLive/database"
	"ProjectGoLive/smtpserver"
	"crypto/rand"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/kennygrant/sanitize"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var nums = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// signup is a handler func to create a new user account.
// Validates user information and creates a new user account.
// Author: Amanda
func signup(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	var (
		clientMsg    string // To display message to the user on the client
		username     string
		password     string
		cmfpassword  string
		firstname    string
		lastname     string
		email        string
		contactno    string
		organisation string
	)

	membertype, _ := database.GetRepMemberType(0)

	// Process form submission
	if req.Method == http.MethodPost {
		// Get form values and sanitize the strings
		username = sanitize.Accents(req.FormValue("username"))
		password = sanitize.Accents(req.FormValue("password"))
		cmfpassword = sanitize.Accents(req.FormValue("cmfpassword"))
		firstname = sanitize.Accents(req.FormValue("firstname"))
		lastname = sanitize.Accents(req.FormValue("lastname"))
		email = sanitize.Accents(req.FormValue("email"))
		contactno = sanitize.Accents(req.FormValue("contactno"))
		organisation = sanitize.Accents(req.FormValue("organisation"))

		for k, v := range membertype {
			membertypeid := req.FormValue("membertype" + strconv.Itoa(k))
			checked := ""
			disabled := ""
			if membertypeid != "" {
				checked = "checked"
			}
			if strings.ToUpper(v.MemberType) == "ADMIN" {
				disabled = "disabled"
			}
			membertype[k] = authenticate.MemberTypeInfo{
				MemberType: v.MemberType,
				Checked:    checked,
				Disabled:   disabled}
		}

		// Validates the input fields from the user
		if err := authenticate.ValidateUserInput(true, username, password, cmfpassword,
			firstname, lastname, email, contactno, organisation); err != nil {

			clientMsg = "ERROR: " + err.Error()
			log.Error(err)
		} else {
			// Check if username exist i.e. exist means it is already taken
			if database.UserNameExist(username) {
				clientMsg = "ERROR: " + "username already taken. Please use another username"
				log.Error("[" + username + "] username already taken")
			} else {
				// Check if email exist i.e. exist means it is already taken
				if database.EmailExist(email, "") {
					clientMsg = "ERROR: " + "email already taken. Please use another email"
					log.Error("[" + username + "] email already taken")
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

					// Store the session token in a map on the server
					authenticate.MapSessions[sessionToken.Value] = username

					// Hashed the password from user input before saving
					bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

					if err != nil {
						clientMsg = "WARNING: " + "internal server error"
						log.Warn("internal server error")
					} else {
						repid, err := GetRepID()

						if err != nil {
							clientMsg = "ERROR: " + "unable to add user"
							log.Error("[" + username + "] unable to add user")
						} else {
							// Insert user into the database
							err = database.AddUser(repid, username, string(bPassword),
								firstname, lastname, email, contactno, organisation)

							for k := range membertype {
								membertypeid := req.FormValue("membertype" + strconv.Itoa(k))

								if membertypeid != "" {
									id, _ := strconv.Atoi(membertypeid)
									// Add member type information by repid into the database
									database.AddRepMemberType(repid, id, username)
								}
							}

							if err != nil {
								log.WithFields(logrus.Fields{
									"repId":    repid,
									"userName": username,
								}).Errorf("[%s] error creating user account ", username)
							} else {
								log.WithFields(logrus.Fields{
									"repId":    repid,
									"userName": username,
								}).Infof("[%s] user account created successfully", username)
							}
							// Redirect to the main index page
							http.Redirect(res, req, "/", http.StatusSeeOther)
							return
						}
					}
				}
			}
		}
	}

	data := struct {
		UserName     string
		Password     string
		CmfPassword  string
		FirstName    string
		LastName     string
		Email        string
		ContactNo    string
		Organisation string
		MemberType   map[int]authenticate.MemberTypeInfo
		ClientMsg    string
	}{
		username,
		password,
		cmfpassword,
		firstname,
		lastname,
		email,
		contactno,
		organisation,
		membertype,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "signup.gohtml", data)
}

// edituser is a handler func to edit user account information.
// Redirects to index page if user has not login.
// Validates user input and updates the information.
// Author: Amanda
func edituser(res http.ResponseWriter, req *http.Request) {
	myUser, validSession := getUser(res, req)

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := "" // To display message to the user on the client

	// Set current user information to be display on the form upon first form load
	repid := myUser.RepID
	username := myUser.UserName
	password := myUser.Password
	cmfpassword := myUser.Password
	firstname := myUser.FirstName
	lastname := myUser.LastName
	email := myUser.Email
	contactno := myUser.ContactNo
	organisation := myUser.Organisation

	membertype, _ := database.GetRepMemberType(repid)

	// Process the form submission
	if req.Method == http.MethodPost {
		// Get form values and sanitize the strings
		username = sanitize.Accents(req.FormValue("username"))
		firstname = sanitize.Accents(req.FormValue("firstname"))
		lastname = sanitize.Accents(req.FormValue("lastname"))
		email = sanitize.Accents(req.FormValue("email"))
		contactno = sanitize.Accents(req.FormValue("contactno"))
		organisation = sanitize.Accents(req.FormValue("organisation"))

		for k, v := range membertype {
			membertypeid := req.FormValue("membertype" + strconv.Itoa(k))

			checked := ""
			disabled := ""
			if membertypeid != "" {
				checked = "checked"
			}

			if strings.ToUpper(v.MemberType) == "ADMIN" {
				checked = "checked"
				disabled = "disabled"
			}

			membertype[k] = authenticate.MemberTypeInfo{
				MemberType: v.MemberType,
				Checked:    checked,
				Disabled:   disabled}
		}

		// Validates the input fields from the user
		if err := authenticate.ValidateUserInput(false, username, password, cmfpassword, firstname, lastname,
			email, contactno, organisation); err != nil {
			clientMsg = "ERROR: " + err.Error()
			log.Error(err)
		} else {
			// Check if email exist i.e. exist means it is already taken
			if database.EmailExist(email, username) {
				clientMsg = "ERROR: " + "email already taken. Please use another email"
				log.Error("[" + username + "] email already taken")
			} else {

				// Update user information into the database
				err = database.UpdateUser(repid, username, firstname, lastname, email, contactno, organisation)

				// Delete member type information from the database
				database.DeleteRepMemberType(repid)

				for k := range membertype {
					membertypeid := req.FormValue("membertype" + strconv.Itoa(k))

					if membertypeid != "" {
						id, _ := strconv.Atoi(membertypeid)
						// Add member type information by repid into the database
						database.AddRepMemberType(repid, id, username)
					}
				}

				if err != nil {
					log.WithFields(logrus.Fields{
						"userName": username,
					}).Errorf("[%s] error updating user account", username)
				} else {
					log.WithFields(logrus.Fields{
						"userName": username,
					}).Infof("[%s] user account updated successfully", username)

					// Redirect to the main index page
					http.Redirect(res, req, "/", http.StatusSeeOther)
					return
				}
			}
		}
	}

	data := struct {
		User         authenticate.User
		UserName     string
		Password     string
		CmfPassword  string
		FirstName    string
		LastName     string
		Email        string
		ContactNo    string
		Organisation string
		MemberType   map[int]authenticate.MemberTypeInfo
		ClientMsg    string
		ValidSession bool
	}{
		myUser,
		username,
		password,
		cmfpassword,
		firstname,
		lastname,
		email,
		contactno,
		organisation,
		membertype,
		clientMsg,
		validSession,
	}
	tpl.ExecuteTemplate(res, "edituser.gohtml", data)
}

// changepwd is a handler func to change user password
// Redirects to index page if user has not login.
// Validates user input and updates the information.
// Author: Amanda
func changepwd(res http.ResponseWriter, req *http.Request) {
	myUser, validSession := getUser(res, req)

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := "" // To display message to the user on the client

	// Set current user information to be display on the form upon first form load
	repid := myUser.RepID
	username := myUser.UserName
	hashpassword := myUser.Password

	password := ""
	newpassword := ""
	cmfpassword := ""

	// Process the form submission
	if req.Method == http.MethodPost {
		// Get form values and sanitize the strings
		password = sanitize.Accents(req.FormValue("password"))
		newpassword = sanitize.Accents(req.FormValue("newpassword"))
		cmfpassword = sanitize.Accents(req.FormValue("cmfpassword"))

		if password == "" {
			err := errors.New("current password cannot be blank")
			clientMsg = "ERROR: " + err.Error()
			log.Error(err)
		} else {
			// Validates the input fields from the user
			if err := authenticate.ValidatePassword(newpassword, cmfpassword); err != nil {
				clientMsg = "ERROR: " + err.Error()
				log.Error(err)
			} else {
				// Matching of password entered
				err := bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(password))

				if err != nil {
					clientMsg = "ERROR: " + "password is incorrect"
					log.Error("password is incorrect")
				} else {

					// Hashed the password from user input before saving
					bPassword, err := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.MinCost)

					if err != nil {
						clientMsg = "WARNING: " + "internal server error"
						log.Warn("internal server error")
					} else {

						// Update password information into the database
						err = database.UpdatePassword(repid, username, string(bPassword))

						if err != nil {
							clientMsg = "ERROR: " + "error updating password"
							log.WithFields(logrus.Fields{
								"userName": username,
							}).Errorf("[%s] error updating password", username)
						} else {
							clientMsg = "Password updated successfully"
							log.WithFields(logrus.Fields{
								"userName": username,
							}).Infof("[%s] password updated successfully", username)

							password = ""
							newpassword = ""
							cmfpassword = ""
						}
					}
				}
			}
		}
	}

	data := struct {
		User         authenticate.User
		Password     string
		NewPassword  string
		CmfPassword  string
		ClientMsg    string
		ValidSession bool
	}{
		myUser,
		password,
		newpassword,
		cmfpassword,
		clientMsg,
		validSession,
	}
	tpl.ExecuteTemplate(res, "changepwd.gohtml", data)
}

// resetpwd is a handler func to reset user password without login
// Validates user input and updates the information.
// Author: Amanda
func resetpwd(res http.ResponseWriter, req *http.Request) {
	clientMsg := "" // To display message to the user on the client
	repid := 0
	username := ""

	v := req.URL.Query()
	if key, ok := v["key"]; ok {
		hashemail := sanitize.Accents(key[0]) // Sanitise the url param string

		users, err := database.GetAllUsers()

		if err != nil {
			clientMsg = "ERROR: " + "unable to reset password"
			log.WithFields(logrus.Fields{
				"hashemail": hashemail,
			}).Errorf("unable to reset password")
		} else {
			for _, v := range users {
				// Matching of password entered
				err := bcrypt.CompareHashAndPassword([]byte(hashemail), []byte(v.Email))

				if err == nil {
					repid = v.RepID
					username = v.UserName
					break
				}
			}
		}

		if username == "" {
			clientMsg = "ERROR: " + "user account not found to reset password"
			log.WithFields(logrus.Fields{
				"hashemail": hashemail,
			}).Errorf("user account not found to reset password")
		}
	}

	newpassword := ""
	cmfpassword := ""

	// Process the form submission
	if req.Method == http.MethodPost {
		// Get form values and sanitize the strings
		newpassword = sanitize.Accents(req.FormValue("newpassword"))
		cmfpassword = sanitize.Accents(req.FormValue("cmfpassword"))

		// Validates the input fields from the user
		if err := authenticate.ValidatePassword(newpassword, cmfpassword); err != nil {
			clientMsg = "ERROR: " + err.Error()
			log.Error(err)
		} else {

			// Hashed the password from user input before saving
			bPassword, err := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.MinCost)

			if err != nil {
				clientMsg = "WARNING: " + "internal server error"
				log.Warn("internal server error")
			} else {
				// Update password information into the database
				err = database.UpdatePassword(repid, username, string(bPassword))

				if err != nil {
					clientMsg = "ERROR: " + "error resetting password"
					log.WithFields(logrus.Fields{
						"username": username,
					}).Errorf("[%s] error resetting password", username)
				} else {
					clientMsg = "Password reset successfully. Please login with your new password."
					log.WithFields(logrus.Fields{
						"userName": username,
					}).Infof("[%s] password reset successfully", username)

					newpassword = ""
					cmfpassword = ""
				}
			}
		}
	}

	// This is an intended empty struct to get pass header.gohtml checks
	myUser := struct {
		UserName    string
		LastLoginDT time.Time
	}{
		"",
		time.Now(),
	}

	data := struct {
		User struct {
			UserName    string
			LastLoginDT time.Time
		}
		UserName    string
		NewPassword string
		CmfPassword string
		ClientMsg   string
	}{
		myUser,
		username,
		newpassword,
		cmfpassword,
		clientMsg,
	}
	tpl.ExecuteTemplate(res, "resetpwd.gohtml", data)
}

// resetpwd is a handler func to reset user password
// It triggers an email to user registered email address to activate password reset
// Author: Amanda
func resetpwdreq(res http.ResponseWriter, req *http.Request) {
	clientMsg := "" // To display message to the user on the client

	email := ""

	// Process the form submission
	if req.Method == http.MethodPost {
		// Get form values and sanitize the strings
		email = sanitize.Accents(req.FormValue("email"))

		// Validates the input fields from the user
		if err := authenticate.ValidateEmail(email); err != nil {
			clientMsg = "ERROR: " + err.Error()
			log.Error(err)
		} else {

			err = smtpserver.EmailResetPassword(email)

			if err != nil {
				clientMsg = "ERROR: " + err.Error()
				log.WithFields(logrus.Fields{
					"email": email,
				}).Errorf(err.Error())
			} else {
				clientMsg = "Reset password email sent"
				log.WithFields(logrus.Fields{
					"email": email,
				}).Infof("Reset password email sent")

				email = ""
			}
		}
	}

	// This is an intended empty struct to get pass header.gohtml checks
	myUser := struct {
		UserName    string
		LastLoginDT time.Time
	}{
		"",
		time.Now(),
	}

	data := struct {
		User struct {
			UserName    string
			LastLoginDT time.Time
		}
		Email     string
		ClientMsg string
	}{
		myUser,
		email,
		clientMsg,
	}
	tpl.ExecuteTemplate(res, "resetpwdreq.gohtml", data)
}

// logout func is a handler to logout the current user. Redirects to index page if user has not login.
// Otherwise, delete session token from server and client, then redirects to index page.
func logout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	sessionToken, _ := req.Cookie("sessionToken")

	// Get username before session is deleted
	username := authenticate.MapSessions[sessionToken.Value]

	// Delete the session token from the server
	delete(authenticate.MapSessions, sessionToken.Value)
	// Remove the cookie from the client
	sessionToken = &http.Cookie{
		Name:   "sessionToken",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, sessionToken)

	log.WithFields(logrus.Fields{
		"userName": username,
	}).Infof("[%s] user logout successfully", username)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// getUser func checks for valid session token.
// Add a new session token cookie to the client if one is not found.
// Return user and true if session is set
// Author: Amanda
func getUser(res http.ResponseWriter, req *http.Request) (authenticate.User, bool) {
	// Get current session cookie
	sessionToken, err := req.Cookie("sessionToken")
	// No session token found
	if err != nil {
		// Add new session token cookie
		id, _ := uuid.NewV4()
		// Set an expiry time of 120 seconds for the cookie, the same as the cache
		sessionToken = &http.Cookie{
			Name:     "sessionToken",
			Value:    id.String(),
			Expires:  time.Now().Add(120 * time.Second),
			HttpOnly: true,
			Path:     "/",
			Domain:   "localhost",
			Secure:   true,
		}
	}
	http.SetCookie(res, sessionToken)

	// If the user exists already, get user
	var myUser authenticate.User

	if _, ok := authenticate.MapSessions[sessionToken.Value]; ok {
		username := authenticate.MapSessions[sessionToken.Value]
		myUser, err = database.GetUser(username) // Get user from database
		if err != nil {
			return myUser, false
		}
		return myUser, ok
	}
	return myUser, false
}

// alreadyLoggedIn func checks if a user has already logged in. Checks for valid session token.
// Returns true if already logged in, false otherwise.
// Author: Amanda
func alreadyLoggedIn(req *http.Request) bool {
	sessionToken, err := req.Cookie("sessionToken")
	if err != nil {
		return false
	}

	// Get username from session map
	username := authenticate.MapSessions[sessionToken.Value]
	return database.UserNameExist(username)
}

func GetRepID() (int, error) {
	x := 4 // change this value for the number of digits to generate for the random no
	randomNo := GenerateRandomNo(x)

	return strconv.Atoi(randomNo)
}

func GenerateRandomNo(x int) string {
	buffer := make([]byte, x)
	n, err := io.ReadAtLeast(rand.Reader, buffer, x)

	if err != nil || n != x {
		return GenerateRandomNo(x)
	}

	for i := 0; i < len(buffer); i++ {
		buffer[i] = nums[int(buffer[i])%len(nums)]
	}

	if database.RepIDExist(string(buffer)) {
		return GenerateRandomNo(x)
	}

	return string(buffer)
}
