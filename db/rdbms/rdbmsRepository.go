package rdbms

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
	"teamC/models"
)

type repository struct {
	DB *gorm.DB
}

func getConnectionFunc(name string) (func(dsn string) gorm.Dialector, error) {
	switch strings.ToLower(name) {
	case "postgres":
		return postgres.Open, nil
	case "sqlite":
		return sqlite.Open, nil
	default:
		return nil, errors.New("unrecognized DBMS")
	}
}

func NewRdbmsRepository(dbURL string, DBMSName string) (*repository, error) {
	connectionFunction, err := getConnectionFunc(DBMSName)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(connectionFunction(dbURL), &gorm.Config{
		PrepareStmt: true, FullSaveAssociations: true,
	})

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
