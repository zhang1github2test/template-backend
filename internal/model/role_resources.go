package model

type RoleResource struct {
	ResourceID uint `gorm:"primaryKey"`
	RoleID     uint `gorm:"primaryKey"`
}
