package models

import "gorm.io/gorm"

type ProductReview struct {
	gorm.Model

	UserID    uint
	ProductID uint
	Rating    uint

	User    User    `gorm:"foreignKey:UserID"`
	Product Product `gorm:"foreignKey:ProductID"`
}
