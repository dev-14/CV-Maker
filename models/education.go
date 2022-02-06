package models

import (
	"time"
)

type Education struct {
	ID          int       `json:"id"`
	Institution string    `json:"institution"`
	Degree      string    `json:"degree"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Grade       string    `json:"grade"`
	Description string    `json:"description"`
	UserID      int
	User        User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// func (e []Education) Pagination(c *gin.Context, page int) {
// 	var edu Education
// 	// DB.Where("user_id = ?", e.UserID).Find(&edu).Limit(2).Offset((page - 1) * 5)

// 	c.JSON(200, edu)
// }
