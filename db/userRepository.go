package db

import "teamC/models"

type UserRepository interface {
	SaveUser(*models.User) (uint, error)
	GetUserById(*models.User, uint) error
	GetUserByUsername(*models.User, string) error
	GetUserBySubId(*models.User, string) error
	GetAllUsers(*[]models.User) error
}
