package repository

import (
	"go-quantus-service/src/entities"
	"gorm.io/gorm"
)

type RepositoryUser interface {
	CreateUser(tx *gorm.DB, user *entities.User) error
	UpdateUser(tx *gorm.DB, user *entities.User) error
	DeleteUser(tx *gorm.DB, idUser int64) error
	GetUserByEmail(tx *gorm.DB, email string) (*entities.User, error)
	GetUserByID(tx *gorm.DB, id int64) (*entities.User, error)
	ListUser(tx *gorm.DB) ([]entities.User, error)
}

type UserRepositoryImpl struct {
}

func NewUSerRepository() RepositoryUser {
	return &UserRepositoryImpl{}
}
