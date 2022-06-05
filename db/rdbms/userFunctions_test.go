package rdbms

import (
	"gorm.io/gorm"
	"teamC/models"
	"testing"
)

func TestSaveTwoAndGetAllUsers(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.User{}, "1=1")

	idA, err := repo.SaveUser(&models.User{
		//Model:     gorm.Model{ID: 1},
		Username:  "1-Username",
		SubId:     "1-SubscriberId",
		FirstName: "1-FirstName",
		LastName:  "1-LastName",
	})
	if err != nil {
		t.Fatal(err)
	}
	idB, err := repo.SaveUser(&models.User{
		//Model:     gorm.Model{ID: 2},
		Username:  "2-Username",
		SubId:     "2-SubscriberId",
		FirstName: "2-FirstName",
		LastName:  "2-LastName",
	})
	if err != nil {
		t.Fatal(err)
	}
	allUsers := make([]models.User, 0)
	err = repo.GetAllUsers(&allUsers)

	if len(allUsers) != 2 || err != nil || allUsers[0].ID != idA || allUsers[1].ID != idB {
		t.Fatalf("We got %#v, but we wanted 2 users and Ids of 1,2", allUsers)
	}
}
func TestSaveAndGetUserBySubId(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.User{}, "1=1")

	subId := "SubscriberID"

	repo.SaveUser(&models.User{
		//Model:     gorm.Model{ID: 1},
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
	repo := getRepo(t)
	repo.DB.Delete(&models.User{}, "1=1")

	username := "Username"

	_, err := repo.SaveUser(&models.User{
		//Model:     gorm.Model{ID: 1},
		Username:  username,
		SubId:     "SubscriberID",
		FirstName: "FirstName",
		LastName:  "LastName",
	})
	if err != nil {
		t.Fatal(err)
	}
	var testUser models.User

	err = repo.GetUserByUsername(&testUser, username)

	if testUser.Username != username || err != nil {
		t.Fatalf("User is %#v, and we wanted a user with an username == %s", testUser, username)
	}
}
func TestSaveUser(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.User{}, "1=1")

	repo.SaveUser(&models.User{
		Model:     gorm.Model{ID: 1},
		Username:  "Username",
		SubId:     "SubscriberID",
		FirstName: "FirstName",
		LastName:  "LastName",
	})

	var users []models.User
	repo.GetAllUsers(&users)
	if len(users) >= 1 {
		t.Fatal("User was not stored in slice")
	}
}
func TestSaveAndGetUserById(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.User{}, "1=1")

	idA, err := repo.SaveUser(&models.User{
		//Model:     gorm.Model{ID: 1},
		Username:  "Username",
		SubId:     "SubscriberID",
		FirstName: "FirstName",
		LastName:  "LastName",
	})
	if err != nil {
		t.Fatal(err)
	}
	var testUser models.User
	err = repo.GetUserById(&testUser, idA)

	if testUser.ID != idA || err != nil {
		t.Fatalf("User is %#v, and we wanted a user with an ID == 1", testUser)
	}

}
