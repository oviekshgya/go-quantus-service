package repository

import (
	"go-quantus-service/src/entities"
	"go-quantus-service/src/pkg"
	"gorm.io/gorm"
	"time"
)

func (r *ContentRepositoryImpl) CreateContent(tx *gorm.DB, content *entities.Content) error {
	content.CreatedAt = time.Now()
	content.UpdatedAt = time.Now()
	return tx.Create(content).Error
}

func (r *ContentRepositoryImpl) UpdateContent(tx *gorm.DB, content *entities.Content) error {
	content.UpdatedAt = time.Now()

	var existing entities.Content
	if err := tx.First(&existing, content.ID).Error; err != nil {
		return err
	}

	data := pkg.UpdateFieldsDynamic(content)
	data["updated_at"] = content.UpdatedAt

	return tx.Model(&entities.Content{}).
		Where("id = ? AND user_id = ?", content.ID, content.UserID).
		Updates(data).Error
}

func (r *ContentRepositoryImpl) DeleteContent(tx *gorm.DB, contentID int64) error {
	return tx.Delete(&entities.Content{}, contentID).Error
}

func (r *ContentRepositoryImpl) GetContentByID(tx *gorm.DB, contentID int64) (*entities.Content, error) {
	var content entities.Content
	if err := tx.First(&content, contentID).Error; err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *ContentRepositoryImpl) ListContentByUserID(tx *gorm.DB, userID int64) ([]entities.Content, error) {
	var contents []entities.Content
	err := tx.Where("user_id = ?", userID).Order("created_at desc").Find(&contents).Error
	return contents, err
}
