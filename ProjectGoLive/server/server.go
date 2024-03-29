/*
Author: Ahmad Bahrudin, Amanda Soh, Huang Yanping, Tan Jun Jie.

Package server initialises the handler functions for the server web pages
and implements the functionalities each handler.

It is separated into serveral .go files to segregate the functionalities of the application done
by each team members.

	server.go: Initialises the templates and handler functions, then starts the server to run
	on the designated port.

	handler.go: Implements the handler functions for displaying the web pages of the server
*/
package server

import (
	"ProjectGoLive/authenticate"
	"ProjectGoLive/database"
	"ProjectGoLive/smtpserver"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	filename "github.com/keepeye/logrus-filename"
	"github.com/sirupsen/logrus"
)

var (
	tpl  *template.Template
	log  = logrus.New()
	file *os.File
)

const cookieName = "sessionToken"

// req struct for storing request information
type newRequest struct {
	RepresentativeId int // id of the coordinator/representative
	/*
		RequestCategoryId
		1 (item donation)
		2 (errands)
	*/
	RequestCategoryId int
	RecipientId       int // id of recipient who receives the aid
	/*
		RequestStatus
		0 (pending/waiting to be matched to a helper)
		1 (being handled)
		2 (completed)
	*/
	RequestStatus  int
	RequestDetails requestDetails
	CreatedBy      string
	CreatedDT      time.Time
	LastModifiedBy string
	LastModifiedDT time.Time
}

//requestDetails struct for storing request detail information
type requestDetails struct {
	RequestDescription string
	ToCompleteBy       time.Time
	FulfilAt           string
}

// viewRecipient struct for storing a view on Recipient details
type viewRecipient struct {
	RecipientID int
	Name        string
}

// viewRequest struct for storing a view on Request details
type viewRequest struct {
	RequestID     int
	Category      string
	RecipientName string
	Description   string
	ToCompleteBy  string
	FulfillAt     string
	Status        string
}

// Author: Amanda Soh.
// InitServer initialises the templates for displaying the web pages at the server.
// It also creates and opens the log file for events logging.
func InitServer() {

	// Log file name is based on current day. Thus a new file is created for each day.
	date := time.Now().Format("20060102")
	logFileName := "log/" + date + "_events.log"

	// Create a new log file for append
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("FATAL: OpenFile - ", err)
	}

	// Log to events.log file
	log.SetOutput(file)
	// Set log formatter
	log.SetFormatter(&logrus.JSONFormatter{})
	// Set log level from Info level onwards
	log.SetLevel(logrus.InfoLevel)

	// Use 3rd party package filename to display filename and line no during logging
	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	log.AddHook(filenameHook)
}

// Author: Amanda Soh.
// StartServer initialises the database and handler functions then
// listens on the designated port to start the server running.
func StartServer() {
	// Initialise the database
	initDB()

	initFieldsLen()

	router := mux.NewRouter()

	// Initialise the handlers
	initaliseHandlers(router)

	// Testing functions created by Ahmad
	//test()
	//testDelete()

	// Set the listen port
	fmt.Println("Listening at port 5221")
	err := http.ListenAndServeTLS(":5221", "certs//cert.pem", "certs//key.pem", router)
	//err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("FATAL: ListenAndServeTLS - ", err)
	}

	// Defer file closure to the end
	defer file.Close()
	defer database.DB.Close()
}

