package inMemory

import (
	"log"
	"teamC/models"
)

/*
type UserRepository interface {
	SaveUser(*models.User) (uint, error)
	GetUserById(*models.User, uint) error
	GetUserByUsername(*models.User, string) error
	GetUserBySubId(*models.User, string) error
	GetAllUsers(*[]models.User) error
}
*/
type repository struct {
	users          map[uint]*models.User
	currentHighest uint
}

type userMap map[uint]*models.User

func NewInMemoryRepository() *repository {
	return &repository{
		users:          make(userMap),
		currentHighest: 0,
	}
}
func (m *repository) PrintAllUsers() {
	log.Println(m.users)
}

func (m *repository) SaveUser(user *models.User) (uint, error) {
	//m.currentHighest
	if user.Id == 0 {
		user.Id = m.currentHighest
		m.currentHighest++
	} else if user.Id > m.currentHighest {
		m.currentHighest = user.Id + 1
	}
	m.users[user.Id] = user

	return user.Id, nil
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
