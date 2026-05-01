package models

import "time"

type Application struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Status        string    `json:"status"`
	TanggalDaftar time.Time `gorm:"autoCreateTime" json:"tanggal_daftar"`
	FreelancerID  uint      `json:"freelancer_id"`
	JobID         uint      `json:"job_id"`

	Job Job `gorm:"foreignKey:JobID;references:ID" json:"job"`

	Freelancer Freelancer `gorm:"foreignKey:FreelancerID;references:ID" json:"freelancer"`
}
