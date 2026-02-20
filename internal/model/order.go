package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID    uint        `gorm:"not null" json:"user_id"`
	ProductID uint        `gorm:"not null" json:"product_id"`
	Quantity  int         `gorm:"not null" json:"quantity"`
	Status    OrderStatus `gorm:"not null;type:varchar(20);default:'pending'" json:"status"`
	User      User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product   Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusFailed     OrderStatus = "failed"
)
