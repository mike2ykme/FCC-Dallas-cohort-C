package models

type User struct {
	Id        uint   `json:"id"`
	Username  string `json:"username"`
	SubId     string `json:"subId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (uRef *User) CopyReferences(val *User) {
	uRef.SubId = val.SubId
	uRef.Id = val.Id
	uRef.FirstName = val.FirstName
	uRef.LastName = val.LastName
	uRef.Username = val.Username
}
