/*
Package server initialises the handler functions for the server web pages
and implements ......
It is separated into x .go files to segregate the functionalities of the application.

	server.go: Initialises the templates and handler functions, then starts the server to run
	on the designated port.

	handler.go: Implements the handler functions for displaying the web pages of the server
*/
package server

import (
	"ProjectGoLive/authenticate"
	"ProjectGoLive/database"
	"ProjectGoLive/testing"
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

	//bFirst = true
)

// user struct for storing user account information
type user struct {
	UserName       string
	Password       []byte
	FirstName      string
	LastName       string
	Email          string
	IsAdmin        bool
	CreatedDT      time.Time
	LastModifiedDT time.Time
	CurrentLoginDT time.Time
	LastLoginDT    time.Time
}

// req struct for storing request information
type newRequest struct {
	RepresentativeId int // id of the coordinator/representative
	/*
		RequestCategoryId
		1 (monetary donation)
		2 (item donation)
		3 (errands)
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

type viewRecipient struct {
	RecipientID int
	Name        string
}

type viewRequest struct {
	RequestID     int
	Category      string
	RecipientName string
	Description   string
	ToCompleteBy  string
}

// InitServer initialises the templates for displaying the web pages at the server.
// It also creates and opens the log file for events logging.
func InitServer() {
	// Parse templates
	tpl = template.Must(template.ParseGlob("templates/*"))

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
}

// initaliseHandlers initialises the handlers for the server.
func initaliseHandlers(router *mux.Router) {

	router.HandleFunc("/", index)

	// ADD HANDLERFUNC BELOW
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/signup", signup)
	router.HandleFunc("/edituser", edituser)
	router.HandleFunc("/changepwd", changepwd)
	router.HandleFunc("/categorytable", categorytable)
	router.HandleFunc("/membertypetable", membertypetable)
	router.HandleFunc("/requeststatustable", requeststatustable)
	//router.HandleFunc("/delcourse", delcourse)
	//router.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./img"))))
	router.Handle("/favicon.ico", http.NotFoundHandler())
}

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
}

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
}

// Author: Ahmad Bahrudin
func test() {
	testing.TestCatInsert()
	testing.TestCatUpdate()
	testing.TestCatGet()
	testing.TestCatGetAll()

	testing.TestMemTInsert()
	testing.TestMemTUpdate()
	testing.TestMemTGet()
	testing.TestMemTGetAll()

	testing.TestReqSInsert()
	testing.TestReqSUpdate()
	testing.TestReqSGet()
	testing.TestReqSGetAll()
}

// Author: Ahmad Bahrudin
func testDelete() {
	testing.TestCatDelete()
	testing.TestMemTDelete()
	testing.TestReqSDelete()
}
