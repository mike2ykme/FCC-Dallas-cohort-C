package db

import "teamC/models"

type FlashcardRepository interface {
	SaveFlashcard(*models.FlashCard) (uint, error)
	GetFlashcardById(*models.FlashCard, uint) error
	GetAllFlashcardByDeckId(*[]models.FlashCard, uint) error
	GetAllFlashcards(*[]models.FlashCard) error
}
