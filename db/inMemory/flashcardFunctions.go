package inMemory

import (
	"errors"
	"teamC/models"
)

/*
type FlashcardRepository interface {
	SaveFlashcard(*models.FlashCard) (uint, error)
	GetFlashcardById(*models.FlashCard, uint) error
	GetFlashcardsByDeck(card *models.FlashCard) error
	GetAllFlashcards(*models.FlashCard) error
}

*/

func (m *repository) SaveFlashcard(fc *models.FlashCard) (uint, error) {
	if fc.DeckId == 0 {
		return 0, errors.New("DeckId cannot be 0")
	}
	if fc.Id == 0 {
		fc.Id = m.currentHighestCardId
		m.currentHighestCardId++
	} else if fc.Id > m.currentHighestCardId {
		m.currentHighestCardId = fc.Id + 1
	}
	m.flashcards[fc.Id] = fc
	for _, answer := range fc.Answers {
		m.SaveAnswer(&answer)
	}

	return fc.Id, nil
}
func (m *repository) GetFlashcardById(fc *models.FlashCard, id uint) error {
	if val, ok := m.flashcards[fc.Id]; ok {
		if val.Id == id {
			fc.CopyRef(val)
		}
	}
	return nil
}
func (m *repository) GetFlashcardsByDeck(fcs *[]models.FlashCard, id uint) error {
	if id == 0 {
		return errors.New("deck Id cannot be 0")
	}

	for _, card := range m.flashcards {
		if card.DeckId == id {
			*fcs = append(*fcs, card.Copy())
		}
	}
	return nil
}
func (m *repository) GetAllFlashcards(fcs *[]models.FlashCard) error {
	if len(*fcs) == 0 {
		*fcs = make([]models.FlashCard, len(m.flashcards))
		for idx, card := range m.flashcards {
			(*fcs)[idx] = card.Copy()
		}
		return nil
	}
	for _, card := range m.flashcards {
		*fcs = append(*fcs, card.Copy())
	}
	return nil
}
