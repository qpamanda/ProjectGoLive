package database

import "strings"

// Author: Ahmad Bahrudin.
// Category struct for storing category information
type Category struct {
	CategoryID      int    // To store category's id
	Category        string // To store category's name
	CreatedBy       string // To store category's created by
	Created_dt      string // To store category's created date/time
	LastModifiedBy  string // To store category's last modified by
	LastModified_dt string // To store category's last modified date/time
}

var (
	Cat = &Category{0, "", "", "", "", ""}
)

// Author: Ahmad Bahrudin.
// GetNextID function that get next category's id from the database
func (cat Category) GetNextID() (int, error) {
	query := "SELECT CategoryID " +
		"FROM Category"

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	}

	for results.Next() {
		err := results.Scan(&cat.CategoryID)

		if err != nil {
			panic("error getting results from sql select")
		}
	}
	return (cat.CategoryID + 1), nil
}

// Author: Ahmad Bahrudin.
// Insert function that adds new category to the database
func (*Category) Insert(categoryID int, category string, userName string) error {
	stmt, err := DB.Prepare("INSERT INTO Category " +
		"(CategoryID, Category, CreatedBy, Created_dt, LastModifiedBy, LastModified_dt) " +
		"VALUES (?, ?, ?, DATE_ADD(NOW(), INTERVAL 8 HOUR), ?, DATE_ADD(NOW(), INTERVAL 8 HOUR))")

	if err != nil {
		panic("error preparing sql insert")
	}

	_, err = stmt.Exec(categoryID, category, userName, userName)
	if err != nil {
		panic("error executing sql insert")
	}

	return nil
}

// Author: Ahmad Bahrudin.
// Update function that update exist category from the database
func (*Category) Update(categoryID int, category string, userName string) error {
	stmt, err := DB.Prepare("UPDATE Category " +
		"SET Category=?, LastModifiedBy=?, LastModified_dt=DATE_ADD(NOW(), INTERVAL 8 HOUR) " +
		"WHERE CategoryID=?")

	if err != nil {
		panic("error preparing sql update")
	}

	_, err = stmt.Exec(category, userName, categoryID)
	if err != nil {
		panic("error executing sql update")
	}

	return nil
}

// Author: Ahmad Bahrudin.
// GetAll function that get all category from the database
func (cat Category) GetAll() (map[int]Category, error) {
	query := "SELECT * " +
		"FROM Category "

	results, err := DB.Query(query)
	if err != nil {
		panic("error executing sql select")
	}

	c := make(map[int]Category)
	for results.Next() {
		err := results.Scan(&cat.CategoryID, &cat.Category, &cat.CreatedBy, &cat.Created_dt, &cat.LastModifiedBy, &cat.LastModified_dt)

		if err != nil {
			panic("error getting results from sql select")
		}
		cat.Created_dt = strings.Replace(cat.Created_dt, "T", " ", -1)
		cat.Created_dt = strings.Replace(cat.Created_dt, "Z", " ", -1)
		cat.LastModified_dt = strings.Replace(cat.LastModified_dt, "T", " ", -1)
		cat.LastModified_dt = strings.Replace(cat.LastModified_dt, "Z", " ", -1)

		c[cat.CategoryID] = cat
	}
	return c, nil
}

// Author: Ahmad Bahrudin.
// Delete function that delete category from the database
// IMPORTANT: ONLY USE FOR TESTING
func (cat Category) Delete(categoryID int) error {
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
