package models

import "time"

type WorkExperiences struct {
	ID             int       `json:"id"`
	CompanyName    string    `json:"company_name"`
	EmploymentType string    `json:"employment_type"`
	From           time.Time `json:"from"`
	To             time.Time `json:"to"`
	JobRole        string    `json:"job_role"`
	JobLocation    string    `json:"job_location"`
	Description    string    `json:"description"`
	UserID         int
	User           User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
