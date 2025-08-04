package pkg

import "time"

type RawUser struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	FullName  string    `gorm:"column:full_name" json:"fullName"`
	Email     string    `gorm:"column:email;uniqueIndex" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	Role      string    `gorm:"column:role" json:"role"`
	IsActive  bool      `gorm:"column:is_active" json:"isActive"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
