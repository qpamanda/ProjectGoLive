/*
Author: Huang Yanping
Last Updated: 26-Jun-2021
*/
package database

import (
	"fmt"
	"time"
)

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

func GetMyRecipient(RepID int) ([]Recipient, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var recipients []Recipient

	results, err := DB.Query("SELECT RecipientID, RepID_FK, Name, Category, Profile,"+
		"ContactNo FROM Recipients WHERE RepID_FK = ?", RepID)

	if err != nil {
		return nil, err
	} else {
		for results.Next() {
			var recipient Recipient
			err := results.Scan(&recipient.RecipientID, &recipient.RepID, &recipient.Name,
				&recipient.Category, &recipient.Profile, &recipient.ContactNo)
			if err != nil {
				return nil, err
			}
			recipients = append(recipients, recipient)
		}
		return recipients, nil
	}
}

func AddRecipient(RepID int, Name string, Category bool, Profile string, ContactNo string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("INSERT INTO Recipients (RepID_FK, Name, Category, Profile, ContactNo, CreatedDT, LastModifiedDT) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(RepID, Name, Category, Profile, ContactNo, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}
func GetRecipient(RepID int, RecipientID int64) (Recipient, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	var recipient Recipient

	results, err := DB.Query("SELECT RecipientID, RepID_FK, Name, Category, Profile,"+
		"ContactNo FROM Recipients WHERE RepID_FK = ? AND RecipientID = ?", RepID, RecipientID)

	if err != nil {
		return recipient, err
	} else {
		for results.Next() {
			err := results.Scan(&recipient.RecipientID, &recipient.RepID, &recipient.Name,
				&recipient.Category, &recipient.Profile, &recipient.ContactNo)
			if err != nil {
				return recipient, err
			}
		}
		return recipient, nil
	}
}

func UpdateRecipient(RepID int, RecipientID int64, Name string, Category bool, Profile string, ContactNo string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("UPDATE Recipients SET Name=?,Category=?, Profile=?, ContactNo=?, LastModifiedDT=? WHERE RepID_FK=? AND RecipientID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(Name, Category, Profile, ContactNo, time.Now(), RepID, RecipientID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRecipient(RepID int, RecipientID int64) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	stmt, err := DB.Prepare("DELETE FROM Recipients WHERE RepID_FK=? AND RecipientID=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(RepID, RecipientID)
	if err != nil {
		return err
	}
	return nil
}