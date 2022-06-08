package models

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//ID       uint   `json:"id"`
	Username  string `json:"username"`
	SubId     string `json:"subId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (uRef *User) CopyReferences(val *User) {
	uRef.SubId = val.SubId
	//uRef.ID = val.ID
	uRef.Model = val.Model
	uRef.FirstName = val.FirstName
	uRef.LastName = val.LastName
	uRef.Username = val.Username
}

func (uRef User) GetUserFromJWT(token *jwt.Token) User {
	claims := token.Claims.(jwt.MapClaims)
	id := claims["id"].(uint)
	username, _ := claims["username"].(string)
	//subId, _ := claims["subId"].(string)
	first, _ := claims["firstName"].(string)
	last, _ := claims["lastName"].(string)

	return User{
		//ID:        id,
		Model:     gorm.Model{ID: id},
		Username:  username,
		FirstName: first,
		LastName:  last,
	}
	//claims := user.Claims.(jwt.MapClaims)
}
