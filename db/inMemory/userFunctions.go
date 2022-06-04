package inMemory

import (
	"fmt"
	"teamC/models"
)

func (m *repository) GetAllUsersString() string {
	return fmt.Sprintf("%#v", m.users)
}

func (m *repository) SaveUser(user *models.User) (uint, error) {
	//m.currentHighestUserId
	if user.ID == 0 {
		user.ID = m.currentHighestUserId
		m.currentHighestUserId++
	} else if user.ID > m.currentHighestUserId {
		m.currentHighestUserId = user.ID + 1
	}
	m.users[user.ID] = user

	return user.ID, nil
}

func (m *repository) GetUserById(uRef *models.User, id uint) error {
	if val, ok := m.users[id]; ok {
		uRef.CopyReferences(val)
	}

	return nil
}
func (m *repository) GetUserByUsername(uRef *models.User, username string) error {
	for _, val := range m.users {
		if val.Username == username {
			uRef.CopyReferences(val)
		}
	}
	return nil
}
func (m *repository) GetUserBySubId(uRef *models.User, subId string) error {
	for _, val := range m.users {
		if val.SubId == subId {
			uRef.CopyReferences(val)
		}
	}
	return nil
}
func (m *repository) GetAllUsers(usersRef *[]models.User) error {
	for _, user := range m.users {
		var newUser models.User
		newUser.CopyReferences(user)
		*usersRef = append(*usersRef, newUser)
	}
	return nil
}
