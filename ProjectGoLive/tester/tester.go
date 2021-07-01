package main

// Author: Ahmad Bahrudin
import (
	"ProjectGoLive/database"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func db() {
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

	// Get database connection string
	connectionString := database.GetConnectionString(config)

	// Connect to database
	err = database.Connect(connectionString)
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

func TestCatInsert() {
	if database.Cat.Insert(90001, "Category TESTING", "admin") != nil &&
		database.Cat.Insert(90002, "Category TESTING", "admin") != nil &&
		database.Cat.Insert(90003, "Category TESTING", "admin") != nil {
		fmt.Println("Category Insert Function: Test Failed!")
	} else {
		fmt.Println("Category Insert Function: Test Success!")
	}
}

func TestCatUpdate() {
	if database.Cat.Update(90001, "Category TESTING 1", "admin") != nil &&
		database.Cat.Update(90002, "Category TESTING 1", "admin") != nil &&
		database.Cat.Update(90003, "Category TESTING 1", "admin") != nil {
		fmt.Println("Category Update Function: Test Failed!")
	} else {
		fmt.Println("Category Update Function: Test Success!")
	}
}

func TestCatGetAll() {
	if _, err := database.Cat.GetAll(); err != nil {
		fmt.Println("Category GetAll Function: Test Failed!")
	} else {
		fmt.Println("Category GetAll Function: Test Success!")
	}
}

func TestCatDelete() {
	if database.Cat.Delete(90001) != nil &&
		database.Cat.Delete(90002) != nil &&
		database.Cat.Delete(90003) != nil {
		fmt.Println("Category Delete Function: Test Failed!")
	} else {
		fmt.Println("Category Delete Function: Test Success!")
	}
}

func TestMemTInsert() {
	if database.MemT.Insert(90001, "Member Type TESTING", "admin") != nil &&
		database.MemT.Insert(90002, "Member Type TESTING", "admin") != nil &&
		database.MemT.Insert(90003, "Member Type TESTING", "admin") != nil {
		fmt.Println("Member Type Insert Function: Test Failed!")
	} else {
		fmt.Println("Member Type Insert Function: Test Success!")
	}
}

func TestMemTUpdate() {
	if database.MemT.Update(90001, "Member Type TESTING 1", "admin") != nil &&
		database.MemT.Update(90002, "Member Type TESTING 1", "admin") != nil &&
		database.MemT.Update(90003, "Member Type TESTING 1", "admin") != nil {
		fmt.Println("Member Type Update Function: Test Failed!")
	} else {
		fmt.Println("Member Type Update Function: Test Success!")
	}
}

func TestMemTGetAll() {
	if _, err := database.MemT.GetAll(); err != nil {
		fmt.Println("Member Type GetAll Function: Test Failed!")
	} else {
		fmt.Println("Member Type GetAll Function: Test Success!")
	}
}

func TestMemTDelete() {
	if database.MemT.Delete(90001) != nil &&
		database.MemT.Delete(90002) != nil &&
		database.MemT.Delete(90003) != nil {
		fmt.Println("Member Type Delete Function: Test Failed!")
	} else {
		fmt.Println("Member Type Delete Function: Test Success!")
	}
}

func TestReqSInsert() {
	if database.ReqS.Insert(91, "Request Status TESTING", "admin") != nil &&
		database.ReqS.Insert(92, "Request Status TESTING", "admin") != nil &&
		database.ReqS.Insert(93, "Request Status TESTING", "admin") != nil {
		fmt.Println("Request Status Insert Function: Test Failed!")
	} else {
		fmt.Println("Request Status Insert Function: Test Success!")
	}
}

func TestReqSUpdate() {
	if database.ReqS.Update(91, "Request Status TESTING 1", "admin") != nil &&
		database.ReqS.Update(92, "Request Status TESTING 1", "admin") != nil &&
		database.ReqS.Update(93, "Request Status TESTING 1", "admin") != nil {
		fmt.Println("Request Status Update Function: Test Failed!")
	} else {
		fmt.Println("Request Status Update Function: Test Success!")
	}
}

func TestReqSGetAll() {
	if _, err := database.ReqS.GetAll(); err != nil {
		fmt.Println("Request Status GetAll Function: Test Failed!")
	} else {
		fmt.Println("Request Status GetAll Function: Test Success!")
	}
}

func TestReqSDelete() {
	if database.ReqS.Delete(91) != nil &&
		database.ReqS.Delete(92) != nil &&
		database.ReqS.Delete(93) != nil {
		fmt.Println("Member Type Delete Function: Test Failed!")
	} else {
		fmt.Println("Member Type Delete Function: Test Success!")
	}
}

func main() {
	fmt.Println("Begin Test...")
	db()
	defer database.DB.Close()

	TestCatInsert()
	TestCatUpdate()
	TestCatGetAll()
	TestCatDelete()

	TestMemTInsert()
	TestMemTUpdate()
	TestMemTGetAll()
	TestMemTDelete()

	TestReqSInsert()
	TestReqSUpdate()
	TestReqSGetAll()
	TestReqSDelete()

	fmt.Println("Test Ended")
}
