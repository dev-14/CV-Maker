package models

type Project struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Role        string `json:"role"`
	Language    string `json:"language"`
	Client      bool   `json:"client"`
	ClientName  string `json:"client_name"`
	BelongsToID int
	User        User `gorm:"foreignKey:BelongsToID" json:"-"`
}
