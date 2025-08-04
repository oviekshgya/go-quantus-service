package repository

import (
	"go-quantus-service/src/entities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

func (u *UserRepositoryImpl) CreateUser(tx *gorm.DB, user *entities.User) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPwd)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return tx.Create(user).Error
}

func (u *UserRepositoryImpl) UpdateUser(tx *gorm.DB, user *entities.User) error {
	user.UpdatedAt = time.Now()
	return tx.Model(&entities.User{}).Where("idUser = ?", user.ID).Updates(user).Error
}

func (u *UserRepositoryImpl) DeleteUser(tx *gorm.DB, idUser int64) error {
	return tx.Delete(&entities.User{}, idUser).Error
}

func (u *UserRepositoryImpl) GetUserByEmail(tx *gorm.DB, email string) (*entities.User, error) {
	var user entities.User
	if err := tx.Where("email = ?", email).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepositoryImpl) ListUser(tx *gorm.DB) ([]entities.User, error) {
	var User []entities.User
	err := tx.Find(&User).Error
	return User, err
}
