package routes

import (
	"crud-ukom/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	// User routes
	r.GET("/users", controllers.GetUsers)
	r.GET("/users/:id", controllers.GetUserByID)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	// Question routes
	r.GET("/questions", controllers.GetQuestions)
	r.GET("/questions/:id", controllers.GetQuestionByID)
	r.POST("/questions", controllers.CreateQuestion)
	r.PUT("/questions/:id", controllers.UpdateQuestion)
	r.DELETE("/questions/:id", controllers.DeleteQuestion)
	// Packet routes
	r.GET("/packets", controllers.GetPackets)
	r.GET("/packets-detail/:id", controllers.GetPacketByID)
	r.GET("/packets/:packet_id/questions", controllers.GetQuestionsByPacketID)
	r.GET("/packets-purchased/:id", controllers.GetPacketsByUser)
	r.POST("/packets", controllers.CreatePacket)
	r.PUT("/packets/:id", controllers.UpdatePacket)
	r.DELETE("/packets/:id", controllers.DeletePacket)

	return r
}
