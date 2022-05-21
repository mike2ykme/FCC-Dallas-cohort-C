package inMemory

import (
	"errors"
	"teamC/models"
)

func (m *repository) SaveDeck(deck *models.Deck) (uint, error) {
	if deck.OwnerId == 0 {
		return 0, errors.New("deck must have an owwner")
	}
	if deck.Id == 0 {
		deck.Id = m.currentHighestDeckId
		m.currentHighestDeckId++
	} else if deck.Id > m.currentHighestDeckId {
		m.currentHighestDeckId = deck.Id + 1
	}
	copy := deck.Copy()
	copy.FlashCards = nil

	for i := 0; i < len(deck.FlashCards); i++ {
		deck.FlashCards[i].DeckId = deck.Id
		if _, err := m.SaveFlashcard(&deck.FlashCards[i]); err != nil {
			return 0, err
		}
	}
	m.decks[deck.Id] = &copy

	return deck.Id, nil
}

func (m *repository) GetDeckById(d *models.Deck, id uint) error {
	if val, ok := m.decks[id]; ok {
		if val.Id == id {
			m.GetAllFlashcardByDeckId(&val.FlashCards, val.Id)
			d.CopyReferences(val)
		}
	}
	return nil
}

func (m *repository) GetAllDecks(userDeck *[]models.Deck) error {
	for id, deck := range m.decks {
		var newDeck models.Deck
		m.GetAllFlashcardByDeckId(&newDeck.FlashCards, id)
		newDeck.CopyReferences(deck)

		*userDeck = append(*userDeck, newDeck)
	}
	return nil
}
