package rdbms

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
