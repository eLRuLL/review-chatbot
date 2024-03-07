package routes

import (
	"backend/handlers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.CORSMiddleware())

	r.POST("/conversation", handlers.CreateConversation)
	r.GET("/conversation/:id", handlers.GetConversation)
	r.POST("/conversation/:id/message", handlers.CreateMessage)

	return r
}
