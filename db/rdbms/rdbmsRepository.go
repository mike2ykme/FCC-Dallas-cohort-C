package rdbms

import (
	"strings"
	"teamC/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type repository struct {
	DB *gorm.DB
}

func NewRdbmsRepository(dbURL string, DBMSName string) (*repository, error) {
	var connectionFunction func(dsn string) (gorm.Dialector)
	if strings.ToLower(DBMSName) == "postgres" {
		connectionFunction = postgres.Open
	} else if strings.ToLower(DBMSName) == "sqlite" {
		connectionFunction = sqlite.Open
	} else {
		panic("Unrecognized DBMS")
	}
	db, err := gorm.Open(connectionFunction(dbURL), &gorm.Config{PrepareStmt: true})
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
