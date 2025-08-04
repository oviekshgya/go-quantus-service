package entities

import "time"

const CONTENT = "contents"

type Content struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Title     string    `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Body      string    `gorm:"column:body;type:text;not null" json:"body"`
	UserID    int64     `gorm:"column:user_id;not null;index" json:"userId"`                             // Foreign key to users
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"` // Optional eager loading
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (Content) TableName() string {
	return CONTENT
}
