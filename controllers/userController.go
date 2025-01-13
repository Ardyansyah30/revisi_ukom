package controllers

import (
	"crud-ukom/config"
	"crud-ukom/models"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_jwt_secret")

func isValidPhoneNumber(phone string) bool { // Fungsi untuk memvalidasi apakah nomor telepon hanya berisi angka
	regex := regexp.MustCompile(`^\+?[0-9]+$`)
	return regex.MatchString(phone)
}
func hashPassword(password string) (string, error) { // Hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func checkPasswordHash(password, hash string) bool { // Fungsi untuk memvalidasi apakah nomor telepon hanya berisi angka
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate JWT token
func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Signup user baru
func Signup(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isValidPhoneNumber(input.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number must contain only numeric characters"})
		return
	}
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user := models.User{
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Password:    hashedPassword,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":           user.ID,
			"name":         user.Name,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"created_at":   user.CreatedAt,
		},
	})
}

// Login user
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	if !checkPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// Return login success response along with user data and token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":           user.ID,
			"name":         user.Name,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"created_at":   user.CreatedAt,
		},
	})
}

// Get all users
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Get a user by ID
func GetUserByID(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Update User
func UpdateUser(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isValidPhoneNumber(user.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number must contain only numeric characters"})
		return
	}

	// Parse date jika DateOfBirth diisi
	if user.DateOfBirth != "" {
		parsedDate, err := time.Parse("2006-01-02", user.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
			return
		}
		user.DateOfBirth = parsedDate.Format("2006-01-02")
	}

	// Update password hanya jika diisi
	if user.Password != "" {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = hashedPassword
	}
	user.UpdatedAt = time.Now()
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete a user by ID
func DeleteUser(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
