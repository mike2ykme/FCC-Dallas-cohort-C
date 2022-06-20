package inMemory

import (
	"errors"
	"fmt"
	"teamC/models"
)

func (r *repository) GetAllUsersString() string {
	return fmt.Sprintf("%#v", r.users)
}

func (r *repository) SaveUser(user *models.User) (uint, error) {
	//m.currentHighestUserId
	if user.ID == 0 {
		user.ID = r.currentHighestUserId
		r.currentHighestUserId++
	} else if user.ID > r.currentHighestUserId {
		r.currentHighestUserId = user.ID + 1
	}
	r.users[user.ID] = user

	return user.ID, nil
}

func (r *repository) GetUserById(uRef *models.User, id uint) error {
	if val, ok := r.users[id]; ok {
		uRef.CopyReferences(val)
		return nil
	}

	return errors.New("unable to find user")
}
func (r *repository) GetUserByUsername(uRef *models.User, username string) error {
	for _, val := range r.users {
		if val.Username == username {
			uRef.CopyReferences(val)
			return nil
		}
	}
	return errors.New("unable to find user")
}
func (r *repository) GetUserBySubId(uRef *models.User, subId string) error {
	for _, val := range r.users {
		if val.SubId == subId {
			uRef.CopyReferences(val)
			return nil
		}
	}
	return errors.New("unable to find user")
}
func (r *repository) GetAllUsers(usersRef *[]models.User) error {
	for _, user := range r.users {
		var newUser models.User
		newUser.CopyReferences(user)
		*usersRef = append(*usersRef, newUser)
	}
	return nil
}

func (r *repository) GetUsernameById(id uint) (string, error) {
	temp := models.User{}
	err := r.GetUserById(&temp, id)
	return temp.Username, err
}
