package model

import "time"

type Category struct {
	ID        int64      `gorm:"id"`
	UserId    int64      `gorm:"user_id"`
	Title     string     `gorm:"name"`
	Slug      string     `gorm:"slug"`
	User      User       `gorm:"foreignKey:UserId"`
	CreatedAt time.Time  `gorm:"created_at"`
	UpdatedAt *time.Time `gorm:"updated_at"`
}
