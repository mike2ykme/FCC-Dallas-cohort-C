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
	if fc.Id == 0 {
		fc.Id = m.currentHighestCardId
		m.currentHighestCardId++
	} else if fc.Id > m.currentHighestCardId {
		m.currentHighestCardId = fc.Id + 1
	}

	var copy models.FlashCard
	copy.CopyRef(fc)

	m.flashcards[fc.Id] = &copy
	for _, answer := range fc.Answers {
		answer.FlashCardId = fc.Id
		m.SaveAnswer(&answer)
	}

	return fc.Id, nil
}
func (m *repository) GetFlashcardById(fc *models.FlashCard, id uint) error {
	if val, ok := m.flashcards[id]; ok {
		fc.CopyRef(val)
		//if val.Id == id {
		//	m.GetAnswersByFlashcardId(&val.Answers, fc.Id)
		//	fc.CopyRef(val)
		//}
	}
	return nil
}
func (m *repository) GetAllFlashcardByDeckId(fcs *[]models.FlashCard, id uint) error {
	if id == 0 {
		return errors.New("deck Id cannot be 0")
	}

	for _, card := range m.flashcards {
		if card.DeckId == id {
			m.GetAnswersByFlashcardId(&card.Answers, card.Id)
			*fcs = append(*fcs, card.Copy())
		}
	}
	return nil
}
func (m *repository) GetAllFlashcards(fcs *[]models.FlashCard) error {
	//if len(*fcs) == 0 {
	//	*fcs = make([]models.FlashCard, len(m.flashcards))
	//	for idx, card := range m.flashcards {
	//		copy := card.Copy()
	//		copy.CopyRef(card)
	//		(*fcs)[idx] = copy
	//	}
	//	return nil
	//}
	for _, card := range m.flashcards {
		copy := card.Copy()
		copy.CopyRef(card)
		m.GetAnswersByFlashcardId(&copy.Answers, copy.Id)
		*fcs = append(*fcs, copy)
	}
	return nil
}
