package model

type HighSchoolAdmissionPlan struct {
	ID               int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Year             int    `gorm:"not null" json:"year"`
	DistrictType     string `gorm:"type:varchar(50)" json:"district_type"`
	SchoolName       string `gorm:"type:varchar(255);not null" json:"school_name"`
	SchoolLevel      string `gorm:"type:varchar(50)" json:"school_level"`
	OperationNature  string `gorm:"type:varchar(50)" json:"operation_nature"`
	TotalStudents    *int   `json:"total_students"`
	BoardingStudents *int   `json:"boarding_students"`
	DayStudents      *int   `json:"day_students"`
	AdmissionScope   string `gorm:"type:text" json:"admission_scope"`
	Remarks          string `gorm:"type:text" json:"remarks"`
	AcdStudents      int    `gorm:"default:0" json:"acd_students"`
	AcStudents       int    `gorm:"default:0" json:"ac_students"`
	DStudents        int    `gorm:"default:0" json:"d_students"`
}
