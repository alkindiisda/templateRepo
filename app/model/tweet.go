package model

import "time"

type Tweet struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserID    int       `json:"user_id"`
	Tweet     string    `json:"tweet" gorm:"not null"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserTweet struct {
	Tweet string `json:"tweet" gorm:"not null"`
	Image string `json:"image"`
}

type TweetContent struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Fullname  string    `json:"fullname"`
	Username  string    `json:"username"`
	Tweet     string    `json:"tweet" gorm:"not null"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
