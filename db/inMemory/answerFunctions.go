package inMemory

import (
	"errors"
	"teamC/models"
)

/*

type AnswerRepository interface {
	SaveAnswer(answer *models.Answer) (uint, error)
	GetAnswerById(answer *models.Answer) error
	GetAnswersByFlashcardId(answers *[]models.Answer) error
	GetAllAnswers(answers *[]models.Answer) error
}

*/

func (m *repository) SaveAnswer(answer *models.Answer) (uint, error) {
	if answer.FlashCardId == 0 {
		return 0, errors.New("cannot have a 0 flashcard ID")
	}
	if answer.Id == 0 {
		answer.Id = m.currentHighestAnswerId
		m.currentHighestAnswerId++
	} else if answer.Id > m.currentHighestAnswerId {
		m.currentHighestAnswerId = answer.Id + 1
	}
	var copy models.Answer
	copy.CopyRef(answer)

	m.answers[answer.Id] = &copy

	return answer.Id, nil
}

func (m *repository) GetAnswerById(answer *models.Answer, id uint) error {
	if id == 0 {
		return errors.New("id cannot be 0")
	}

	if val, ok := m.answers[id]; ok {
		answer.CopyRef(val)
	}
	return nil
}

func (m *repository) GetAnswersByFlashcardId(answers *[]models.Answer, id uint) error {
	if id == 0 {
		return errors.New("id cannot be 0")
	}
	for _, answer := range m.answers {
		if answer.FlashCardId == id {
			*answers = append(*answers, answer.Copy())
		}
	}

	return nil
}

func (m *repository) GetAllAnswers(answers *[]models.Answer) error {
	for _, answer := range m.answers {
		*answers = append(*answers, answer.Copy())
	}

	return nil
}
