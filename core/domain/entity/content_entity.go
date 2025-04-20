package entity

import "time"

type ContentEntity struct {
	ID          int64
	Title       string
	Excerpt     string
	Description string
	Image       string
	Category    CategoryEntity
	User        UserEntity
	Tags        []string
	Status      string
	CategoryId  int64
	UserID      int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type QueryString struct {
	Limit      int
	Page       int
	Order      string
	OrderBy    string
	OrderType  string
	Search     string
	CategoryId int64
	Status     string
}