// Author: Amanda Soh, Huang Yanping, Tan Jun Jie, Ahmad Bahrudin.
// initaliseHandlers initialises the handlers for the server.
func initaliseHandlers(router *mux.Router) {

	router.HandleFunc("/", index)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/signup", signup)
	router.HandleFunc("/edituser", edituser)
	router.HandleFunc("/changepwd", changepwd)
	router.HandleFunc("/managerecipient", manageRecipient)
	router.HandleFunc("/addrecipient", addRecipient)
	router.HandleFunc("/getrecipient", getRecipient)
	router.HandleFunc("/updaterecipient", updateRecipient)
	router.HandleFunc("/deleterecipient", deleteRecipient)
	router.HandleFunc("/resetpwd", resetpwd)
	router.HandleFunc("/resetpwdreq", resetpwdreq)
	router.HandleFunc("/addrequest", addrequest)
	router.HandleFunc("/deleterequest", deleterequest)
	router.HandleFunc("/selectrequest", selectrequest)
	router.HandleFunc("/fulfilrequest", fulfilrequest)
	router.HandleFunc("/requestcompleted", requestcompleted)
	router.HandleFunc("/selecteditrequest", selecteditrequest)
	router.HandleFunc("/editrequest", editrequest)
	router.HandleFunc("/viewrequest", viewrequest)
	router.HandleFunc("/managerequest", managerequest)
	router.HandleFunc("/aaCatAdd", aaCatAdd)
	router.HandleFunc("/aaCatUpdate", aaCatUpdate)
	router.HandleFunc("/aaCatUpdate2", aaCatUpdate2)
	router.HandleFunc("/aaCatView", aaCatView)
	router.HandleFunc("/aaMemTypeAdd", aaMemTypeAdd)
	router.HandleFunc("/aaMemTypeUpdate", aaMemTypeUpdate)
	router.HandleFunc("/aaMemTypeUpdate2", aaMemTypeUpdate2)
	router.HandleFunc("/aaMemTypeView", aaMemTypeView)
	router.HandleFunc("/aaReqSAdd", aaReqSAdd)
	router.HandleFunc("/aaReqSUpdate", aaReqSUpdate)
	router.HandleFunc("/aaReqSUpdate2", aaReqSUpdate2)
	router.HandleFunc("/aaReqSView", aaReqSView)

	router.Handle("/favicon.ico", http.NotFoundHandler())

	router.Handle("/static/main.css", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	router.Handle("/static/img/logo.jpg", http.StripPrefix("/static/img", http.FileServer(http.Dir("./static/img"))))
	router.Handle("/static/img/user.jpg", http.StripPrefix("/static/img", http.FileServer(http.Dir("./static/img"))))

	// Parse templates
	tpl = template.Must(template.ParseGlob("templates/*"))
}

// Author: Amanda Soh.
// initDB initialises the database
func initDB() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	config := getDBConfig()

	// Get database connection string
	connectionString := database.GetConnectionString(config)

	// Connect to database
	err := database.Connect(connectionString)
	if err != nil {
		panic("error connecting to database")
	}

	// Test connection to database
	err = database.DB.Ping()
	if err != nil {
		panic("error pinging to database")
	} else {
		fmt.Println("Ping to database success")
	}

	database.DB.SetMaxOpenConns(0)
}

// Author: Amanda Soh.
// getDBConfig retrieves the database configurations and returns a struct.
func getDBConfig() database.Config {
	// Load setup.env file from same directory
	err := godotenv.Load("setup.env")
	if err != nil {
		log.Fatal("FATAL: Error loading .env file")
	}

	// Get env variables for database configuration
	serverName := os.Getenv("SERVER_NAME")
	dbName := os.Getenv("DB_NAME")
	dbUserName := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	config :=
		database.Config{
			ServerName: serverName,
			User:       dbUserName,
			Password:   dbPassword,
			DB:         dbName,
		}

	return config
}

// Author: Amanda Soh.
// Initiate fields from .env file
func initFieldsLen() {
	// Load setup.env file from same directory
	err := godotenv.Load("setup.env")
	if err != nil {
		log.Fatal("FATAL: Error loading .env file")
	}

	// Set the min characters for username
	authenticate.MinUserName, _ = strconv.Atoi(os.Getenv("MIN_USERNAME"))

	// Set the max characters for username
	authenticate.MaxUserName, _ = strconv.Atoi(os.Getenv("MAX_USERNAME"))

	// Set the min characters for password
	authenticate.MinPassword, _ = strconv.Atoi(os.Getenv("MIN_PASSWORD"))

	// Set the max characters for password
	authenticate.MaxPassword, _ = strconv.Atoi(os.Getenv("MAX_PASSWORD"))

	// Setup fields for email sending feature
	smtpserver.HostPath = os.Getenv("HOST_PATH")
	smtpserver.SMTPHost = os.Getenv("SMTP_HOST")
	smtpserver.SMTPPort = os.Getenv("SMTP_PORT")
	smtpserver.EmailPassword = os.Getenv("EMAIL_PASSWORD")
	smtpserver.FromEmail = os.Getenv("FROM_EMAIL")
}
