package inMemory

import "teamC/models"

func (m *repository) SaveDeck(deck *models.Deck) (uint, error) {
	if deck.Id == 0 {
		deck.Id = m.currentHighestDeckId
		m.currentHighestDeckId++
	} else if deck.Id > m.currentHighestDeckId {
		m.currentHighestDeckId = deck.Id + 1
	}
	m.decks[deck.Id] = deck

	return deck.Id, nil
}

func (m *repository) GetDeckById(d *models.Deck, id uint) error {
	if val, ok := m.decks[d.Id]; ok {
		if val.Id == id {
			d.CopyReferences(val)
		}
	}
	return nil
}

func (m *repository) GetAllDecks(userDeck *[]models.Deck) error {
	for _, deck := range m.decks {
		var newDeck models.Deck
		newDeck.CopyReferences(deck)
		*userDeck = append(*userDeck, newDeck)
	}
	return nil
}
