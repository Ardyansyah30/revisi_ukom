package controllers

import (
	"crud-ukom/config"
	"crud-ukom/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Create a new question
func CreateQuestion(c *gin.Context) {
	var input struct {
		PacketID      int64  `json:"packet_id" binding:"required"`
		Question      string `json:"question" binding:"required"`
		Answer        string `json:"answer"`
		CorrectAnswer string `json:"correct_answer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format answer to have newline between options
	formattedAnswer := strings.ReplaceAll(input.Answer, ", ", "\n")

	question := models.Question{
		PacketID:      input.PacketID,
		Question:      input.Question,
		Answer:        formattedAnswer,
		CorrectAnswer: input.CorrectAnswer,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := config.DB.Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, question)
}

// Get all questions
func GetQuestions(c *gin.Context) {
	var questions []models.Question
	if err := config.DB.Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Periksa jika tidak ada data
	if len(questions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No questions found"})
		return
	}

	// Konversi setiap pertanyaan menggunakan ToResponse
	var responses []struct {
		PacketID      int64    `json:"packet_id"`
		Question      string   `json:"question"`
		Answer        []string `json:"answer"`
		CorrectAnswer string   `json:"correct_answer"`
	}

	// Format dan masukkan ke dalam responses
	for _, question := range questions {
		answers := strings.Split(question.Answer, "\n")
		responses = append(responses, struct {
			PacketID      int64    `json:"packet_id"`
			Question      string   `json:"question"`
			Answer        []string `json:"answer"`
			CorrectAnswer string   `json:"correct_answer"`
		}{
			PacketID:      question.PacketID,
			Question:      question.Question,
			Answer:        answers,
			CorrectAnswer: question.CorrectAnswer,
		})
	}

	// Kirim response dalam bentuk array
	c.JSON(http.StatusOK, responses)
}

// Get questions by PacketID
func GetQuestionsByPacketID(c *gin.Context) {
	// Ambil parameter packet_id
	PacketIDStr := c.Param("packet_id")

	// Periksa apakah kosong
	if PacketIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PacketID is required"})
		return
	}

	// Konversi string ke int64
	PacketID, err := strconv.ParseInt(PacketIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PacketID format"})
		return
	}

	// Ambil pertanyaan berdasarkan PacketID
	var questions []models.Question
	if err := config.DB.Where("packet_id = ?", PacketID).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Periksa apakah ada hasil
	if len(questions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No questions found for the given PacketID"})
		return
	}

	// Format response
	var responses []struct {
		PacketID      int64    `json:"packet_id"`
		Question      string   `json:"question"`
		Answer        []string `json:"answer"`
		CorrectAnswer string   `json:"correct_answer"`
	}

	// Loop untuk menambahkan data pertanyaan dalam format yang diinginkan
	for _, question := range questions {
		answers := strings.Split(question.Answer, "\n")
		responses = append(responses, struct {
			PacketID      int64    `json:"packet_id"`
			Question      string   `json:"question"`
			Answer        []string `json:"answer"`
			CorrectAnswer string   `json:"correct_answer"`
		}{
			PacketID:      question.PacketID,
			Question:      question.Question,
			Answer:        answers,
			CorrectAnswer: question.CorrectAnswer,
		})
	}

	// Kirim response dalam bentuk array
	c.JSON(http.StatusOK, responses)
}

// Get a question by ID
func GetQuestionByID(c *gin.Context) {
	var question models.Question
	if err := config.DB.Where("id = ?", c.Param("id")).First(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, question.ToResponse())
}

// Update a question by ID
func UpdateQuestion(c *gin.Context) {
	var question models.Question
	if err := config.DB.Where("id = ?", c.Param("id")).First(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	var input struct {
		PacketID      int64  `json:"packet_id" binding:"required"`
		Question      string `json:"question" binding:"required"`
		Answer        string `json:"answer"`
		CorrectAnswer string `json:"correct_answer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format answer to have newline between options
	formattedAnswer := strings.ReplaceAll(input.Answer, ", ", "\n")

	// Check if the answer is correct
	question.PacketID = input.PacketID
	question.Question = input.Question
	question.Answer = formattedAnswer
	question.CorrectAnswer = input.CorrectAnswer
	question.UpdatedAt = time.Now()

	config.DB.Save(&question)
	c.JSON(http.StatusOK, question)
}

// Delete a question by ID
func DeleteQuestion(c *gin.Context) {
	var question models.Question
	if err := config.DB.Where("id = ?", c.Param("id")).First(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	config.DB.Delete(&question)
	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}
