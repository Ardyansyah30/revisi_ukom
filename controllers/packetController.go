package controllers

import (
	"crud-ukom/config"
	"crud-ukom/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get all packets
func GetPackets(c *gin.Context) {
	var packets []models.Packet
	config.DB.Find(&packets)
	c.JSON(http.StatusOK, packets)
}

// Get packet by ID
func GetPacketByID(c *gin.Context) {
	var packet models.Packet
	id := c.Param("id")
	if err := config.DB.First(&packet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Packet not found!"})
		return
	}
	c.JSON(http.StatusOK, packet)
}

// Create new packet
func CreatePacket(c *gin.Context) {
	var input models.Packet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&input)
	c.JSON(http.StatusOK, input)
}

// Update packet by ID
func UpdatePacket(c *gin.Context) {
	var packet models.Packet
	id := c.Param("id")
	if err := config.DB.First(&packet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Packet not found!"})
		return
	}

	var input models.Packet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&packet).Updates(input)
	c.JSON(http.StatusOK, packet)
}

// Delete packet by ID
func DeletePacket(c *gin.Context) {
	var packet models.Packet
	id := c.Param("id")
	if err := config.DB.First(&packet, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Packet not found!"})
		return
	}
	config.DB.Delete(&packet)
	c.JSON(http.StatusOK, gin.H{"message": "Packet deleted successfully!"})
}

func GetPacketsByUser(c *gin.Context) {
	var orders []models.Order
	userID := c.Param("id") // Mendapatkan userID dari parameter URL

	// Menemukan semua order yang memiliki id_user = userID dan status = 1
	// Melakukan join ke tabel packet berdasarkan id_packet
	if err := config.DB.
		Where("id_user = ? AND status = ?", userID, 1).
		Joins("JOIN packets ON packets.id = orders.id_packet").
		Find(&orders).
		Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada paket ditemukan untuk pengguna ini!"})
		return
	}

	// Menyiapkan slice untuk menyimpan data packet yang terkait
	var packets []models.Packet

	// Ambil data packet terkait dari setiap order
	for _, order := range orders {
		var packet models.Packet
		// Mendapatkan packet berdasarkan id_packet dari order
		if err := config.DB.First(&packet, order.IDPacket).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Paket tidak ditemukan!"})
			return
		}
		// Menambahkan packet ke slice packets
		packets = append(packets, packet)
	}

	// Mengembalikan data paket
	c.JSON(http.StatusOK, packets)
}
