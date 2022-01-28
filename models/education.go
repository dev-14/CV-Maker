package models

import "time"

type Education struct {
	ID          int       `json:"id"`
	Institution string    `json:"institution"`
	Degree      string    `json:"degree"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Grade       string    `json:"grade"`
	Description string    `json:"description"`
	BelongsToID int
	User        User `gorm:"foreignKey:BelongsToID" json:"-"`
}
