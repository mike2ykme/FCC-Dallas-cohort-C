package inMemory

import (
	"errors"
	"teamC/models"
)

/*
type FlashcardRepository interface {
	SaveFlashcard(*models.FlashCard) (uint, error)
	GetFlashcardById(*models.FlashCard, uint) error
	GetAllFlashcardByDeckId(card *models.FlashCard) error
	GetAllFlashcards(*models.FlashCard) error
}

*/

func (m *repository) SaveFlashcard(fc *models.FlashCard) (uint, error) {
	if fc.DeckId == 0 {
		return 0, errors.New("DeckId cannot be 0")
	}
	if fc.ID == 0 {
		fc.ID = m.currentHighestCardId
		m.currentHighestCardId++
	} else if fc.ID > m.currentHighestCardId {
		m.currentHighestCardId = fc.ID + 1
	}

	var copy models.FlashCard
	copy.CopyRef(fc)
	copy.Answers = nil
	m.flashcards[fc.ID] = &copy

	for i := 0; i < len(fc.Answers); i++ {
		fc.Answers[i].FlashCardId = fc.ID
		if _, err := m.SaveAnswer(&fc.Answers[i]); err != nil {
			return 0, err
		}
	}

	return fc.ID, nil
}
func (m *repository) GetFlashcardById(fc *models.FlashCard, id uint) error {
	if val, ok := m.flashcards[id]; ok {
		fc.CopyRef(val)
		//if val.ID == id {
		//	m.GetAnswersByFlashcardId(&val.Answers, fc.ID)
		//	fc.CopyRef(val)
		//}
	}
	return nil
}
func (m *repository) GetAllFlashcardByDeckId(fcs *[]models.FlashCard, id uint) error {
	if id == 0 {
		return errors.New("deck ID cannot be 0")
	}

	for _, card := range m.flashcards {
		if card.DeckId == id {
			m.GetAnswersByFlashcardId(&card.Answers, card.ID)
			*fcs = append(*fcs, card.Copy())
		}
	}
	return nil
}
func (m *repository) GetAllFlashcards(fcs *[]models.FlashCard) error {
	if *fcs == nil || len(*fcs) == 0 {
		*fcs = make([]models.FlashCard, len(m.flashcards))
		idx := 0
		for _, card := range m.flashcards {
			copy := card.Copy()
			copy.CopyRef(card)
			(*fcs)[idx] = copy
			idx++
		}
		return nil
	}
	for _, card := range m.flashcards {
		copy := card.Copy()
		copy.CopyRef(card)
		m.GetAnswersByFlashcardId(&copy.Answers, copy.ID)
		*fcs = append(*fcs, copy)
	}
	return nil
}
