package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name   string  `gorm:"not null" json:"name"`
	Price  float64 `gorm:"not null" json:"price"`
	Stock  int     `gorm:"not null;default:0" json:"stock"`
	Orders []Order `gorm:"foreignKey:ProductID"`
}
