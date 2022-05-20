package inMemory

import (
	"teamC/models"
	"testing"
)

/*
type UserRepository interface {
	SaveUser(*models.User) (uint, error) // Done
	GetUserById(*models.User, uint) error // Done
	GetUserByUsername(*models.User, string) error //Done
	GetUserBySubId(*models.User, string) error // Done
	GetAllUsers(*[]models.User) error // Done
}
*/

func TestSaveTwoAndGetAllUsers(t *testing.T) {
	repo := NewInMemoryRepository()

	repo.SaveUser(&models.User{
		Id:        1,
		Username:  "1-Username",
		SubId:     "1-SubscriberId",
		FirstName: "1-FirstName",
		LastName:  "1-LastName",
	})
	repo.SaveUser(&models.User{
		Id:        2,
		Username:  "2-Username",
		SubId:     "2-SubscriberId",
		FirstName: "2-FirstName",
		LastName:  "2-LastName",
	})

	allUsers := make([]models.User, 0)
	err := repo.GetAllUsers(&allUsers)

	if len(allUsers) != 2 || err != nil || allUsers[0].Id != 1 || allUsers[1].Id != 2 {
		t.Fatalf("We got %#v, but we wanted 2 users and Ids of 1,2", allUsers)
	}
}
func TestSaveAndGetUserBySubId(t *testing.T) {
	repo := NewInMemoryRepository()
	subId := "SubscriberID"

	repo.SaveUser(&models.User{
		Id:        1,
		Username:  "Username",
		SubId:     subId,
		FirstName: "FirstName",
		LastName:  "LastName",
	})
	var testUser models.User

	err := repo.GetUserBySubId(&testUser, subId)

	if testUser.SubId != subId || err != nil {
		t.Fatalf("User is %#v, and we wanted a user with an subId == %s", testUser, subId)
	}
}
func TestSaveAndGetUserByUsername(t *testing.T) {
	repo := NewInMemoryRepository()
	username := "Username"

	repo.SaveUser(&models.User{
		Id:        1,
		Username:  username,
		SubId:     "SubscriberID",
		FirstName: "FirstName",
		LastName:  "LastName",
	})
	var testUser models.User

	err := repo.GetUserByUsername(&testUser, username)

	if testUser.Username != username || err != nil {
		t.Fatalf("User is %#v, and we wanted a user with an username == %s", testUser, username)
	}
}
func TestSaveUser(t *testing.T) {
	repo := NewInMemoryRepository()

	repo.SaveUser(&models.User{
		Id:        1,
		Username:  "Username",
		SubId:     "SubscriberID",
		FirstName: "FirstName",
		LastName:  "LastName",
	})

	if len(repo.users) != 1 {
		t.Fatal("User was not stored in slice")
	}
}

func TestSaveAndGetUserById(t *testing.T) {
	repo := NewInMemoryRepository()
	repo.SaveUser(&models.User{
		Id:        1,
		Username:  "Username",
		SubId:     "SubscriberID",
		FirstName: "FirstName",
		LastName:  "LastName",
	})
	var testUser models.User
	err := repo.GetUserById(&testUser, 1)

	if testUser.Id != 1 || err != nil {
		t.Fatalf("User is %#v, and we wanted a user with an Id == 1", testUser)
	}

}
