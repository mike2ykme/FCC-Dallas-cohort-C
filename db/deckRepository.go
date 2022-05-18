package db

import "teamC/models"

type DeckRepository interface {
	SaveDeck(*models.Deck) (uint, error)
	GetDeckById(*models.Deck, uint) error
	GetAllDecks(*models.Deck) error
}
