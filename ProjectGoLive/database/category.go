package database

import "errors"

// Author: Ahmad Bahrudin
// category struct for storing category information
type category struct {
	CategoryID      int    // To store category's id
	Category        string // To store category's name
	CreatedBy       string // To store category's created by
	Created_dt      string // To store category's created date/time
	LastModifiedBy  string // To store category's last modified by
	LastModified_dt string // To store category's last modified date/time
}

var (
	Cat = &category{0, "", "", "", "", ""}
)

func (*category) Insert(categoryID int, category string, userName string) error {
	stmt, err := DB.Prepare("INSERT INTO Category " +
		"(CategoryID, Category, CreatedBy, Created_dt, LastModifiedBy, LastModified_dt) " +
		"VALUES (?, ?, ?, now(), ?, now())")

	if err != nil {
		panic("error preparing sql insert")
	}

	_, err = stmt.Exec(categoryID, category, userName, userName)
	if err != nil {
		panic("error executing sql insert")
	}

	return nil
}

func (*category) Update(categoryID int, category string, userName string) error {
	stmt, err := DB.Prepare("UPDATE Category " +
		"SET CategoryID=?, Category=?, LastModifiedBy=?, LastModified_dt=now() " +
		"WHERE CategoryID=?")

	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(categoryID, category, userName, categoryID)
	if err != nil {
		panic("error executing sql update")
	}

	return nil
}

func (cat category) Get(categoryID int) (category, error) {
	query := "SELECT * " +
		"FROM Category " +
		"WHERE categoryID=?"

	results, err := DB.Query(query, categoryID)
	if err != nil {
		panic("error executing sql select")
	}

	if results.Next() {
		err := results.Scan(&cat.CategoryID, &cat.Category, &cat.CreatedBy, &cat.Created_dt, &cat.LastModifiedBy, &cat.LastModified_dt)

		if err != nil {
			panic("error getting results from sql select")
		}
	}
	return cat, errors.New("category not found")
}

func (cat category) GetAll() (map[int]category, error) {
	query := "SELECT * " +
		"FROM Category "

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	}

	c := make(map[int]category)
	for results.Next() {
		err := results.Scan(&cat.CategoryID, &cat.Category, &cat.CreatedBy, &cat.Created_dt, &cat.LastModifiedBy, &cat.LastModified_dt)

		if err != nil {
			panic("error getting results from sql select")
		}

		c[cat.CategoryID] = cat
	}
	return c, nil
}

func (cat category) Delete(categoryID int) error {
	stmt, err := DB.Prepare("DELETE FROM Category " +
		"WHERE categoryID=?")

	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(categoryID)
	if err != nil {
		panic("error executing sql update")
	}

	return nil
}
