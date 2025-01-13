package models

import (
	"strings"
	"time"
)

// Question struct
type Question struct {
	ID            int64     `json:"id"`
	PacketID      int64     `json:"packet_id"`
	Question      string    `json:"question"`
	Answer        string    `json:"answer"` // Answer is a single string with newlines
	CorrectAnswer string    `json:"correct_answer"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Custom output untuk JSON
type QuestionResponse struct {
	ID            int64    `json:"id"`
	PacketID      int64    `json:"packet_id"`
	Question      string   `json:"question"`
	Answer        []string `json:"answer"` // Array untuk format baris
	CorrectAnswer string   `json:"correct_answer"`
}

// Konversi Question ke QuestionResponse
func (q *Question) ToResponse() QuestionResponse {
	return QuestionResponse{
		ID:            q.ID,
		PacketID:      q.PacketID,
		Question:      q.Question,
		Answer:        strings.Split(q.Answer, "\n"), // Split jawaban berdasarkan newline
		CorrectAnswer: q.CorrectAnswer,
	}
}
