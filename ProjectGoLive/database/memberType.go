package database

import "errors"

// Author: Ahmad Bahrudin
// memberType struct for storing member type information
type memberType struct {
	MemberTypeID    int    // To store memberType's id
	MemberType      string // To store memberType's name
	CreatedBy       string // To store memberType's created by
	Created_dt      string // To store memberType's created date/time
	LastModifiedBy  string // To store memberType's last modified by
	LastModified_dt string // To store memberType's last modified date/time
}

var (
	MemT = &memberType{0, "", "", "", "", ""}
)

func (*memberType) Insert(memberTypeID int, memberType string, userName string) error {
	stmt, err := DB.Prepare("INSERT INTO MemberType " +
		"(MemberTypeID, MemberType, CreatedBy, Created_dt, LastModifiedBy, LastModified_dt) " +
		"VALUES (?, ?, ?, now(), ?, now())")

	if err != nil {
		panic("error preparing sql insert")
	}

	_, err = stmt.Exec(memberTypeID, memberType, userName, userName)
	if err != nil {
		panic("error executing sql insert")
	}

	return nil
}

func (*memberType) Update(memberTypeID int, memberType string, userName string) error {
	stmt, err := DB.Prepare("UPDATE MemberType " +
		"SET MemberTypeID=?, MemberType=?, LastModifiedBy=?, LastModified_dt=now() " +
		"WHERE MemberTypeID=?")

	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(memberTypeID, memberType, userName, memberTypeID)
	if err != nil {
		panic("error executing sql update")
	}

	return nil
}

func (memT memberType) Get(memberTypeID int) (memberType, error) {
	query := "SELECT * " +
		"FROM MemberType " +
		"WHERE MemberTypeID=?"

	results, err := DB.Query(query, memberTypeID)
	if err != nil {
		panic("error executing sql select")
	}

	if results.Next() {
		err := results.Scan(&memT.MemberTypeID, &memT.MemberType, &memT.CreatedBy, &memT.Created_dt, &memT.LastModifiedBy, &memT.LastModified_dt)

		if err != nil {
			panic("error getting results from sql select")
		}
	}
	return memT, errors.New("member type not found")
}

func (memT memberType) GetAll() (map[int]memberType, error) {
	query := "SELECT * " +
		"FROM MemberType "

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	}

	m := make(map[int]memberType)
	for results.Next() {
		err := results.Scan(&memT.MemberTypeID, &memT.MemberType, &memT.CreatedBy, &memT.Created_dt, &memT.LastModifiedBy, &memT.LastModified_dt)

		if err != nil {
			panic("error getting results from sql select")
		}

		m[memT.MemberTypeID] = memT
	}
	return m, nil
}

func (memT memberType) Delete(categoryID int) error {
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
