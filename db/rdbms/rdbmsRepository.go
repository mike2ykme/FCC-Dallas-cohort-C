package rdbms

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"teamC/models"
)

type repository struct {
	DB *gorm.DB
}

func NewRdbmsRepository(dbURL string) (*repository, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{PrepareStmt: true})
	if err != nil {
		return nil, err
	}
	return &repository{
		DB: db,
	}, nil
}

func (r *repository) AutoMigrate() error {
	if err := r.DB.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	if err := r.DB.AutoMigrate(&models.Deck{}); err != nil {
		return err
	}
	if err := r.DB.AutoMigrate(&models.FlashCard{}); err != nil {
		return err
	}
	if err := r.DB.AutoMigrate(&models.Answer{}); err != nil {
		return err
	}

	return nil
}
