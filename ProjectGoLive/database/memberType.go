package database

import "strings"

// Author: Ahmad Bahrudin.
// MemberType struct for storing member type information
type MemberType struct {
	MemberTypeID    int    // To store member type's id
	MemberType      string // To store member type's name
	CreatedBy       string // To store member type's created by
	Created_dt      string // To store member type's created date/time
	LastModifiedBy  string // To store member type's last modified by
	LastModified_dt string // To store member type's last modified date/time
}

var (
	MemT = &MemberType{0, "", "", "", "", ""}
)

// Author: Ahmad Bahrudin.
// GetNextID function that get next member type's id from the database
func (memT MemberType) GetNextID() (int, error) {
	query := "SELECT MemberTypeID " +
		"FROM MemberType"

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	}

	for results.Next() {
		err := results.Scan(&memT.MemberTypeID)

		if err != nil {
			panic("error getting results from sql select")
		}
	}
	return (memT.MemberTypeID + 1), nil
}

// Author: Ahmad Bahrudin.
// Insert function that adds new member type to the database
func (*MemberType) Insert(memberTypeID int, memberType string, userName string) error {
	stmt, err := DB.Prepare("INSERT INTO MemberType " +
		"(MemberTypeID, MemberType, CreatedBy, Created_dt, LastModifiedBy, LastModified_dt) " +
		"VALUES (?, ?, ?, DATE_ADD(NOW(), INTERVAL 8 HOUR), ?, DATE_ADD(NOW(), INTERVAL 8 HOUR))")

	if err != nil {
		panic("error preparing sql insert")
	}

	_, err = stmt.Exec(memberTypeID, memberType, userName, userName)
	if err != nil {
		panic("error executing sql insert")
	}

	return nil
}

// Author: Ahmad Bahrudin.
// Update function that update exist member type from the database
func (*MemberType) Update(memberTypeID int, memberType string, userName string) error {
	stmt, err := DB.Prepare("UPDATE MemberType " +
		"SET MemberType=?, LastModifiedBy=?, LastModified_dt=DATE_ADD(NOW(), INTERVAL 8 HOUR) " +
		"WHERE MemberTypeID=?")

	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(memberType, userName, memberTypeID)
	if err != nil {
		panic("error executing sql update")
	}

	return nil
}

// Author: Ahmad Bahrudin.
// GetAll function that get all member type from the database
func (memT MemberType) GetAll() (map[int]MemberType, error) {
	query := "SELECT * " +
		"FROM MemberType "

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	}

	m := make(map[int]MemberType)
	for results.Next() {
		err := results.Scan(&memT.MemberTypeID, &memT.MemberType, &memT.CreatedBy, &memT.Created_dt, &memT.LastModifiedBy, &memT.LastModified_dt)

		if err != nil {
			panic("error getting results from sql select")
		}
		memT.Created_dt = strings.Replace(memT.Created_dt, "T", " ", -1)
		memT.Created_dt = strings.Replace(memT.Created_dt, "Z", " ", -1)
		memT.LastModified_dt = strings.Replace(memT.LastModified_dt, "T", " ", -1)
		memT.LastModified_dt = strings.Replace(memT.LastModified_dt, "Z", " ", -1)

		m[memT.MemberTypeID] = memT
	}
	return m, nil
}

// Author: Ahmad Bahrudin.
// Delete function that delete member type from the database
// IMPORTANT: ONLY USE FOR TESTING
func (memT MemberType) Delete(categoryID int) error {
	stmt, err := DB.Prepare("DELETE FROM MemberType " +
		"WHERE MemberTypeID=?")

	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(categoryID)
	if err != nil {
		panic("error executing sql update")
	}

	return nil
}
