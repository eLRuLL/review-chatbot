package db

import (
	"backend/db/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func ConnectDatabase() {
	DB, err = gorm.Open(sqlite.Open("chatbot.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	// Migration to create tables for Conversation models
	DB.AutoMigrate(&models.Conversation{}, &models.User{}, &models.Message{}, &models.Product{}, &models.ProductReview{})
}
