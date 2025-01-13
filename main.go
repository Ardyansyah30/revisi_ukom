package main

import (
	"crud-ukom/config"
	"crud-ukom/models"
	"crud-ukom/routes"
)

func main() {
	config.ConnectDB()

	config.DB.AutoMigrate(&models.User{}, &models.Question{}, &models.Packet{})

	router := routes.SetupRoutes()
	router.Run("0.0.0.0:8080")
}
