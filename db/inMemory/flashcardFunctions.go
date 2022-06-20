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

func (r *repository) SaveFlashcard(fc *models.FlashCard) (uint, error) {
	if fc.DeckId == 0 {
		return 0, errors.New("DeckId cannot be 0")
	}
	if fc.ID == 0 {
		fc.ID = r.currentHighestCardId
		r.currentHighestCardId++
	} else if fc.ID > r.currentHighestCardId {
		r.currentHighestCardId = fc.ID + 1
	}

	var copy models.FlashCard
	copy.CopyRef(fc)
	copy.Answers = nil
	r.flashcards[fc.ID] = &copy

	for i := 0; i < len(fc.Answers); i++ {
		fc.Answers[i].FlashCardId = fc.ID
		if _, err := r.SaveAnswer(&fc.Answers[i]); err != nil {
			return 0, err
		}
	}

	return fc.ID, nil
}
func (r *repository) GetFlashcardById(fc *models.FlashCard, id uint) error {
	if val, ok := r.flashcards[id]; ok {
		fc.CopyRef(val)
		return nil
	}
	return errors.New("there was no flashcard found")
}
func (r *repository) GetAllFlashcardByDeckId(fcs *[]models.FlashCard, id uint) error {
	if id == 0 {
		return errors.New("deck ID cannot be 0")
	}

	for _, card := range r.flashcards {
		if card.DeckId == id {
			newCard := card.Copy()
			r.GetAnswersByFlashcardId(&newCard.Answers, card.ID)
			*fcs = append(*fcs, newCard)
		}
	}
	return nil
}
func (r *repository) GetAllFlashcards(fcs *[]models.FlashCard) error {
	if *fcs == nil || len(*fcs) == 0 {
		*fcs = make([]models.FlashCard, len(r.flashcards))
		idx := 0
		for _, card := range r.flashcards {
			copy := card.Copy()
			copy.CopyRef(card)
			(*fcs)[idx] = copy
			idx++
		}
		return nil
	}
	for _, card := range r.flashcards {
		copy := card.Copy()
		copy.CopyRef(card)
		r.GetAnswersByFlashcardId(&copy.Answers, copy.ID)
		*fcs = append(*fcs, copy)
	}
	return nil
}

//func (m *repository) DeleteDeckById(id uint) error {
//	delete(m.decks, id)
//	return nil
//}

func (r *repository) DeleteFlashcardById(id uint) error {
	var currentCard models.FlashCard

	err := r.GetFlashcardById(&currentCard, id)
	if err != nil {
		return err
	}

	for _, answer := range currentCard.Answers {
		delErr := r.DeleteAnswerById(answer.ID)
		if delErr != nil {
			return delErr
		}
	}

	delete(r.flashcards, id)

	return nil
}
