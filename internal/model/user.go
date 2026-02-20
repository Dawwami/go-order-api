package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string  `gorm:"uniqueIndex;not null" json:"email"`
	Password string  `gorm:"not null" json:"-" binding:"required"`
	Orders   []Order `gorm:"foreignKey:UserID"`
}
