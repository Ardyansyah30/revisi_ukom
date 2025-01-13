package models

import "time"

type User struct {
	ID                     uint      `gorm:"primaryKey"`
	Name                   string    `json:"name"`
	Email                  string    `json:"email"`
	Password               string    `json:"password"`
	PhoneNumber            string    `json:"phone_number"`
	DateOfBirth            string    `json:"date_of_birth"`
	Gender                 string    `json:"gender"`
	EducationalInstitution string    `json:"educational_institution"`
	Profession             string    `json:"profession"`
	Address                string    `json:"address"`
	Province               string    `json:"province"`
	City                   string    `json:"city"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
