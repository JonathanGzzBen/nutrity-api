package models

import "time"

type Article struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"userId"`
	User       User      `json:"user"`
	CategoryID uint      `json:"categoryId"`
	Category   Category  `json:"category"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updtedAt"`
	Body       string    `json:"body"`
	Title      string    `json:"title"`
	ImageURL   string    `json:"imageUrl"`
	// Tags is a comma separated string of tags
	Tags string `json:"tags"`
}
