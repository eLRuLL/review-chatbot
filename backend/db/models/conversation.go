package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WorkflowStatus string

const (
	WELCOME                     WorkflowStatus = "WELCOME"
	PRODUCT_REVIEW              WorkflowStatus = "PRODUCT_REVIEW"
	PRODUCT_REVIEW_CONFIRMATION WorkflowStatus = "PRODUCT_REVIEW_CONFIRMATION"
	PRODUCT_REVIEW_RATING       WorkflowStatus = "PRODUCT_REVIEW_RATING"
	PRODUCT_REVIEW_FINISHED     WorkflowStatus = "PRODUCT_REVIEW_FINISHED"
)

type Conversation struct {
	gorm.Model

	UserID         uint
	WorkflowStatus WorkflowStatus
	Context        datatypes.JSON

	User User `gorm:"foreignKey:UserID"`
}
