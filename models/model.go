package models

import "time"

type User struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Post struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	AuthorId  uint64    `json:"user_id"`
	TimeStamp time.Time `json:"timestamp"`
}
