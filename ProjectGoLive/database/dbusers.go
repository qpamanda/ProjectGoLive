package database

import (
	"ProjectGoLive/authenticate"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// GetUser implements the sql operations to retrieve a user.
// Author: Amanda
func GetUser(username string) (authenticate.User, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	// Instantiate user
	var (
		user         authenticate.User
		repid        int
		password     string
		firstname    string
		lastname     string
		email        string
		contactno    string
		organisation string
		lastloginDT  time.Time
	)
	query := "SELECT RepID, Password, " +
		"FirstName, LastName, Email, " +
		"ContactNo, Organisation, LastLogin_dt " +
		"FROM Representatives WHERE UserName=?"

	results, err := DB.Query(query, username)
	if err != nil {
		panic("error executing sql select")
	} else {
		if results.Next() {
			err := results.Scan(&repid, &password, &firstname, &lastname,
				&email, &contactno, &organisation, &lastloginDT)

			if err != nil {
				panic("error getting results from sql select")
			}
		} else {
			return user, errors.New("user not found")
		}

		user = authenticate.User{
			RepID:        repid,
			UserName:     username,
			Password:     password,
			FirstName:    firstname,
			LastName:     lastname,
			Email:        email,
			ContactNo:    contactno,
			Organisation: organisation,
			LastLoginDT:  lastloginDT,
		}
		return user, nil
	}
}

// GetAllUsers implements the sql operations to retrieve all users.
// Author: Amanda
func GetAllUsers() ([]authenticate.User, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	// Instantiate user
	var (
		user         authenticate.User
		repid        int
		username     string
		password     string
		firstname    string
		lastname     string
		email        string
		contactno    string
		organisation string
		lastloginDT  time.Time
	)
	users := make([]authenticate.User, 0)

	query := "SELECT RepID, UserName, Password, " +
		"FirstName, LastName, Email, " +
		"ContactNo, Organisation, LastLogin_dt " +
		"FROM Representatives"

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	} else {
		for results.Next() {
			err := results.Scan(&repid, &username, &password, &firstname, &lastname,
				&email, &contactno, &organisation, &lastloginDT)

			if err != nil {
				panic("error getting results from sql select")
			}

			user = authenticate.User{
				RepID:        repid,
				UserName:     username,
				Password:     password,
				FirstName:    firstname,
				LastName:     lastname,
				Email:        email,
				ContactNo:    contactno,
				Organisation: organisation,
				LastLoginDT:  lastloginDT,
			}
			users = append(users, user)
		}
		return users, nil
	}
}

// GetHashPassword implements the sql operations to retrieve user hashed password.
// Author: Amanda
func GetHashPassword(username string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var hashpassword string

	query := "SELECT Password " +
		"FROM Representatives WHERE UserName=?"

	results, err := DB.Query(query, username)
	if err != nil {
		panic("error executing sql select")
	} else {
		if results.Next() {
			err := results.Scan(&hashpassword)

			if err != nil {
				panic("error getting results from sql select")
			}
		} else {
			return hashpassword, errors.New("user not found")
		}

		return hashpassword, nil
	}
}

// AddUser implements the sql operations to insert a new user to the database.
// Author: Amanda
func AddUser(repid int, username string, password string, firstname string, lastname string,
	email string, contactno string, organisation string) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("INSERT INTO Representatives (RepID, UserName, Password, FirstName, LastName, " +
		"Email, ContactNo, Organisation, CurrentLogin_dt, LastLogin_dt, " +
		"CreatedBy, Created_dt, LastModifiedBy, LastModified_DT) VALUES " +
		"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic("error preparing sql insert")
	}

	_, err = stmt.Exec(repid, username, password, firstname, lastname,
		email, contactno, organisation, time.Now(), time.Now(),
		username, time.Now(), username, time.Now())
	if err != nil {
		panic("error executing sql insert")
	}
	return nil
}

// UpdateUser implements the sql operations to update a user from the database.
// Author: Amanda
func UpdateUser(repid int, username string, firstname string, lastname string,
	email string, contactno string, organisation string) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("UPDATE Representatives SET FirstName=?, LastName=?, Email=?, " +
		"ContactNo=?, Organisation=?, LastModifiedBy=?, LastModified_dt=? " +
		"WHERE RepID=?")
	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(firstname, lastname, email,
		contactno, organisation, username, time.Now(),
		repid)
	if err != nil {
		panic("error executing sql update")
	}
	return nil
}

// DeleteUser implements the sql operations to delete a user from the database.
// Author: Amanda
func DeleteUser(userName string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("DELETE FROM Representatives WHERE UserName=?")
	if err != nil {
		panic("error preparing sql delete")
	}

	_, err = stmt.Exec(userName)
	if err != nil {
		panic("error executing sql delete")
	}
	return nil
}

// UpdateLoginDate updates the LastLoginDT to previous CurrentLoginDT.
// Then updates the CurrentLoginDt to time.Now(). No changes to all other information.
// Author: Amanda
func UpdateLoginDate(username string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("UPDATE Representatives SET LastLogin_dt=CurrentLogin_dt, CurrentLogin_dt=? " +
		"WHERE UserName=?")

	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(time.Now(), username)
	if err != nil {
		panic("error executing sql update")
	}
	return nil
}

// UpdatePassword updates the password of a user
// Author: Amanda
func UpdatePassword(repid int, username string, password string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("UPDATE Representatives SET Password=?, LastModifiedBy=?, LastModified_dt=? " +
		"WHERE RepID=?")

	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(password, username, time.Now(), repid)
	if err != nil {
		panic("error executing sql update")
	}
	return nil
}

