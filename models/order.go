package models

import (
	"time"
)

type Order struct {
	ID        int       `gorm:"primaryKey"`
	IDUser    int       `gorm:"not null"`
	IDPacket  int       `gorm:"not null"`
	Status    int       `gorm:"not null"`
	OrderDate time.Time `gorm:"type:date;not null"` 
	Amount    float64   `gorm:"type:decimal(10,2)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
