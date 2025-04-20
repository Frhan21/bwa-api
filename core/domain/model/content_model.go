package model

import "time"

type Content struct {
	ID          int64      `gorm:"id"`
	CategoryID  int64      `gorm:"category_id"`
	UserID      int64      `gorm:"user_id"`
	Title       string     `gorm:"title"`
	Excerpt     string     `gorm:"excerpt"`
	Tags        string     `gorm:"tags"`
	Status      string     `grom:"status"`
	Description string     `gorm:"description"`
	Image       string     `gorm:"image"`
	CreatedAt   time.Time  `gorm:"created_at"`
	UpdatedAt   *time.Time `gorm:"updated_at"`
	Category    Category   `gorm:"foreignKey:CategoryID"`
	User        User       `gorm:"foreignKey:UserID"`
}
