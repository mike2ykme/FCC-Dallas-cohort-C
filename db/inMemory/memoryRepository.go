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
	users                  userMap
	currentHighestUserId   uint
	decks                  deckMap
	currentHighestDeckId   uint
	flashcards             flashcardMap
	currentHighestCardId   uint
	answers                answerMap
	currentHighestAnswerId uint
}
type userMap map[uint]*models.User
type deckMap map[uint]*models.Deck
type flashcardMap map[uint]*models.FlashCard
type answerMap map[uint]*models.Answer

func NewInMemoryRepository() *repository {
	return &repository{
		users:                  make(userMap),
		currentHighestUserId:   0,
		decks:                  make(deckMap),
		currentHighestDeckId:   0,
		flashcards:             make(flashcardMap),
		currentHighestCardId:   0,
		answers:                make(answerMap),
		currentHighestAnswerId: 0,
	}
}
