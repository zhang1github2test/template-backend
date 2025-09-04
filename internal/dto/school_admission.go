package dto

type SchoolAdmissionQueryRequest struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	SchoolName string `form:"schoolName"`
	Year       int    `form:"year"`
	Category   string `form:"category"`
}
