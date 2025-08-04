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

type RepositoryContent interface {
	CreateContent(tx *gorm.DB, content *entities.Content) error
	UpdateContent(tx *gorm.DB, content *entities.Content) error
	DeleteContent(tx *gorm.DB, contentID int64) error
	GetContentByID(tx *gorm.DB, contentID int64) (*entities.Content, error)
	ListContentByUserID(tx *gorm.DB, userID int64) ([]entities.Content, error)
}

type ContentRepositoryImpl struct{}

func NewContentRepository() RepositoryContent {
	return &ContentRepositoryImpl{}
}
