package models

import "time"

type Transaction struct {
	Id	       	int      	`json:"id"`
	UserId    	int      	`json:"user_id"`
	Amount    	float64   `json:"amount"`
	Type      	string    `json:"type"`
	CategoryId	int      	`json:"category_id"`
	Category		string    `json:"category"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
}