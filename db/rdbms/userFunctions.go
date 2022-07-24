package rdbms

import (
	"errors"
	"fmt"
	"teamC/models"
)

func (r *repository) GetAllUsersString() string {
	var users = make([]models.User, 0)
	r.DB.Find(&users)

	return fmt.Sprintf("%#v", users)
}

func (r *repository) SaveUser(user *models.User) (uint, error) {
	err := r.DB.Save(&user).Error
	return user.ID, err
}

func (r *repository) GetUserById(uRef *models.User, id uint) error {
	if id <= 0 {
		return errors.New("user ID cannot be <= 0")
	}
	return r.DB.Take(uRef, id).Error
}

func (r *repository) GetUserByUsername(uRef *models.User, username string) error {
	if len(username) == 0 {
		return errors.New("username cannot be empty")
	}
	return r.DB.Where("username = ?", username).First(&uRef).Error
}
func (r *repository) GetUserBySubId(uRef *models.User, subId string) error {
	if len(subId) == 0 {
		return errors.New("sub ID cannot be empty")
	}
	return r.DB.Where("sub_id = ?", subId).First(&uRef).Error
}
func (r *repository) GetAllUsers(usersRef *[]models.User) error {
	return r.DB.Find(usersRef).Error
}
func (r *repository) GetUsernameById(id uint) (string, error) {
	temp := models.User{}
	err := r.GetUserById(&temp, id)
	return temp.Username, err
}
