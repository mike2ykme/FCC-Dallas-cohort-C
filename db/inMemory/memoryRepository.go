package inMemory

import (
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
	users                  map[uint]*models.User
	currentHighestUserId   uint
	decks                  map[uint]*models.Deck
	currentHighestDeckId   uint
	cards                  map[uint]*models.FlashCard
	currentHighestCardId   uint
	answers                map[uint]*models.Answer
	currentHighestAnswerId uint
}

func NewInMemoryRepository() *repository {
	return &repository{
		users:                make(userMap),
		currentHighestUserId: 0,
	}
}
