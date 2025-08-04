package entities

import "time"

const USER = "users"

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	FullName  string    `gorm:"column:full_name" json:"fullName"`
	Email     string    `gorm:"column:email;uniqueIndex" json:"email"`
	Password  string    `gorm:"column:password" json:"-"`
	Role      string    `gorm:"column:role" json:"role"`
	IsActive  bool      `gorm:"column:is_active" json:"isActive"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (User) TableName() string {
	return USER
}
