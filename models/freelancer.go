package models

import "time"

type Freelancer struct {
	ID                  uint    `gorm:"primaryKey"`
	Nama                string  `gorm:"type:varchar(100)" json:"nama"`
	Email               string  `gorm:"unique" json:"email"`
	Password            string  `json:"password"`
	Gender              string  `json:"gender"`
	Age                 int     `json:"age"`
	FotoProfil          string  `json:"foto_profil"`
	Role                string  `json:"role"`
	Location            string  `json:"location"`
	EducationLevel      string  `json:"education_level"`
	YearsOfExperience   int     `json:"years_of_experience"`
	MonthlySalaryExp    string  `json:"monthly_salary_exp"`
	JobInterest         string  `json:"job_interest"`
	Bio                 string  `json:"bio"`
	Skill               string  `json:"skill"`
	Attachments         string  `json:"attachments"`
	WorkPre             string  `json:"work_pre"`
	Resume              string  `json:"resume"`
	Certificates        string  `json:"certificates"`
	Penilaian           float64 `json:"penilaian"`
	JumlahProyekSelesai int     `json:"jumlah_proyek_selesai"`
	AvgRating           float64
	TotalPenilaian      int
	CreatedAt           time.Time `json:"created_at"`
}
