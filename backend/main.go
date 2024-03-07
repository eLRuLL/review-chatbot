package main

import (
	"backend/db"
	"backend/routes"
)

func main() {
	db.ConnectDatabase()

	r := routes.SetupRouter()

	// r.POST("/", func(c *gin.Context) {
	// 	var reqBody struct {
	// 		Message string `json:"message"`
	// 	}
	// 	if err := c.ShouldBindJSON(&reqBody); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	replyMessage := "Received your message: " + reqBody.Message
	// 	c.JSON(http.StatusOK, gin.H{"reply": replyMessage})
	// })

	r.Run()
}
