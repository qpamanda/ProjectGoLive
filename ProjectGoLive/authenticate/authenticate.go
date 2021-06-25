package authenticate

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"unicode"
)

var (
	MapSessions = map[string]string{}

	MinUserName int // Set the min length for new Username
	MaxUserName int // Set the max length for new Username
	MinPassword int // Set the min length for new Password
	MaxPassword int // Set the max length for new Password

	IsAdmin = false
)

// User struct for storing user account information
type User struct {
	RepID        int
	UserName     string
	Password     string
	FirstName    string
	LastName     string
	Email        string
	ContactNo    string
	Organisation string
	LastLoginDT  time.Time
}

// MemberType struct for storing Member Type info
type MemberTypeInfo struct {
	MemberType string
	Checked    string
	Disabled   string
}

// validateUserInput func validates user input
// Author: Amanda
func ValidateUserInput(adduser bool, username string, password string, cmfpassword string,
	firstname string, lastname string, email string,
	contactno string, organisation string) error {

	// Validate username
	if username == "" {
		return errors.New("username cannot be blank")
	} else if len(username) < MinUserName || len(username) > MaxUserName {
		return fmt.Errorf("username must be between %d - %d characters", MinUserName, MaxUserName)
	}

	// Validate password if request to add user
	if adduser {
		// Validate password
		if err := ValidatePassword(password, cmfpassword); err != nil {
			return err
		}
	}

	// Validate first name
	if firstname == "" {
		return errors.New("first name cannot be blank")
	}

	// Validate last name
	if lastname == "" {
		return errors.New("last name cannot be blank")
	}

	// Validate email
	if email == "" {
		return errors.New("email cannot be blank")
	} else if !IsValidEmail(email) {
		return errors.New("invalid email")
	}

	// Validate contact no
	if contactno == "" {
		return errors.New("contact no cannot be blank")
	}

	return nil
}

// validatePassword validates that the input user password must contain as least
// one upper case, lower case, numeric and special characters.
// Author: Amanda
func ValidatePassword(password string, cmfpassword string) error {
	if password == "" {
		return errors.New("password cannot be blank")
	} else if len(password) < MinPassword || len(password) > MaxPassword {
		return fmt.Errorf("password must be between %d - %d characters", MinPassword, MaxPassword)
	}

next:
	for name, classes := range map[string][]*unicode.RangeTable{
		"upper case": {unicode.Upper, unicode.Title},
		"lower case": {unicode.Lower},
		"numeric":    {unicode.Number, unicode.Digit},
		"special":    {unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
	} {
		for _, r := range password {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return fmt.Errorf("password must have at least one %s character", name)
	}

	// Validate confirm password
	if cmfpassword == "" {
		return errors.New("confirm password cannot be blank")
	} else if cmfpassword != password {
		return errors.New("confirm password must be the same as password")
	}

	return nil
}

// isValidEmail validates if the string parameter is a valid email using regexp
// Author: Amanda
func IsValidEmail(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
