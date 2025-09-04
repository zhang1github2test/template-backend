package model

import "time"

type SchoolAdmissionInfo struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	SchoolCode     string    `gorm:"size:20;not null" json:"schoolCode"`
	SchoolName     string    `gorm:"size:100;not null" json:"schoolName"`
	Category       string    `gorm:"size:20;not null" json:"category"`
	TotalScore     int       `gorm:"not null" json:"totalScore"`
	TieBreaker     string    `gorm:"size:255" json:"tieBreaker"`
	AdmissionScope string    `gorm:"size:255" json:"admissionScope"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	Year           int       `gorm:"type:year;not null" json:"year"`
}

func (SchoolAdmissionInfo) TableName() string {
	return "school_admission_info"
}
