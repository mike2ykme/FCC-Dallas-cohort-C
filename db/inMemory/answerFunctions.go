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

func (r *repository) SaveAnswer(answer *models.Answer) (uint, error) {
	if answer.FlashCardId == 0 {
		return 0, errors.New("cannot have a 0 flashcard ID")
	}
	if answer.ID == 0 {
		answer.ID = r.currentHighestAnswerId
		r.currentHighestAnswerId++
	} else if answer.ID > r.currentHighestAnswerId {
		r.currentHighestAnswerId = answer.ID + 1
	}
	var copy models.Answer
	copy.CopyRef(answer)

	r.answers[answer.ID] = &copy

	return answer.ID, nil
}

func (r *repository) GetAnswerById(answer *models.Answer, id uint) error {
	if id == 0 {
		return errors.New("id cannot be 0")
	}

	if val, ok := r.answers[id]; ok {
		answer.CopyRef(val)
	}
	return nil
}

func (r *repository) GetAnswersByFlashcardId(answers *[]models.Answer, id uint) error {
	if id == 0 {
		return errors.New("id cannot be 0")
	}
	for _, answer := range r.answers {
		if answer.FlashCardId == id {
			*answers = append(*answers, answer.Copy())
		}
	}

	return nil
}

func (r *repository) GetAllAnswers(answers *[]models.Answer) error {
	if *answers == nil || len(*answers) == 0 {
		*answers = make([]models.Answer, len(r.answers))

		idx := 0
		for _, ans := range r.answers {
			copy := ans.Copy()
			copy.CopyRef(ans)
			(*answers)[idx] = copy
			idx++
		}
		return nil
	}

	for _, answer := range r.answers {
		*answers = append(*answers, answer.Copy())
	}

	return nil
}

func (r *repository) DeleteAnswerById(id uint) error {
	delete(r.answers, id)
	return nil
}
