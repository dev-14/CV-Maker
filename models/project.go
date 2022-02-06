package models

type Project struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Role        string `json:"role"`
	Language    string `json:"language"`
	Client      bool   `json:"client"`
	ClientName  string `json:"client_name"`
	UserID      int
	User        User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
