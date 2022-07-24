package inMemory

import (
	"errors"
	"teamC/models"
)

func (r *repository) SaveDeck(deck *models.Deck) (uint, error) {
	if deck.OwnerId == 0 {
		return 0, errors.New("deck must have an owwner")
	}
	if deck.ID == 0 {
		deck.ID = r.currentHighestDeckId
		r.currentHighestDeckId++
	} else if deck.ID > r.currentHighestDeckId {
		r.currentHighestDeckId = deck.ID + 1
	}
	copy := deck.Copy()
	copy.FlashCards = nil

	for i := 0; i < len(deck.FlashCards); i++ {
		deck.FlashCards[i].DeckId = deck.ID
		if _, err := r.SaveFlashcard(&deck.FlashCards[i]); err != nil {
			return 0, err
		}
	}
	r.decks[deck.ID] = &copy

	return deck.ID, nil
}

func (r *repository) GetDeckById(d *models.Deck, id uint) error {
	if val, ok := r.decks[id]; ok {
		if val.ID == id {
            newDeck := val.Copy()
			r.GetAllFlashcardByDeckId(&newDeck.FlashCards, newDeck.ID)
			d.CopyReferences(&newDeck)
		}
	} else {
		return errors.New("deck does not exist")
	}
	return nil
}

func (r *repository) GetAllDecks(userDeck *[]models.Deck) error {

	if *userDeck == nil || len(*userDeck) == 0 {
		*userDeck = make([]models.Deck, len(r.decks))
		idx := 0
		for _, deck := range r.decks {
			copy := deck.Copy()
			copy.CopyReferences(deck)
			(*userDeck)[idx] = copy
			idx++
		}
		return nil
	}
	for id, deck := range r.decks {
		var newDeck models.Deck
		r.GetAllFlashcardByDeckId(&newDeck.FlashCards, id)
		newDeck.CopyReferences(deck)

		*userDeck = append(*userDeck, newDeck)
	}
	return nil
}

func (r *repository) GetDecksByUserId(decks *[]models.Deck, userId uint) error {
	for _, deck := range r.decks {
		if deck.OwnerId == userId {
			newDeck := deck.Copy()
			r.GetAllFlashcardByDeckId(&newDeck.FlashCards, newDeck.ID)
			*decks = append(*decks, newDeck)
		}
	}
	return nil

}

func (r *repository) DeleteDeckById(id uint) error {
	delete(r.decks, id)
	return nil
}
