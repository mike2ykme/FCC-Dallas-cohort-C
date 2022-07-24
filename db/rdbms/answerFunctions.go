package rdbms

import (
	"errors"
	"teamC/models"
)

/*

type AnswerRepository interface {
	SaveAnswer(answer *models.Answer) (uint, error)
	GetAnswerById(answer *models.Answer, id uint) error
	GetAnswersByFlashcardId(answers *[]models.Answer) error
	GetAllAnswers(answers *[]models.Answer) error
}

*/

func (r *repository) SaveAnswer(answer *models.Answer) (uint, error) {
	if answer.FlashCardId == 0 {
		return 0, errors.New("cannot have a 0 flashcard ID")
	}

	err := r.DB.Save(answer).Error

	return answer.ID, err
}

func (r *repository) GetAnswerById(answer *models.Answer, id uint) error {
	if id == 0 {
		return errors.New("id cannot be 0")
	}
	return r.DB.First(&answer, id).Error
}

func (r *repository) GetAnswersByFlashcardId(answers *[]models.Answer, id uint) error {
	if id == 0 {
		return errors.New("id cannot be 0")
	}
	return r.DB.Where("flash_card_id = ?", id).Find(answers).Error
}

func (r *repository) GetAllAnswers(answers *[]models.Answer) error {
	return r.DB.Find(answers).Error
}

func (r *repository) DeleteAnswerById(id uint) error {
	return r.DB.Where("id = ?", id).Delete(&models.Answer{}).Error
}