// UserNameExist checks if UserName exists in the database table
// Author: Amanda
func UserNameExist(username string) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	query := "SELECT UserName " +
		"FROM Representatives WHERE UserName=?"

	results, err := DB.Query(query, username)
	if err != nil {
		panic("error executing sql select")
	} else {
		if results.Next() {
			return true
		}
	}
	return false
}

// EmailExist checks if email exists in the database table
// Author: Amanda
func EmailExist(email string, username string) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var (
		query   string
		results *sql.Rows
		err     error
	)
	if username != "" {
		query := "SELECT Email " +
			"FROM Representatives WHERE UPPER(Email)=UPPER(?) " +
			"AND UserName <> ?"

		results, err = DB.Query(query, email, username)
	} else {
		query = "SELECT Email " +
			"FROM Representatives WHERE UPPER(Email)=UPPER(?)"

		results, err = DB.Query(query, email)
	}

	if err != nil {
		panic("error executing sql select")
	} else {
		if results.Next() {
			return true
		}
	}
	return false
}

// RepIDExist checks if RepID exists in the database table
// Author: Amanda
func RepIDExist(repid string) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	query := "SELECT RepID " +
		"FROM Representatives WHERE RepID=?"

	results, err := DB.Query(query, repid)
	if err != nil {
		panic("error executing sql select")
	} else {
		if results.Next() {
			return true
		}
	}
	return false
}

// IsAdmin checks if user is an admin
// Author: Amanda
func IsAdmin(username string) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	query := "SELECT rep.RepID FROM Representatives rep " +
		"INNER JOIN RepMemberType rmt ON rep.RepID = rmt.RepID " +
		"INNER JOIN MemberType mt ON mt.MemberTypeID = rmt.MemberTypeID " +
		"WHERE UPPER(mt.MemberType) = 'ADMIN' AND rep.UserName = ?"

	results, err := DB.Query(query, username)
	if err != nil {
		panic("error executing sql select")
	} else {
		if results.Next() {
			return true
		}
	}
	return false
}

// GetMemberType implements the sql operations to retrieve the member types.
// Author: Amanda
func GetMemberType() (map[int]authenticate.MemberTypeInfo, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	// Instantiate membertype
	var (
		membertype     = make(map[int]authenticate.MemberTypeInfo)
		membertypeid   int
		membertypedesc string
	)
	query := "SELECT MemberTypeID, MemberType " +
		"FROM MemberType "

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	} else {
		checked := ""
		for results.Next() {
			err := results.Scan(&membertypeid, &membertypedesc)
			if err != nil {
				panic("error getting results from sql select")
			}

			checked = "checked"
			if strings.ToUpper(membertypedesc) == "ADMIN" {
				checked = "disabled"
			}
			membertype[membertypeid] = authenticate.MemberTypeInfo{
				MemberType: membertypedesc,
				Checked:    checked}
		}

		return membertype, nil
	}
}

// GetRepMemberType implements the sql operations to retrieve the RepMemberType.
// Author: Amanda
func GetRepMemberType(repid int) (map[int]authenticate.MemberTypeInfo, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	// Instantiate membertype
	var (
		membertype     = make(map[int]authenticate.MemberTypeInfo)
		membertypeid   int
		membertypedesc string
		repid1         int
	)
	query := "SELECT mt.MemberTypeID, mt.MemberType, " +
		"ifnull((SELECT rmt.RepID FROM RepMemberType rmt WHERE rmt.RepID = ? " +
		"AND rmt.MemberTypeID = mt.MemberTypeID), 0) RepID FROM MemberType mt"

	results, err := DB.Query(query, repid)
	if err != nil {
		panic("error executing sql select")
	} else {

		for results.Next() {
			err := results.Scan(&membertypeid, &membertypedesc, &repid1)
			if err != nil {
				panic("error getting results from sql select")
			}

			checked := ""
			disabled := ""
			if repid1 != 0 {
				checked = "checked"
			}

			if strings.ToUpper(membertypedesc) == "ADMIN" {
				disabled = "disabled"
			}

			membertype[membertypeid] = authenticate.MemberTypeInfo{
				MemberType: membertypedesc,
				Checked:    checked,
				Disabled:   disabled,
			}
		}

		return membertype, nil
	}
}

// AddRepMemberType implements the sql operations to insert a new RepMemberType record to the database.
// Author: Amanda
func AddRepMemberType(repid int, membertypeid int, username string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("INSERT INTO RepMemberType (RepID, MemberTypeID, " +
		"CreatedBy, Created_dt, LastModifiedBy, LastModified_DT) VALUES " +
		"(?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic("error preparing sql insert")
	}

	_, err = stmt.Exec(repid, membertypeid, username, time.Now(), username, time.Now())
	if err != nil {
		panic("error executing sql insert")
	}
	return nil
}

// DeleteRepMemberType implements the sql operations to delete membertypes by repid from the database.
// Author: Amanda
func DeleteRepMemberType(repid int) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	// Do not delete admin record
	stmt, err := DB.Prepare("DELETE FROM RepMemberType rmt " +
		"WHERE RepID=? AND MemberTypeID NOT IN " +
		"(SELECT MemberTypeID FROM MemberType " +
		"WHERE UPPER(MemberType) = 'ADMIN')")
	if err != nil {
		panic("error preparing sql delete")
	}

	_, err = stmt.Exec(repid)
	if err != nil {
		panic("error executing sql delete")
	}
	return nil
}
