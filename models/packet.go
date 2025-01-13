package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Packet struct {
	ID           uint      `gorm:"primaryKey"`
	NamePacket   string    `gorm:"type:varchar(255)"`
	Description  string    `gorm:"type:text"`
	Price        float64   `gorm:"type:decimal(10,2)"`
	DurationExam string    `gorm:"type:varchar(100)"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (p *Packet) UnmarshalJSON(data []byte) error {

	var tmp struct {
		ID           uint        `json:"id"`
		NamePacket   string      `json:"name_packet"`
		Description  string      `json:"description"`
		Price        interface{} `json:"price"`
		DurationExam string       `json:"duration_exam"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// Assign the fields
	p.ID = tmp.ID
	p.NamePacket = tmp.NamePacket
	p.Description = tmp.Description
	p.DurationExam = tmp.DurationExam
	p.CreatedAt = tmp.CreatedAt
	p.UpdatedAt = tmp.UpdatedAt

	// Handle Price, which can be either a string or a number
	switch value := tmp.Price.(type) {
	case string:
		parsedPrice, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("could not parse price string: %v", err)
		}
		p.Price = parsedPrice
	case float64:
		p.Price = value
	default:
		return fmt.Errorf("invalid type for price field")
	}

	return nil
}
