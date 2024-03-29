package model

import "time"

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Fullname  string    `json:"fullname" gorm:"type:varchar(255);"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null"`
	Username  string    `json:"username" gorm:"type:varchar(255);not null"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegister struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdate struct {
	Fullname string `json:"fullname" binding:"required"`
	Username string `json:"username" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Image    string `json:"image"`
}
