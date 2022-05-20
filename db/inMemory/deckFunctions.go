package inMemory

import "teamC/models"

func (m *repository) SaveDeck(deck *models.Deck) (uint, error) {
	if deck.Id == 0 {
		deck.Id = m.currentHighestDeckId
		m.currentHighestDeckId++
	} else if deck.Id > m.currentHighestDeckId {
		m.currentHighestDeckId = deck.Id + 1
	}
	copy := deck.Copy()
	for _, card := range deck.Cards {
		card.DeckId = deck.Id
		if _, err := m.SaveFlashcard(&card); err != nil {
			return 0, err
		}
	}
	m.decks[deck.Id] = &copy

	return deck.Id, nil
}

func (m *repository) GetDeckById(d *models.Deck, id uint) error {
	if val, ok := m.decks[d.Id]; ok {
		if val.Id == id {
			m.GetAllFlashcardByDeckId(&val.Cards, val.Id)
			d.CopyReferences(val)
		}
	}
	return nil
}

func (m *repository) GetAllDecks(userDeck *[]models.Deck) error {
	for id, deck := range m.decks {
		var newDeck models.Deck
		m.GetAllFlashcardByDeckId(&newDeck.Cards, id)
		newDeck.CopyReferences(deck)

		*userDeck = append(*userDeck, newDeck)
	}
	return nil
}
