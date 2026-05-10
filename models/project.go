package models

import "time"

type Project struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	JobID        uint       `json:"job_id"`
	Job          Job        `gorm:"foreignKey:JobID" json:"job"`
	ClientID     uint       `json:"client_id"`
	Client       Client     `gorm:"foreignKey:ClientID" json:"client"`
	FreelancerID uint       `json:"freelancer_id"`
	Freelancer   Freelancer `gorm:"foreignKey:FreelancerID" json:"freelancer"`
	Status       string     `gorm:"default:'active'" json:"status"`
	Tasks        []Task     `gorm:"foreignKey:ProjectID" json:"tasks"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	SubmissionLink string `json:"submission_link"`
	SubmissionFile string `json:"submission_file"`
}

type Task struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `json:"project_id"`
	Title     string    `json:"title"`
	Status    string    `gorm:"default:'todo'" json:"status"`
	Priority  string    `gorm:"default:'low'" json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
