package db

import "teamC/models"

type AnswerRepository interface {
	SaveAnswer(answer *models.Answer) (uint, error)
	GetAnswerById(answer *models.Answer, id uint) error
	GetAnswersByFlashcardId(answers *[]models.Answer, id uint) error
	//GetAnswersByDeckId(answers *[]models.Answer, id uint) error
	GetAllAnswers(answers *[]models.Answer) error
	DeleteAnswerById(uint) error
}
