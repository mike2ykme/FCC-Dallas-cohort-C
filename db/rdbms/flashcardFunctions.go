package rdbms

import (
	"errors"
	"gorm.io/gorm/clause"
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
	if fc.DeckId <= 0 {
		return 0, errors.New("DeckId cannot be <= 0")
	}

	err := r.DB.Save(&fc).Error
	return fc.ID, err
}
func (r *repository) GetFlashcardById(fc *models.FlashCard, id uint) error {
	if id <= 0 {
		return errors.New("id cannot be <= 0")
	}
	return r.DB.Preload(clause.Associations).First(fc, id).Error
}
func (r *repository) GetAllFlashcardByDeckId(fcs *[]models.FlashCard, id uint) error {
	if id <= 0 {
		return errors.New("deck ID cannot be <= 0")
	}
	return r.DB.Where("deck_id = ?", id).Preload(clause.Associations).Find(fcs).Error
}
func (r *repository) GetAllFlashcards(fcs *[]models.FlashCard) error {
	return r.DB.Preload(clause.Associations).Find(fcs).Error
}

func (r *repository) DeleteFlashcardById(id uint) error {
	return r.DB.Where("id = ?", id).Delete(&models.FlashCard{}).Error
}
