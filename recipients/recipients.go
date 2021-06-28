package recipients

import "time"

type Request struct {
	RequestID       int
	RepID           int
	CatID           int
	RecID           int
	StatusCode      int
	ReqDesc         string
	ToCompleteBy    time.Time
	FulfilAt        string
	UserName        string
	FirstName       string
	LastName        string
	Email           string
	ContactNo       string
	Organisation    string
	Category        string
	RecName         string
	RecCategory     int
	RecProfile      string
	RecContactNo    string
	Status          string
	RecCategoryDesc string
}
