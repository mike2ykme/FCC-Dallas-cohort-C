package models

type User struct {
	Username  string `json:"username"`
	SubId     string `json:"subId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
