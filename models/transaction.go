package models

import "time"

type Transaction struct {
	Id         int       `json:"id" gorm:"primaryKey"`
	UserId     int       `json:"user_id"`
	User       User      `json:"user" gorm:"foreignKey:UserId"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
	CategoryId int       `json:"category_id"`
	Category   Category  `json:"category" gorm:"foreignKey:CategoryId"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
