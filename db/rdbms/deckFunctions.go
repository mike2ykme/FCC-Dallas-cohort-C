package rdbms

import (
	"errors"
	"gorm.io/gorm/clause"
	"teamC/models"
)

func (r *repository) SaveDeck(deck *models.Deck) (uint, error) {
	if deck.OwnerId <= 0 {
		return 0, errors.New("deck must have an owner")
	}

	err := r.DB.Save(deck).Error
	return deck.ID, err
}

func (r *repository) GetDeckById(d *models.Deck, id uint) error {
	if id <= 0 {
		return errors.New("id cannot be 0")
	}
	return r.DB.Where("id = ?", id).
		Preload("FlashCards.Answers").
		Preload(clause.Associations).
		First(&d).
		Error
}

func (r *repository) GetAllDecks(decks *[]models.Deck) error {
	return r.DB.Preload("FlashCards.Answers").
		Preload(clause.Associations).
		Find(decks).
		Error
}

func (r *repository) GetDecksByUserId(decks *[]models.Deck, userId uint) error {
	if userId <= 0 {
		return errors.New("user ID cannot be 0")
	}
	return r.DB.Where("owner_id = ?", userId).
		Preload("FlashCards.Answers").
		Preload(clause.Associations).
		Find(decks).Error

}

func (r *repository) DeleteDeckById(id uint) error {
	return r.DB.Where("id = ?", id).Delete(&models.Deck{}).Error
}
