package db

import "teamC/models"

type FlashcardRepository interface {
	SaveFlashcard(*models.FlashCard) (uint, error)
	GetFlashcardById(*models.FlashCard, uint) error
	GetFlashcardsByDeck(*models.FlashCard) error
	GetAllFlashcards(*models.FlashCard) error
}
